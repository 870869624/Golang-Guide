> 本文由 [简悦 SimpRead](http://ksria.com/simpread/) 转码， 原文地址 [www.cnblogs.com](https://www.cnblogs.com/jiushao-ing/p/17638824.html)

setUp Thread Group
------------------

有的时候对于 python 的类来说，我们有前置后置动作。jmeter 的线程组也有前置线程组和后置线程组。执行优先级如下：

setup thread group  > thread group > teardown thread group

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230819203642930-135581222.png)![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230819203704271-417902493.png)

可以看得到，setUp Thread Group 和常规 Thread Group 在配置选项方面几乎是相同的。唯一的区别是 setUp Thread Group 将在其他线程组执行之前执

setup thread group 这个线程组主要用于来运行预负载测试操作。

在许多测试情况下，需要为负载测试准备目标测试环境。例如：

*   在数据库中创建一个注册用户列表
*   从数据库中检索所需要的数据，并在负载测试期间使用
*   在负载测试开始之前用测试数据填充数据库
*   运行应与负载测试分开执行的预负载测试计算。例如，一些参数计算一次，然后在负载测试期使用

来看一个例子：

我们定义以下场景:

我们要测试 10 个登录到我们的网站并购买一些衣服的用户。我们需要在执行测试之前创建用户。这可以通过以下方式使用 setUp Thread Group 来完成:

（1）配置数据库  
创建一个名为 “app-sign-ups” 表，它配置了 3 个字段“first-name”, "last-name”, “email”，表中内容为空

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230817232837369-1061753208.png)

（2）配置 setUp Thread Group 线程组

然后，创建一个包含 setUp Thread Group 的测试计划。

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230817233043481-211779383.png)

在 setUp 线程组中包含一个 CSV 数据文件配置，配置我们要注册的用户列表 (名字姓氏和电子邮件地址) 的 CSV 文件

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230817233341838-245676693.png)

然后在我们的请求体中将 “firstName”“lastName” 和 emailAddress”进行参数化。这个时候 JMeter 将解析 CSV 文件，并用与 CSV 文件数据相对应的值填充这些变量。

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230817233539161-304241733.png)

请求接口如下：

POST [https://sampledb-d52b.restdb.io/rest/app-sign-ups](https://sampledb-d52b.restdb.io/rest/app-sign-ups) (data no longer supported)

消息体中要发送的数据如下：

注意，我们需要确保使用与 CSV 中定义的变量名和数据库中定义的字段相同的变量名

然后我们还有添加一个 HTTP 头管理器，将 content-type 指定为 "application/json"，还有 HTTP cookie 管理器

```
{
"first-name":"${firstName}",
"last-name":"${lastName}",
"email":"${emailAddress}"
}

```

（3）配置基础线程组 Thread Group

接下来，添加包含实际测试逻辑和采样器 的常规 Thread Group:

1.  登录
2.  购买
3.  结账
4.  退出

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230817234203576-1646050727.png)

 到这里脚本就配置完毕，当我们运行这个脚本后，第一个运行的线程组是 setUp Thread Group 线程组。因为我们配置了 10 个线程与一个迭代，它将创建 10 个用户。当 setUp Thread Group 线程组完成后，另一个线程组 Thread Group 将开始执行，这就开始了我们的负载测试

 在下面的截图中，我们可以看到我们运行的结果，发送到我们数据库的请求是成功的，我们收到了一个 200 响应代码。

请注意，无论 setUp Thread Group 在 GUI 中的显示顺序如何它都将首先执行。

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230817234449608-1407238237.png)

 JSON 响应指示请求成功:

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230818000826221-752632360.png)

 回到数据库，检查一开始的那张表，看它是否创建了 10 个新用户。

可以看到，CSV 文件中的值填充到了表 "app-sign-ups" 中

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230817234630467-358332558.png)

 小结：

本来登录网站是没有用户的，用不存在的用户登录网站会报错，这里我们用 setUp Thread Group 线程组在数据库里面先新建了 10 个用户，然后在基础线程组 Thread Group 用这 10 个用户进行登录购买结算退出等操作，甚至我们可以再用一个 teardown thread group 线程组来清除数据库的用户。