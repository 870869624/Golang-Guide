> 本文由 [简悦 SimpRead](http://ksria.com/simpread/) 转码， 原文地址 [www.cnblogs.com](https://www.cnblogs.com/jiushao-ing/p/17644150.html)

Ultimate Thread Group
---------------------

### 介绍：

“Ultimate” 意味着将不会有进一步的线程组插件的需要。每个人都可以在 JMeter 用:

*   无限数量的线程场景
*   每个线程场景的 ramp-up time, shutdown time, flight time
*   当然，还有值得信赖的负载预览图

### 添加方式：

右键测试计划 -> 添加 ->Threads(Users)->jp@gc - Ultimate Thread Group

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230820160724366-1550111674.png)

### 页面说明：

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230820160909500-963096260.png)

*   Start Threads Count：当前行的线程总数
*   Initial Delay/sec：延时启动当前行的线程，单位: 秒
*   Startup Time/sec：启动当前行所有线程达峰值所需时间，单位: 秒
*   Hold Load For/sec：当前行线程达到峰值后的稳定加载时间，单位: 秒
*   Shutdown Time：停止当前行所有线程所需时间，单位: 秒

Ultimate Thread Group 线程组灵活度还是相当高的，功能上，相当于把多个不同基础线程组进行组合。

 在 [jmeter 详解 - 线程组详解（3）- 再看 Ramp-Up(seconds)](https://www.cnblogs.com/jiushao-ing/p/17643095.html) 对线程组 Ultimate Thread Group 讲过一部分，这里按测试场景进行一下小结：

Ultimate Thread Group 可以用于以下几个场景的测试：

1.  **创建线性负载**
2.  **创建阶梯负载**
3.  **创建尖峰负载**
4.  **创建波浪形负载**

###  （1）创建线性负载

测试场景：60s 内启动 100 个线程，持续运行 60s，花 10s 的时间结束

脚本配置如下：

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230820161651539-1641025478.png)

添加监听器 jp@gc - Active Threads Over Time，运行后查看线程运行情况

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230820162111158-645431283.png)

### **（2）创建阶梯负载**

测试场景：

测试 100 个用户，我们将逐步地将它们提升。我们将从 25 个用户开始在一定时间内保持一个负载，查看服务器如何处理它。之后我们会再加 25 个到 50 个再加 25 个到 75 个，最后加 25 个到 100 个用户。这种方法效果好得多，也更可靠。

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230820162350457-1690730289.png)

通过以上配置，观察日志和监听器，就可以知道系统在哪个负载下面平稳运行，能承担多大的负载

添加监听器 jp@gc - Active Threads Over Time，运行后查看线程运行情况

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230820163535163-1110978385.png)

### （3） **创建尖峰负载**

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230820163844457-369533213.png)

### （4）**创建波浪形负载**

比如 12306 抢票的时候，每次开放抢票时，有大量用户涌入，等到下次开放时，又有大量用户涌入，这个时候，就像波浪一样，不断敲击服务器，考验服务器的性能

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230820164124197-2097501981.png)

锯齿形

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230820164625107-802839073.png)