> 本文由 [简悦 SimpRead](http://ksria.com/simpread/) 转码， 原文地址 [www.cnblogs.com](https://www.cnblogs.com/jiushao-ing/p/17646737.html)

bzm - Arrivals Thread Group
---------------------------

Arrival：到来，抵达

### 介绍

这个线程组使用 “arrivals” 调度作为一种表达负载的方式。“arrivals”表示线程迭代开始。如果所有现有线程在迭代过程中都很忙，它将创建新线程。  
注意，恒定的到达率意味着增加并发性，所以要小心你输入的值。使用 “Concurrency Limit” 字段作为安全阀，以防止内存不足。

### 主要功能：

*   每秒 / 每分钟 请求数
    
*   阶梯控制

### 页面说明：

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230821225936453-722084074.png)

*   Target Rate (arrivals/min) : 目标率 （每分钟请求数）
*   Ramp UP Time（min） ： 在多少秒内到达目标请求数
*   Ramp-Up Steps Count : 启动之后到达目标并发线程数的 阶梯数
*   Hold Target Rate Time(min) : 到达目标请求数之后，持续运行多长时间
*   Time Unit：minutes seconds : 时间单元：分 / 秒
*   Thread iterations Limit : 线程循环次数限制
*   Log Threads Status into File : 保存线程状态至文件
*   Concurrency Limit : 最大线程数限制

 来看一个例子，设置如下：

每分钟 60 个请求，相当于每秒一个请求，Ramp-Up 启动时间和阶梯数 Ramp-Up step Count 都设置为 0

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230821222402277-1057853148.png)

添加监听器 Active Threads Over Time，聚合报告

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230821222605177-277598581.png)

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230821222706902-471174532.png)

可以看到只启动了一个线程，2 分钟 120 个请求

 再来看一个例子：

每分钟请求 60 次，相当于每秒 12 次，Ramp-Up 启动时间 1 分钟，阶梯数 3，持续 2 分钟

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230821223621980-1665839171.png)

添加监听器 Transactions per Second，即每秒事务数

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230821224342810-940679609.png)

可以看到事务数基本上是在每秒 12 左右