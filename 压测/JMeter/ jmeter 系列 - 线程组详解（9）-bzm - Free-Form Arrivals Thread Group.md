> 本文由 [简悦 SimpRead](http://ksria.com/simpread/) 转码， 原文地址 [www.cnblogs.com](https://www.cnblogs.com/jiushao-ing/p/17647289.html)

bzm - Free-Form Arrivals Thread Group
-------------------------------------

### 介绍：

 顾名思义，相当于自由形式的 Arrivals Thread Group，它只是提供了自由形式的时间表的能力。相当于我们可以更灵活的控制  每分钟 / 每秒钟的请求数。

### 页面说明：

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230821231830267-896811038.png)

Threads Schedule（线程场景）:

*   Start value: 开始时的用户数
*   End value：结束时的用户数
*   Duration：持续时间

Time Unit：时间单位选择（影响所有配置项）

*   minutes：分钟
*   second：秒

Thread Iterations Limit: 线程迭代次数限制。如果我们只需运行每个用户一次以模拟用户的实际行为，则设置为 1；设置为空，表示循环，直到调度结束。

Log Threads Status into File：将线程启动和线程停止事件保存为日志文件。

Concurrency Limit: 最大并发数限制，避免出现内存不足。

 来看 2 个例子，当然实际上会用更多的场景配置

例子 1：

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230821232517541-1906447172.png)

例子 2：

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230821232904833-277694563.png)