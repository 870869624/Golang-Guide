# K8s调度Service至其他节点

## 📋 使用场景

在以下场景中，需要将Service调度到特定的K8s节点：

- **资源隔离**：某些服务需要运行在高性能节点或专用节点上
- **故障排查**：某个节点出现问题，需要将服务迁移到健康节点
- **硬件依赖**：服务需要特定硬件（如GPU、特定CPU架构）
- **网络要求**：服务需要访问节点本地资源或特定网络环境
- **性能优化**：将服务调度到负载较低的节点

---

## 🔧 核心命令

### 方法一：编辑Deployment（推荐）

```bash
kubectl edit deploy agentsvc -n topsec-topaiop
```

![编辑Deployment](image/调度svc至其他节点/1790060618677.png)

### 方法二：直接修改Pod（临时方案）

```bash
# 查看Pod当前所在节点
kubectl get pod <pod-name> -n <namespace> -o wide

# 编辑Pod（不推荐，Pod重启后会失效）
kubectl edit pod <pod-name> -n <namespace>
```

---

## 📝 关键字段说明

在 `kubectl edit` 打开的YAML文件中，找到 `spec.template.spec` 部分，可以添加以下字段：

### 1. nodeName（强制指定节点）

```yaml
spec:
  template:
    spec:
      nodeName: worker-node-01  # 直接指定Pod运行的节点名称
```

**特点**：
- ✅ 最简单直接的方式
- ❌ 绕过调度器，硬性绑定
- ⚠️ 如果指定节点不可用，Pod会一直处于Pending状态

### 2. nodeSelector（标签选择）

```yaml
spec:
  template:
    spec:
      nodeSelector:
        disktype: ssd           # 选择带有disktype=ssd标签的节点
        kubernetes.io/arch: arm64  # 选择arm64架构的节点
```

**特点**：
- ✅ 简单且灵活的标签匹配
- ⚠️ 需要节点事先打好对应标签
- 📌 适合简单的调度需求

**节点打标签命令**：

```bash
# 给节点添加标签
kubectl label nodes <node-name> disktype=ssd

# 查看节点标签
kubectl get nodes --show-labels

# 删除节点标签
kubectl label nodes <node-name> disktype-
```

### 3. nodeAffinity（节点亲和性，推荐）

```yaml
spec:
  template:
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: kubernetes.io/arch
                operator: In
                values:
                - arm64
                - amd64
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 1
            preference:
              matchExpressions:
              - key: node-role.kubernetes.io/worker
                operator: In
                values:
                - high-performance
```

**特点**：
- ✅ 最强大的调度方式，支持硬性要求（required）和偏好（preferred）
- ✅ 支持多种操作符：In, NotIn, Exists, DoesNotExist, Gt, Lt
- 📌 生产环境推荐使用

### 4. tolerations（容忍污点）

如果节点设置了污点（taint），需要配置容忍：

```yaml
spec:
  template:
    spec:
      tolerations:
      - key: "special-node"
        operator: "Equal"
        value: "true"
        effect: "NoSchedule"
```

**查看节点污点**：

```bash
kubectl describe node <node-name> | grep Taints
```

**给节点添加/删除污点**：

```bash
# 添加污点
kubectl taint nodes <node-name> special-node=true:NoSchedule

# 删除污点
kubectl taint nodes <node-name> special-node:NoSchedule-
```

---

## 🎯 完整操作示例

### 示例1：将agentsvc调度到特定节点

```bash
# 1. 查看当前节点列表
kubectl get nodes -o wide

# 2. 查看当前Pod所在节点
kubectl get pods -n topsec-topaiop -o wide | grep agentsvc

# 3. 编辑Deployment
kubectl edit deploy agentsvc -n topsec-topaiop

# 4. 在spec.template.spec下添加nodeName
spec:
  template:
    spec:
      nodeName: worker-node-02  # 添加这一行
      containers:
      - name: agentsvc
        image: xxx

# 5. 保存退出，Pod会自动重建并调度到新节点

# 6. 验证调度结果
kubectl get pods -n topsec-topaiop -o wide | grep agentsvc
```

### 示例2：使用nodeSelector调度到arm64节点

```bash
# 1. 查看节点架构标签
kubectl get nodes --show-labels | grep arch

# 2. 如果没有标签，手动添加
kubectl label nodes arm-worker-01 kubernetes.io/arch=arm64

# 3. 编辑Deployment
kubectl edit deploy agentsvc -n topsec-topaiop

# 4. 添加nodeSelector
spec:
  template:
    spec:
      nodeSelector:
        kubernetes.io/arch: arm64

# 5. 保存并验证
kubectl rollout status deploy agentsvc -n topsec-topaiop
kubectl get pods -n topsec-topaiop -o wide
```

### 示例3：使用节点亲和性（生产推荐）

```bash
kubectl edit deploy agentsvc -n topsec-topaiop
```

