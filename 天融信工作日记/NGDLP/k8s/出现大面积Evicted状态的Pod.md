### 2.8 问题处理
```bash
[root@master ~]# kubectl get po -A
NAMESPACE         NAME                                     READY   STATUS                   RESTARTS        AGE
kube-system       coredns-744574cd6d-qgt6h                 1/1     Running                  4 (80d ago)     237d
kube-system       local-path-provisioner-6f5d79df6-mmcj4   1/1     Running                  6               331d
kube-system       metrics-server-54fd9b65b-zggtx           1/1     Running                  6 (3d20h ago)   304d
topsec-topngdlp   agent-nginx-758c8f8d9b-tr6w4             1/1     Running                  0               16h
topsec-topngdlp   agentserver-5b959445cd-nfmgm             0/1     ContainerStatusUnknown   1               16h
topsec-topngdlp   agentserver-5b959445cd-tkh2b             0/1     Pending                  0               6m26s
topsec-topngdlp   agentsvc-8569b6b5db-8gg46                0/1     Error                    2               16h
topsec-topngdlp   agentsvc-8569b6b5db-flq9b                0/1     Pending                  0               4m11s
topsec-topngdlp   bff-5899c9d4df-nvcx7                     0/1     Pending                  0               3m18s
topsec-topngdlp   bff-5899c9d4df-zcbn2                     0/1     Error                    0               16h
topsec-topngdlp   clickhouse-0                             0/1     Pending                  0               5m17s
topsec-topngdlp   clientscanner-7f6bc68687-s97zh           1/1     Running                  0               16h
topsec-topngdlp   consumersvc-779594bd44-glzg2             0/1     Error                    3               16h
topsec-topngdlp   consumersvc-779594bd44-kn6qd             0/1     Pending                  0               2m12s
topsec-topngdlp   ersvc-6d7d98598d-22r6j                   0/1     Evicted                  0               22s
topsec-topngdlp   ersvc-6d7d98598d-2585p                   0/1     Evicted                  0               43s
topsec-topngdlp   ersvc-6d7d98598d-25kmh                   0/1     Evicted                  0               20s
topsec-topngdlp   ersvc-6d7d98598d-2bs8t                   0/1     Evicted                  0               23s
topsec-topngdlp   ersvc-6d7d98598d-2fz2x                   0/1     Evicted                  0               60s
topsec-topngdlp   ersvc-6d7d98598d-2gq86                   0/1     Evicted                  0               15s
topsec-topngdlp   ersvc-6d7d98598d-2hlmq                   0/1     Evicted                  0               28s
topsec-topngdlp   ersvc-6d7d98598d-2jc6j                   0/1     Evicted                  0               34s
topsec-topngdlp   ersvc-6d7d98598d-2kd76                   0/1     Evicted                  0               11s
topsec-topngdlp   ersvc-6d7d98598d-2lhhc                   0/1     Evicted                  0               7s
topsec-topngdlp   ersvc-6d7d98598d-2nf5f                   0/1     Evicted                  0               5s
topsec-topngdlp   ersvc-6d7d98598d-2nkj2                   0/1     Evicted                  0               51s
topsec-topngdlp   ersvc-6d7d98598d-2rgzf                   0/1     Evicted                  0               54s
topsec-topngdlp   ersvc-6d7d98598d-2rwrp                   0/1     Evicted                  0               34s
topsec-topngdlp   ersvc-6d7d98598d-2sg5s                   0/1     Evicted                  0               16s
topsec-topngdlp   ersvc-6d7d98598d-2tpnk                   0/1     Evicted                  0               1s
```

当出现这种大量的pod错误的时候

#### 一、快速定位根本原因
按顺序执行以下排查步骤，重点关注节点资源和Pod 事件。

1. 检查节点状态与资源压力
```bash
kubectl get nodes -o wide
kubectl top nodes   # 需要 metrics-server 正常工作
```

- 是否有节点 NotReady？

- 节点内存/CPU 使用率是否接近 100%？

- 如果 top nodes 不可用，直接看节点描述：

```bash
kubectl describe node <node-name> | grep -A 10 "Conditions"
```

重点关注 MemoryPressure、DiskPressure、PIDPressure 是否为 True。

2. 查看被驱逐 Pod 的原因
选一个 Evicted 的 Pod 看详情：
```bash
kubectl describe pod ersvc-6d7d98598d-2tpnk -n topsec-topngdlp
```

在输出末尾的 Events 中会看到类似：

```bash
The node was low on resource: memory.
```

3. 查看 Pending Pod 的调度失败原因
```bash
kubectl describe pod agentserver-5b959445cd-tkh2b -n topsec-topngdlp
```

常见原因：

- 0/1 nodes are available: 1 Insufficient memory → 节点内存不足

- 0/1 nodes are available: pod has unbound immediate PersistentVolumeClaims → PVC 未绑定

- 0/1 nodes are available: node(s) had untolerated taint → 存在污点

4. 检查 PVC 是否正常（尤其是 clickhouse）
```bash
kubectl get pvc -n topsec-topngdlp
kubectl describe pvc <clickhouse-pvc-name> -n topsec-topngdlp
```

5. 查看集群整体资源分配
```bash
kubectl describe nodes | grep -E "Allocated|Capacity" -A 5
```

计算每个节点的 Requests 是否接近或超过 Capacity。

6. 检查是否有镜像拉取失败（针对 Error 状态 Pod）
```bash
kubectl describe pod agentsvc-8569b6b5db-8gg46 -n topsec-topngdlp
```
如果 Events 中有 Failed to pull image 或 ImagePullBackOff，则需修正镜像地址或添加 imagePullSecret。

#### 三、紧急恢复建议
1. 清理被驱逐的 Pod（避免占用调度资源）
```bash
kubectl delete pod -n topsec-topngdlp --field-selector=status.phase=Failed

# 或者更精确地删除所有 Evicted Pod
kubectl get pod -n topsec-topngdlp | grep Evicted | awk '{print $1}' | xargs kubectl delete pod -n topsec-topngdlp
```

2.临时降低有问题的 Deployment 副本数（如 ersvc、agentserver）
```bash
kubectl scale deployment ersvc -n topsec-topngdlp --replicas=0
```

3. 调整资源 requests/limits
如果节点内存紧张，降低 Pod 的 requests.memory，或为内存限制设置合适的值（例如从 2Gi 降到 1Gi）。

4. 扩容节点
添加新节点或为现有节点增加资源（云上可垂直扩容实例规格）。

5. 修复 PVC/StorageClass
如果 clickhouse 等需要持久存储，确保有可用的 StorageClass 并且 PVC 能成功绑定。

6. 一键清理 Docker 无用数据（最有效）
```bash
docker system prune -a -f --volumes
```

这会删除：

- 停止的容器
- 未使用的网络
- 悬空镜像
- 未使用的卷（注意：--volumes 会删除未被任何容器引用的卷，谨慎但通常安全）

如果不确定，可以先只清理镜像和容器：

```bash
docker system prune -a -f
```

7. 清理 K3s 的陈旧镜像和缓存
```bash
# 清理已退出的容器
docker rm $(docker ps -aq --filter status=exited) 2>/dev/null

# 清理悬挂镜像
docker image prune -a -f

# 清理未使用的卷
docker volume prune -f
```
