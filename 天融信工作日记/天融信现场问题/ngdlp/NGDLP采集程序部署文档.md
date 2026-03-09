### 安装部署采集程序sfc

1. 解压采集程序包

   ```shell
   tar -zxvf ngdlp-sfc-amd64.tar.gz

   cd ngdlp-sfc-amd64
   ```
2. 加载sfc镜像

   ```shell
   /home/topsec/cache/gcrpi -f ngdlp-sfc-1.0.26-sfc.tar
   ```
3. 修改./chart/values.yaml sfc的nfs、部署节点等配置

   ```shell
   storage:
     create: enable
     pvName: topsec-cloud-sfc
     pvcName: topsec-cloud-sfc
     storage: 4T
     # 修改储存方式: "local-path", "nfs" [default: local-path]
     storageClass: local-path 
     # nfs地址及路径
     nfsServer: ""
     nfsPath: ""
     localPath: /home/topsec/data/
     localPathNode: master
     ngdlpPVCName: topsec-cloud
   sfc:
     # sfc运行的节点
     nodeName: ""
     image: ngdlp/ngdlp
     tag: "1.0.26-sfc"
     resources:
       requests:
         cpu: "0.2"
         memory: "200Mi"
       limits:
         cpu: "2"
         memory: "2Gi"
   ```
4. 修改采集程序配置./chart/templates/configmap.yaml

   ```yaml
   apiVersion: v1
   kind: ConfigMap
   metadata:
   name: ngdlp-sfc-config
   namespace: {{ .Release.Namespace }}
   data:
   sensitive_collect.yaml: |
    extract:
      enable: true  # 启用停用
      bootrun: true  # 是否在启动时运行一次采集程序，默认true
      cron_spec: "@every 1h"  # cron表达式， @every 1h 每小时执行一次，@daily 每天0点执行，@weekly 每周一0点执行
      force_extract: false  # 是否强制提取，会给正在提取的文件发送提取任务
      extract_task_interval: "5s"  # 提取间隔 5s，如果管理中心上传压力大，可以调大
      file_suffix: ""  # 文件后缀 多个以英文逗号分隔
      max_file_size: ""  # 文件最大值 100M
      min_file_size: ""  # 文件最小值 1M
      start_at: ""  # 文件最新发现开始时间 格式：2023-01-01 00:00:00
      end_at: ""  # 文件最新发现结束时间 格式：2023-01-01 23:59:59
      since: "3h" # 提取3小时之前到现在的文件，配置了这个参数start_at/end_at会失效
    archive:
      hash_file: "/home/topsec/data/sfc_archive/workspace/hash.txt"  # 哈希文件路径,支持md5、文件名过滤，默认按换行读取
   ```
5. 安装采集程序sfc（默认启动sfc就会下发采集任务，并且每小时执行一次）

   ```shell
   helm install -n topsec-topngdlp sfc ./chart/
   ```

### 升级采集程序

1. 修改./chart/values.yaml sfc的nfs、部署节点等配置，跟第一次安装一样，修改完了再升级
2. 加载sfc镜像

   ```shell
   /home/topsec/cache/gcrpi -f ngdlp-sfc-1.0.26-sfc.tar
   ```
3. 升级程序

```shell
helm upgrade -n topsec-topngdlp sfc ./chart/

# 查看sfc是否重启了，如果没重启需要手动重启: kubectl rollout restart -n topsec-topngdlp deployment sfc
kubectl get pod -n topsec-topngdlp | grep sfc


```

    

### 卸载采集程序

```shell
helm uninstall -n topsec-topngdlp sfc

# 先删除pvc 再删除pv
kubectl delete -n topsec-topngdlp pvc topsec-cloud-sfc
kubectl delete pv topsec-cloud-sfc
```

### 采集程序命令

```shell
# 查看采集程序日志
kubectl logs -f -n topsec-topngdlp -l app=sfc

# 打印提取任务进度
kubectl exec -it -n topsec-topngdlp deployments/sfc -- bin/ngdlp --sp

# 归档数据，默认在/home/topsec/data/sfc_archive/目录，执行之前需要将指定的hash.txt放到/home/topsec/data/sfc_archive/workspace/下
# hash.txt 默认换行读取，支持文件md5跟文件名
kubectl exec -it -n topsec-topngdlp deployments/sfc -- bin/ngdlp --sa

# 归档数据，数据多时，最好后台运行(推荐)
nohup kubectl exec -it -n topsec-topngdlp deployments/sfc -- bin/ngdlp --sa > archive.log 2>&1 &
```

### 修改采集程序配置

```shell
# 只要不执行helm upgrade则会一直生效
kubectl edit cm -n topsec-topngdlp ngdlp-sfc-config
```

### 采集程序修改日志级别

```shell
# 修改为info
kubectl edit -n topsec-topngdlp deploy sfc
=>
- name: LOGGER_LEVEL
  value: info
```
