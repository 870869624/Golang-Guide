# NFT / 区块链系统设计

> 核心文件：`app/service/digitalService/createTokenId.go`、`digitalMedal.go`、`transferNft.go`
> 合约ABI：`build/NFT.abi`、`PohAbi/Pohabi.abi`、`pridowtoken/PridowToken.abi`
> 区块链SDK：go-conflux-sdk v1.5.9

---

## 一、整体架构

```
用户支付成功
    │
    ▼
创建铸造记录（DB状态：铸造中）
    │
    ▼
发送消息到 RabbitMQ（custingNft 队列）
    │
    ▼
MQ消费者（独立goroutine）
    │
    ▼
1. 解析铸造ID，查询铸造记录/部落/勋章
    │
    ▼
2. 调用 Conflux 智能合约 Mint
    │  ├── sdk.NewAccountManager(keystorePath, networkID)
    │  ├── client.AccountManager.Unlock(address, password)
    │  └── instance.Mint(to, tokenId, uri)  // ERC721 Mint
    │
    ▼
3. 等待交易回执（WaitForTransactionReceipt）
    │  ├── pending → 继续轮询
    │  ├── success → 解析Logs获取TokenId
    │  └── failed → 标记铸造失败
    │
    ▼
4. 上传元数据到 IPFS（Pinata API）
    │  ├── 构建metadata JSON（name, description, image, attributes）
    │  └── 返回CID（内容标识符）
    │
    ▼
5. 回写DB：tokenId、交易hash、铸造状态
    │
    └── 铸造完成
```

---

## 二、Conflux 区块链集成

### 连接方式

```go
client := sdk.NewClient("https://main.confluxrpc.com")
accountManager := sdk.NewAccountManager("./formal-keystore/", networkID)
accountManager.Unlock(address, password)
```

- **RPC节点**：Conflux主网 `https://main.confluxrpc.com`
- **账户管理**：Keystore文件 + 密码解锁（以太坊标准Keystore格式）
- **网络ID**：区分主网/测试网

### 合约交互

通过 `build.NewNFT(contractAddr, client)` 加载编译好的ABI绑定：

```go
// 获取合约实例
instance := build.NewNFT(common.HexToAddress(contractAddr), client)

// 铸造NFT（Mint）
txHash, err := instance.Mint(toAddress, big.NewInt(tokenId), tokenURI)

// 转移NFT（TransferFrom）
txHash, err := instance.TransferFrom(from, to, big.NewInt(tokenId))

// 授权合约操作
instance.SetApprovalForAll(contractAddress, true)

// 检查授权状态
approved, _ := instance.IsApprovedForAll(owner, contractAddress)
```

### 地址类型转换

Conflux使用base32编码的地址格式，需要转换：

```go
// base32 → hex（合约调用需要hex地址）
hexAddress := cfxaddress.MustNew(base32Address).MustGetCommonAddress().Hex()
```

---

## 三、交易回执解析（关键步骤）

```go
receipt, err := client.Transaction().GetTransactionReceipt(txHash)

// 遍历Logs，解析Transfer事件
for _, log := range receipt.Logs {
    // ParseTransfer 解析ERC721 Transfer事件
    // 获取 TokenId
    tokenId := parseTransferEvent(log)
}
```

**这是标准的区块链交互模式**：
1. 发送交易 → 获取 txHash
2. 轮询交易回执（WaitForTransactionReceipt）
3. 解析事件日志 → 获取链上状态变更
4. 持久化到数据库

---

## 四、NFT铸造的异步化（MQ解耦）

### 为什么用MQ？

- **链上操作耗时长**：交易确认需要几秒到几分钟
- **不可阻塞用户响应**：用户支付后应立即看到"铸造中"状态
- **失败可重试**：MQ消息重入队即可

### 消费者设计

```go
func ConsumeMessagesCreateTokenId() {
    // 订阅 custingNft 队列
    mq.SubscribeToQueue("custingNft", func(msg amqp.Delivery) {
        defer func() {
            if r := recover(); r != nil {
                log.Error("panic recovered:", r)
            }
        }()
        
        // 1. 解析消息
        // 2. 查DB获取铸造记录
        // 3. 调用合约Mint
        // 4. 等待回执
        // 5. 上传IPFS
        // 6. 更新DB
        
        msg.Ack(false)  // 手动ACK
    })
}
```

