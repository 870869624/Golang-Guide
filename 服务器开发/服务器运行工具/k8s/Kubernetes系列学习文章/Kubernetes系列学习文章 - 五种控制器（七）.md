# 一、为什么要有控制器

K8S是容器资源管理和调度平台，容器跑在Pod里，Pod是K8S里最小的单元。所以，这些Pod作为一个个单元我们肯定需要去操作它的状态和生命周期。那么如何操作？这里就需要用到控制器了。

这里一个比较通俗的公式：应用APP = 网络 + 载体 + 存储

![1770274138082](image/Kubernetes系列学习文章-五种控制器（七）/1770274138082.png)

这里应用一般分为无状态应用、有状态应用、守护型应用、批处理应用这四种。

**无状态应用**：应用实例不涉及事务交互，不产生持久化数据存储在本地，并且多个应用实例对于同一个请求响应的结果是完全一致的。举例：nginx或者tomcat

**有状态应用**：有状态服务可以说是需要数据存储功能的服务或者指多线程类型的服务、队列等。举例：mysql数据库、kafka、redis、zookeeper等。

**守护型应用**：类似守护进程一样，长期保持运行，监听持续的提供服务。举例：ceph、logstash、fluentd等。

**批处理应用**：工作任务型的服务，通常是一次性的。举例：运行一个批量改文件夹名字的脚本。

这些类型的应用服务如果是安装在传统的物理机或者虚拟机上，那么我们一般会通过人肉方式或者自动化工具的方式去管理编排。但是这些服务一旦容器化了跑在了Pod里，那么就应该按照K8S的控制方式来管理了。上一篇文章我们讲到了编排，那么K8S靠什么具体的操作来做编排？答案就是这些控制器。

# 二、K8S有哪些控制器
既然应用的类型有上面说的这些无状态、有状态的，那么K8S肯定要实现一些控制器来专门处理对应类型的应用。总体来说，K8S有五种控制器，分别对应处理无状态应用、有状态应用、守护型应用和批处理应用。

## 1. Deployment
Deployment中文意思为部署、调度，通过Deployment我们能操作RS（ReplicaSet），你可以简单的理解为它是一种通过yml文件的声明，在Deployment 文件里可以定义Pod数量、更新方式、使用的镜像，资源限制等。无状态应用都用Deployment来创建，例：

```yaml
apiVersion: extensions/v1beta1
kind: Deployment   # 定义是Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 2
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.8.0
        ports:
        - containerPort: 80
```

## 2. StatefulSet
StatefulSet的出现是K8S为了解决 “有状态” 应用落地而产生的，Stateful这个单词本身就是“有状态”的意思。之前大家一直怀疑有状态应用落地K8S的可行性，StatefulSet很有效解决了这个问题。有状态应用一般都需要具备一致性，它们有固定的网络标记、持久化存储、顺序部署和扩展、顺序滚动更新等等。总结两个词就是需要稳定、有序。

那么StatefulSet如何做到Pod的稳定、有序？具体有了哪些内在机制和方法？主要概况起来有这几个方面：

- 给Pod一个唯一和持久的标识（例:Pod name）
- 给予Pod一份持久化存储
- 部署Pod都是顺序性的，0 ~ N-1
- 扩容Pod必须前面的Pod还存在着
- 终止Pod，后面Pod也一并终止
- 举个例子：创建了zk01、zk02、zk03 三个Pod，zk01就是给的命名，如果要扩容zk04，那么前面01、02、03必须存在，否则不成功；如果删除了zk02，那么zk03也会被删除。

## 3. DaemonSet
Daemon本身就是守护进程的意思，那么很显然DaemonSet就是K8S里实现守护进程机制的控制器。比如我们需要在每个node里部署fluentd采集容器日志，那么我们完全可以采用DaemonSet机制部署。它的作用就是能确保全部（或者你指定的node数里）运行一个fluentd Pod副本。当有 node加入集群时，也会为他们新增一个 Pod 。当有 node从集群移除时，这些 Pod 也会被回收。删除 DaemonSet 将会删除它创建的所有 Pod。

所以，你可以想象，DaemonSet 特别适合运行那些静默后台运行的应用，而且是连带性质的，非常方便。

## 4. Job
Job就是任务，我们不用K8S，批处理的运行一些自动化脚本或者跑下ansible也是经常的事儿。那么在K8S里运行批处理任务我们用Job即可。执行一次的任务，它保证批处理任务的一个或多个Pod成功结束。

## 5. CronJob
在IT环境里，经常遇到一些需要定时启动运行的任务。传统的linux里我们执行定义crontab即可，那么在K8S里我们就可以用到CronJob控制器。其实它就是上面Job的加强版，带时间定点运行的。

例子：每一分钟输出一句，2019-08-25 08:08:08 Hello K8S！


```yaml
apiVersion: batch/v1beta1
kind: CronJob  # 定义CronJob类型
metadata:
  name: hello
spec:
  schedule: "*/1 * * * *"   # 定义定时任务运行
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: hello
            image: busybox
            args:
            - /bin/sh
            - -c
            - date; echo Hello K8S!
          restartPolicy: OnFailure
```

**总结**：以上就是K8S五种控制器的介绍，这五种控制器的存在对标的就是四种类型应用的编排处理。有人会问这五种控制器到底怎么用呢？很简单，还是通过编写运行YAML文件来操控。网上有很多控制器运行的YAML例子，你参考一个部署一个，基本其他的你也会了。