> 本文由 [简悦 SimpRead](http://ksria.com/simpread/) 转码， 原文地址 [www.cnblogs.com](https://www.cnblogs.com/jiushao-ing/p/17644388.html)

bzm - Concurrency Thread Group
------------------------------

### 介绍：

Concurrency Thread Group 中文翻译就是并发线程组。此线程组提供了配置线程调度的简化方法。它旨在维护并发级别，这意味着如果没有足够的线程并行运行，则在运行时启动额外的线程。与标准 Thread Group 不同，它不会预先创建所有线程，因此不会使用额外的内存。对于 Stepping Thread Group 线程组来说，这是一个很好的改进，因为它允许线程优雅地完成它们的工作。可以看做是 Stepping Thread Group 的一个替代。

### 添加方式：

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230820184245637-504513142.png)

### 页面说明：

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230821162419751-1629277449.png)

*   **Target Concurrency**：目标并发（线程数）
*   **Ramp Up Time**：启动时间；若设置 1 min，则目标线程在 1 imn 内全部启动
*   **Ramp-Up Steps Count**：阶梯次数；若设置 3 ，则目标线程在 1min 内分六次阶梯加压（启动线程）；**每次启动的线程数** = 目标线程数 / 阶梯次数 = 12 / 3 = 4
*   **Hold Target Rate Time**：持续负载运行时间；若设置 2 ，则启动完所有线程后，持续负载运行 2 min，然后再结束
*   **Time Unit**：时间单位（分钟或者秒）
*   **Thread Iterations Limit：**线程迭代次数限制（循环次数）；默认为空，理解成永远，如果**运行时间到达** Ramp Up Time + Hold Target Rate Time，则停止运行线程**【不建议设置该值】**
*   **Log Threads Status into File：**将线程状态记录到文件中（将线程启动和线程停止事件保存为日志文件）；

**注意点：**

1.  Target Concurrency 只是个**期望值**，实际不一定可以达到这个并发数，还得看配置如电脑性能、网络、内存、CPU 等因素都会影响最终并发线程数
2.  Jmeter 会根据 Target Concurrency 的值和当前处于**活动状态的线程数**来判断当前并发线程数是否达到了 Target Concurrency；若没有，则会不断启动线程，尽力让并发线程数达到 Target Concurrency 的值

添加监听器 Active Threads Over Time，可以看到活跃线程数维持在了 12

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230821163650432-1087344492.png)

**注意点：**

虽然从 Concurrency Thread Group 负载预览图和 Active Threads Over Time 看，每次阶梯增压都是瞬时增压的，但是实际测试结果可以看到它也是有一个过渡期，并不是瞬  
增压。上图 Active Threads Over Time 看起来是瞬间启动的原因是线程数太少，对机器没有压力直接就瞬发了，如果线程数过大，是可以看到一个线程数爬坡过程的。

实际上若将 Ramp-Up Steps Count 设置为 0，就成了线性负载

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230821171225039-1411430770.png)

**小结：**

Concurrency Thread Group 从一开始的启动线程数是尽力去达到这个启动线程数（第一个阶梯的线程数），到了阶梯的最上面（最后一个阶梯）也是尽量去保持这个并发数（即便有波动时也会尽量保持目标并发），线程最后释放的时候也是尽力去达到一个瞬时释放。

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230821170253735-595049686.png)

### Concurrency Thread Group 和 Stepping Thread Group 的区别

*   Concurrency Thread Group 不能设置启动延迟、Ramp-Up 时间等，而 Stepping Thread Group 对阶梯加压的过程控制得更为细腻
*   Concurrency Thread Group 到了结束时间点瞬间释放线程，而 Stepping Thread Group 可以设置阶梯释放时间
*   Concurrency Thread Group 会**尽力启动线程达到** Target Concurrency 值，而 Stepping Thread Group 设置了多少个线程就会**严格执行**

### 与 Throughput Shaping Timer（吞吐量计时器）一起使用

*   当 Concurrency Thread Group 与 Throughput Shaping Timer（吞吐量计时器）一起使用时，可以用 tstFeedback 函数的调用来替换 Target 并发值，以动态维护实现目标 RPS 所需的线程数，使用此方法时， 需要将 Ramp Up Time 和 Ramp-Up Steps Count 置空，但要确保 Hold Target Rate Time ≥ Throughput Shaping Timer 时间表中指定的总持续时间值（Duration）
*   有关功能使用详细信息，请参阅 Throughput Shaping Time 计划反馈功能。