添加以下配置：

```yaml
spec:
  template:
    spec:
      affinity:
        nodeAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
            nodeSelectorTerms:
            - matchExpressions:
              - key: kubernetes.io/arch
                operator: In
                values:
                - arm64
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            preference:
              matchExpressions:
              - key: node-role.kubernetes.io/ai
                operator: In
                values:
                - "true"
      tolerations:
      - key: "ai-workload"
        operator: "Equal"
        value: "true"
        effect: "NoSchedule"
```

---

## ✅ 验证调度结果

```bash
# 1. 查看Pod所在节点
kubectl get pods -n <namespace> -o wide | grep <service-name>

# 2. 查看Pod详细信息
kubectl describe pod <pod-name> -n <namespace> | grep Node:

# 3. 查看调度事件
kubectl get events -n <namespace> --sort-by='.lastTimestamp' | grep <pod-name>

# 4. 查看滚动更新状态
kubectl rollout status deploy <deploy-name> -n <namespace>

# 5. 验证服务是否正常
curl http://<service-ip>:<port>/health
```

---

## ↩️ 回滚操作

### 方法一：删除nodeName/nodeSelector

```bash
# 1. 编辑Deployment
kubectl edit deploy agentsvc -n topsec-topaiop

# 2. 删除之前添加的nodeName、nodeSelector或affinity字段

# 3. 保存退出，Pod会重新由调度器分配
```

### 方法二：使用rollout undo

```bash
# 查看历史版本
kubectl rollout history deploy agentsvc -n topsec-topaiop

# 回滚到上一个版本
kubectl rollout undo deploy agentsvc -n topsec-topaiop

# 回滚到指定版本
kubectl rollout undo deploy agentsvc -n topsec-topaiop --to-revision=2
```

---

## ⚠️ 注意事项

1. **nodeName vs nodeSelector**：
   - `nodeName` 是硬绑定，绕过调度器
   - `nodeSelector` 和 `affinity` 会通过调度器，更灵活

2. **Pod重建**：
   - 修改Deployment后，旧Pod会被终止，新Pod在新节点创建
   - 如果有数据持久化需求，确保使用PersistentVolume

3. **节点资源**：
   - 目标节点必须有足够的CPU和内存资源
   - 使用 `kubectl describe node <node-name>` 查看资源使用情况

4. **网络策略**：
   - 确认目标节点的网络策略允许Pod通信
   - 检查Service的Endpoints是否正常更新

5. **污点与容忍**：
   - 如果节点有污点，必须配置对应的tolerations
   - 常见污点：`NoSchedule`、`PreferNoSchedule`、`NoExecute`

6. **多副本调度**：
   - 多个副本可能会被调度到不同节点
   - 使用 `podAntiAffinity` 可以避免副本在同一节点

---

## 🔍 故障排查

### Pod一直处于Pending状态

```bash
# 查看Pod事件
kubectl describe pod <pod-name> -n <namespace>

# 常见原因：
# - 节点资源不足
# - 节点标签不匹配
# - 节点有污点但未配置tolerations
# - 指定的nodeName不存在
```

### Pod调度到节点但无法启动

```bash
# 查看Pod日志
kubectl logs <pod-name> -n <namespace>

# 查看节点状态
kubectl get nodes
kubectl describe node <node-name>

# 检查镜像拉取
docker images | grep <image-name>  # 在目标节点执行
```

### Service无法访问

```bash
# 检查Endpoints
kubectl get endpoints <service-name> -n <namespace>

# 检查Service配置
kubectl get svc <service-name> -n <namespace> -o yaml

# 测试网络连通性
kubectl exec -it <pod-name> -n <namespace> -- curl <service-ip>:<port>
```

---

## 💡 最佳实践

1. **优先使用nodeAffinity**：比nodeSelector更灵活，支持软硬结合
2. **避免硬编码nodeName**：除非有特殊需求，否则让调度器决定
3. **合理设置资源限制**：确保调度器能正确评估节点资源
4. **使用标签管理节点**：为不同用途的节点打上清晰标签
5. **监控节点负载**：避免所有服务都调度到同一节点

---

## 📚 相关命令速查

```bash
# 查看节点信息
kubectl get nodes -o wide
kubectl describe node <node-name>

# 节点标签管理
kubectl get nodes --show-labels
kubectl label nodes <node-name> key=value
kubectl label nodes <node-name> key-

# 节点污点管理
kubectl taint nodes <node-name> key=value:effect
kubectl taint nodes <node-name> key:effect-

# Pod调度查看
kubectl get pods -o wide
kubectl describe pod <pod-name> | grep Node:

# Deployment管理
kubectl edit deploy <name> -n <namespace>
kubectl rollout status deploy <name> -n <namespace>
kubectl rollout history deploy <name> -n <namespace>
kubectl rollout undo deploy <name> -n <namespace>
```