### 失败处理

```
铸造失败
    │
    ├── 合约调用失败 → Nack重入队 → 下次重试
    │
    ├── 交易pending超时 → 标记铸造失败 → 触发退款
    │
    └── IPFS上传失败 → 重试或使用备用存储
```

---

## 五、NFT转移

```go
func TransferNft(from string, to string, tokenId int64) (string, error) {
    fromAddress, instance, client, _ := nftserver.Transfer(from)
    
    // 1. 授权合约可以操作NFT
    instance.SetApprovalForAll(contractAddress, true)
    time.Sleep(2 * time.Second)  // 等待授权交易上链
    
    // 2. 验证授权状态
    approved, _ := instance.IsApprovedForAll(fromAddress, contractAddress)
    if !approved {
        return "", errors.New("approval failed")
    }
    
    // 3. 执行转移
    txHash, _ := instance.TransferFrom(fromAddress, toAddress, big.NewInt(tokenId))
    return txHash, nil
}
```

**ERC721 标准三步曲**：`SetApprovalForAll` → `IsApprovedForAll` → `TransferFrom`

---

## 六、铸造冷却机制

```go
// 判断铸造是否完成（30分钟冷却）
isMinted := record.CreatedAt.TimeX().Before(time.Now().Add(-30 * time.Minute))
```

用户支付后30分钟内，前端展示"铸造中"状态。超时后标记为铸造失败，触发退款流程。

---

## 七、IPFS 元数据上传

```go
// Pinata API 上传
func UploadToIPFS(metadata map[string]interface{}) (string, error) {
    // 1. 构建 metadata JSON
    //    {
    //      "name": "UEFUN 数字藏品 #001",
    //      "description": "...",
    //      "image": "ipfs://Qm...",
    //      "attributes": [{"trait_type": "series", "value": "..."}]
    //    }
    // 2. POST 到 Pinata API（pinJson 端点）
    // 3. 返回 CID（如 QmXxx...）
}
```

---

## 八、三个智能合约

| 合约 | 用途 |
|------|------|
| NFT | ERC721标准，铸造、转移、授权 |
| PoH | Proof of Humanity 证明 |
| PridowToken | ERC20代币合约 |

### NFT 合约 ABI 方法

| 方法 | 类型 | 说明 |
|------|------|------|
| `Mint(to, tokenId, uri)` | 写 | 铸造NFT |
| `TransferFrom(from, to, tokenId)` | 写 | 转移NFT |
| `SetApprovalForAll(operator, approved)` | 写 | 授权/取消授权 |
| `IsApprovedForAll(owner, operator)` | 读 | 查询授权状态 |
| `BalanceOf(owner)` | 读 | 查询持有数量 |
| `TokenURI(tokenId)` | 读 | 查询元数据URI |

---

## 九、面试常见问题

### Q1: 铸造失败但用户已付款怎么办？
> 标记铸造记录为"失败"，通过退款流程退还用户支付金额。30分钟冷却机制兜底超时场景。

### Q2: 区块链交易pending状态持续很久怎么办？
> 设置30分钟超时，超时标记为失败。正常Conflux主网交易确认在几秒到几十秒内。

### Q3: Gas费谁出？
> 平台代付。使用平台Keystore账户发起交易，Gas从平台账户扣除。

### Q4: 为什么选Conflux而不是以太坊？
> Conflux交易费极低（几乎为零），适合高频铸造场景。且Conflux兼容EVM，go-ethereum的ABI绑定可以直接使用。

### Q5: IPFS存储可靠吗？
> 使用Pinata（商业IPFS Pinning服务），保证元数据的持久可用性。同时元数据也存储在数据库中做备份。

### Q6: NFT转移为什么需要先授权？
> ERC721标准要求：如果合约A要操作用户的NFT，用户必须先授权（`SetApprovalForAll`）。这是区块链的安全机制，防止未经授权的转移。