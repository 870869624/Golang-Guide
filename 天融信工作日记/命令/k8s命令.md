下面把“sdsp 命名空间下的 cleansvc 服务”最常用的 10 条命令做成一条“口诀”，  
背下来就能覆盖 90 % 日常排查与变更场景。口诀分三行，每行 3~4 个关键字，对应一条完整命令。

------------------------------------------------
# 口诀（从左到右读，空格就是分隔）：

1 看  2 进  3 志  4 重启  
5 扩  6 缩  7 灰  8 镜  
9 配  0 密  d 删  s 选  

------------------------------------------------
## 逐字对应命令（把 %% 替换成 cleansvc，-n sdsp 固定放在尾缀）：

### 1 看  → 看现状  
kubectl get deploy %% -n sdsp

### 2 进  → 进容器  
kubectl exec -it deploy/%% -n sdsp -- bash

### 3 志  → 看日志  
kubectl logs -f deploy/%% -n sdsp --tail=200

### 4 重启 → 滚动重启  
kubectl rollout restart deploy %% -n sdsp

### 5 扩  → 扩容副本  
kubectl scale deploy %% --replicas=5 -n sdsp

### 6 缩  → 缩容副本  
kubectl scale deploy %% --replicas=1 -n sdsp

### 7 灰  → 灰度镜像  
kubectl set image deploy %% *=registry.example.com/cleansvc:v2.1.3 -n sdsp

### 8 镜  → 查镜像  
kubectl get deploy %% -n sdsp -o jsonpath='{.spec.template.spec.containers[0].image}'

### 9 配  → 改配置（ConfigMap）  
kubectl edit cm %% -n sdsp          # 改完执行 4 重启

### 0 密  → 改秘钥（Secret）  
kubectl edit secret %% -n sdsp      # 改完执行 4 重启

### d 删  → 删服务（谨慎）  
kubectl delete deploy %% -n sdsp

### s 选  → 快速给后面所有命令打别名  
alias k='kubectl -n sdsp'           # 以后直接 k get deploy cleansvc

------------------------------------------------
记忆技巧  
- 第一行“看进志重启”是排查四板斧，顺序记牢。  
- 第二行“扩缩灰镜像”是变更四板斧。  
- 第三行“配密删选”是高级/危险操作。  
- 所有命令统一带 `-n sdsp`，不会误操作别的命名空间。

把口诀抄到终端 profile 或者贴到桌面，一周后就忘不掉了。