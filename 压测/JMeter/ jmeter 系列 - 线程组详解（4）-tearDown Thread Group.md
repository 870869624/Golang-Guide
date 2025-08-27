> 本文由 [简悦 SimpRead](http://ksria.com/simpread/) 转码， 原文地址 [www.cnblogs.com](https://www.cnblogs.com/jiushao-ing/p/17643481.html)

tearDown Thread Group 线程组：
--------------------------

在测试任务线程组运行结束后被运行。通常用来做清理测试脏数据、登出、关闭资源等工作。

应用场景举例：

A、测试数据库操作功能时，用于执行关闭数据库连接的操作。  
B、测试用户购物功能时，用于执行用户的退出等操作。

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230820154835900-1992521068.png)![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230820154758444-1670642381.png)

tearDown Thread Group 线程的行为与普通的线程组元素完全一样。不同之处在于这些类型的线程在测试完成执行其常规线程组之后执行。

注意点：

请注意，默认情况下，如果 Test 正常关闭，它将不会运行，如果您想让它在这种情况下运行，请确保在测试计划元素上选中 “在主线程关闭后运行 tearDown 线程组” 选项。如果停止测试计划，即使选中选项，tearDown 也不会运行。

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230820155024531-308098949.png)

 关闭主线程后运行 tearDown Thread Groups