> 本文由 [简悦 SimpRead](http://ksria.com/simpread/) 转码， 原文地址 [blog.csdn.net](https://blog.csdn.net/beginningL/article/details/144116706)

### 一. 普通性能场景设计

![](https://i-blog.csdnimg.cn/direct/8345c37030384798981c124273b69f85.png)

**线程数**：jmeter 本身对线程数没有限制, 但是 jmeter 钱这些并发用户数的时候，需要消耗资源，收到 电脑 [spu](https://so.csdn.net/so/search?q=spu&spm=1001.2101.3001.7020) 的主频限制，一条电脑不可能无限制的启动线程数，实际情况下，http 协议大概可以产生 1500 左右，保守一点是 1000 左右，如果要模拟超过 1000 的并发线程数，可能需要考虑分布式。  
**ramp-up**：启动所有线程数的时间，500 以内的并发用户，ramp-up 大概是 2-3 秒；500-1000 ramp-up 大概是 5 秒；大于 1000 ramp-up 大概是 5-8 秒；一个原则：ramp-up 时间在总的执行时间中占比要很低。  
**循环次数**：默认必须大于 1。  
一个简单的例子：并发用户数是 40，运行 1 分钟，看到吞吐率没有网络瓶颈，所以聚合报告中的吞吐量等于 tps

![](https://i-blog.csdnimg.cn/direct/8a6afbd8c0ea4b8da83d597941801166.png)tps 是 9.7 左右，并发用户数是 30，吞吐量小于用户数，说明每秒不能处理一个请求，这个接口的 tps 较低；  
90% 的响应时间是 3.035 秒已经超过用户满意度指数 1.5s 所 i 响应时间有点慢；  
由此得出此用户的并发用户数小于 30

### 二. 阶梯性能场景设计（负载场景设计）

**应用场景**：如果接口从来没有做过性能测试，不知道一些指标信息，可以使用负载场景测试逐渐增加并发用户数来找出性能瓶颈。逐步增加并发用户数，是缓起步快结束，并不是瞬间结束  
**测试前的准备**：使用 jmeter 中的插件安装 jp@gc 插件的安装并重启 jmeter  
**负载测试的执行**：添加线程组–jp@gc Stepping Thread Group  
![](https://i-blog.csdnimg.cn/direct/751fe75d2ebd46358a7307f893ed4092.png)  
这个配置的测试场景是从 0 线程组进行增加，每次增加 1 个线程组，达到一个线程组就持续运行 30 秒时间，然后知道达到最大的线程数 10，持续运行 60s 时间，最后在 5s 内将线程停掉看测试结果  
添加监听器有三个, 三个图要结合一起看：  
jp@gc -Active Threads Over Time: 查看并发用户线程的变化  
jp@gc -Transactions per Second: 查看 tps 的变化趋势图  
jp@gc -Response Times Over Time: 查看响应时间的变化  
![](https://i-blog.csdnimg.cn/direct/e02082f9b4364aaa81b982dda50a523b.png)  
![](https://i-blog.csdnimg.cn/direct/d15878cd1982459ebb19383d93db9c42.png)![](https://i-blog.csdnimg.cn/direct/f000203a7b9d491ba7306d2e18b1e6c3.png)  
三个图结合一起看下，可以看到并发数是 10，tps 发改是在 18 左右，响应时间 0.5s 左右，所以通过上面三张图可以看到并发数 10 是完全可以达到的，tps 也大于 18。  
注意：  
（1）jp@gc -Transactions per Second 这个图表可以看到 tps 的变化趋势，也可以看到失败的数据  
（2）jp@gc -Response Times Over Time 这个图可以看到响应时间的变化，如果响应时间是一条直线的时候此时大部分是因为接口超时了。  
（3）负载测试是不可以看聚合报告的，因为聚合报告中的线程数要一样，所以阶梯性能长江不能看聚合报告

### 三. 压力测试场景设计

**应用场景**：通过负载测试找到对应的性能的最大并发用户数，然后设置最大用户数的 20% 和 80% 进行长时间的测试。  
**执行测试**：普通线程组设置持续时间即可，持续时间设置久一点。  
![](https://i-blog.csdnimg.cn/direct/d1c1f016cffe404697f9125cd2672fe8.png)

### 四. 面向目标场景设计

**应用场景**：期望项目的接口达到多少 tps，或者期望项目的接口达到多少并发用户数，这种即是面向目标场景的性能测试场景。  
**场景一**：面向多少 tps：期望我项目的接口都要满足 50tps，中小企业如果日均访问量是千万，基本 50tps 就可以满足了。  
**执行测试**：添加线程组–bzm-Arrivals Tjread Group  
![](https://i-blog.csdnimg.cn/direct/442dda38f01d4f7c97e2c40d01a0e078.png)  
Target Rate(arrivals/sec)：目标的 tps 值  
Ramp Up Timde(sec)：达到这个目标 tps 的总时间  
Ramp-Up steps Count：分几次达到这个目标值  
Hold Target Rate Time：达到目标值以后持续运行多久  
![](https://i-blog.csdnimg.cn/direct/2e63e89eccc349b7a87d8c7137753c4b.png)  
![](https://i-blog.csdnimg.cn/direct/1a6f88b92bd443eca2d4b88d1a4b98fc.png)  
![](https://i-blog.csdnimg.cn/direct/03c2c66a748c4f4999759eda39c44df9.png)  
分析：从上面三张图可以看出来，tps 达到 50tps 以后，并发用户数再 160 左右，但是响应时间是一条直线并且时间都大于 3s，所以这个接口的性能达不到 50tps。  
举个例子：要做一个秒杀，能支持 1000 个人同时秒杀，我们的系统不能崩溃怎么设计性能场景？？1000 个人访问我们系统持续运行，系统不能崩溃，用户对秒杀的理解，我要在 1 秒钟内收到处理结果 1000tps 转换为面向目标的场景。  
**场景二：面向多少并发用户数 / 线程数**  
执行测试：添加线程组 - bzm-Concurrency Thread Group  
![](https://i-blog.csdnimg.cn/direct/95d62260fc0447238b6b6dde8bd8fe63.png)  
Target Concurrency：目标的并发用户数  
Ramp Up Timde(sec)：达到这个目标的总时间  
Ramp-Up steps Count：分几次达到这个目标值  
Hold Target Rate Time：达到目标值以后持续运行多久  
![](https://i-blog.csdnimg.cn/direct/33b3c24f25fc49b8bf35a6268c901135.png)  
![](https://i-blog.csdnimg.cn/direct/c3f68976d9c747568c8b123533faf25c.png)  
![](https://i-blog.csdnimg.cn/direct/3943045ab36d4d47917cbfb94e7d0685.png)  
测试结论：并发用户数设置 80 以后，tps 大概是 60 左右，小于并发用户数，并且响应时间在 3s 左右近似一条直线，可以判断接口超时了，所以这个接口的并发用户数达不到 80。

### 五. 有时间规律场景（波浪型场景）

**应用场景**：比如点餐系统饭点会比较频繁调用，打卡的时候指定时间点会频繁使用

**执行测试**：添加线程组–Ultimate Thread Group  
![](https://i-blog.csdnimg.cn/direct/de746c26533e4d27bf1a910d39fdc6e1.png)start_threads count：达到的总的线程数

Initial Delay：起始的时间（如果是第二行新增的起始时间应该是上一行所有时间之和）

Startup Time：达到线程数的总时间

Hold Load For：持续运行的时间

### 六. 混合场景测试

应用场景：接口 1 线程数是 10，接口 2 线程数是 20，接口 3 线程数是 30，并且接口之间有关联，分别设置三个线程组，每个线程组的线程数不一样，三个线程组之间传参可以跨线程组传参  
执行测试：有两个接口一个是登录接口并发用户数是 80，另外一个是商品查询 6 合一接口并发用户数是 18 怎么设计性能测试场景？？  
分析：商品 6 合一接口要使用到登录接口的 token，所以要将 token 设置成 jmeter 的属性进行跨线程组进行调用。  
（1）编写登录的接口，并且根据响应结果使用 json 提取器将对应的 token 进行提取出来  
![](https://i-blog.csdnimg.cn/direct/a9a4cef04db047b2a30645b96b9c1c8d.png)  
（2）使用函数助手将提取的 token 应用 setProperty 函数给设置成 jmeter 的属性

![](https://i-blog.csdnimg.cn/direct/60b885c682764c7e81d072dda0b85d91.png)  
![](https://i-blog.csdnimg.cn/direct/7729f0b414354547ac03528c55f58ca5.png)（3）将这个属性通过 P 函数将对应的值取出来用于下一个线程组的接口中  
![](https://i-blog.csdnimg.cn/direct/cd25ea46a488479baeee76f390519de0.png)  
（4）将登录接口的线程数设置成 80，将商品 6 合一接口的线程数设置成 18，并且添加监听器进行查看

（5）执行对应的性能测试

```
登录接口的测试结果图如下

```

![](https://i-blog.csdnimg.cn/direct/539156a12a774153be1e49987fa02b18.png)  
![](https://i-blog.csdnimg.cn/direct/3c0ae89eaa2a47ba8108dcf30d7bf5c5.png)  
![](https://i-blog.csdnimg.cn/direct/cffd640f2061492587516d0a7e402e84.png)  
性能测试结论：响应时间在 1s 左右，对应的 tps 大概在 85 左右，所以这个接口的并发数可以达到 80，tps 可以达到 85，还可以继续进行往上施压进行测试。

商品 6 合一接口对应的性能测试结果如下图：  
![](https://i-blog.csdnimg.cn/direct/7a627e4af246423196241a95d572f661.png)  
![](https://i-blog.csdnimg.cn/direct/13ac9d253317436bbbc61c068150cf28.png)  
![](https://i-blog.csdnimg.cn/direct/a867405f0e4948cda4b7d1ff135c2d66.png)  
测试结论：并发数是 18，对应的响应时间是 1.2s 左右。tps 在 20 左右，所以这个接口的性能是可以达到 18 并发数，tps 可以达到 20，性能还可以继续往上进行施压得到更好的数据。

### 七. 拓展知识

```
   在性能测试过程中，能不启动监听器则不启动，但是没有监听器我们怎么知道性能测试结果？？

```

（1）生成 jmeter 的 html 报告，与是否启用添加监听器无关

（2）添加结果数将测试结果存储在文件中  
![](https://i-blog.csdnimg.cn/direct/0b7b99f253614210a59692e421cc03f9.png)  
（3）运行完成以后，可以使用 jmeter 中的工具–Generate html report 功能将刚刚生成的结果文件保存成 html 文件格式进行打开

![](https://i-blog.csdnimg.cn/direct/950a4317c26f4fb880f4f66183a5799c.png)  
![](https://i-blog.csdnimg.cn/direct/c5a80809d78e450295629e88526eafac.png)

