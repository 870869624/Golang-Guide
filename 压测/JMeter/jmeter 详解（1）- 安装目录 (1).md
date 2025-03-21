> 本文由 [简悦 SimpRead](http://ksria.com/simpread/) 转码， 原文地址 [www.cnblogs.com](https://www.cnblogs.com/jiushao-ing/p/17628366.html)

JMeter 官网：

https://jmeter.apache.org/

GitHub:

https://github.com/apache/jmeter

用户文档（英文）：

https://jmeter.apache.org/usermanual/index.html

 ApacheJMeter 可用于测试静态和动态资源、Web 动态应用程序的性能。它可以用来模拟一台服务器、一组服务器、网络或对象上的重负载，以测试其强度或分析不同负载类型下的整体性能。

ApacheJMeter 的功能包括：

（1）能够对许多不同的应用程序 / 服务器 / 协议类型进行负载和性能测试:

*   Web - HTTP, HTTPS (Java, NodeJS, PHP, ASP.NET, …)
*   SOAP / REST Webservices
*   FTP
*   Database via JDBC
*   LDAP
*   Message-oriented middleware (MOM) via JMS
*   Mail - SMTP(S), POP3(S) and IMAP(S)
*   Native commands or shell scripts
*   TCP
*   Java Objects

（2）全功能的测试 IDE，允许快速测试计划记录 (从浏览器或本地应用程序)，构建和调试

（3）CLI 模式 (命令行模式 (以前称为非 GU) /headless 模式)，用于从任 Java 兼容的操作系统 (Linux,Windows,Mac OSX,..

（4）一个完整且随时可以呈现动态 HTML 报告

（5）能从最流行的响应格式中提取数据并进行关联，如： **[HTML](https://jmeter.apache.org/usermanual/component_reference.html#CSS/JQuery_Extractor), [JSON](https://jmeter.apache.org/usermanual/component_reference.html#JSON_Extractor) , [XML](https://jmeter.apache.org/usermanual/component_reference.html#XPath_Extractor) or [any textual format](https://jmeter.apache.org/usermanual/component_reference.html#Regular_Expression_Extractor)**

（6）完全的可移植性和 100% 的 Java 代码编写的纯应用

（7）完整的多线程框架允许多个线程的并发采样，以及通过单独的线程组同时采样不同的功能。

（8）缓存和离线分析 / 重放测试结果。

（9）高度可扩展的核心:

*   可插拔采样器允许无限的测试能力。
*   脚本采样器 (兼容 JSR223 语言，如 Groovy 和 BeanShell)
*   可插拔定时器可以选择几个负载统计信息。
*   数据分析和可视化插件允许很大的扩展性以及个性化。
*   函数可用于为测试提供动态输入或提供数据操作。
*   通过面向 Maven、Gradle 和 Jenkins 的第三方开源库轻松实现持续集成

 **安装目录详解：**

版本：apache-jmeter-5.5

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230814124536615-1438857554.png)

**安装目录说明：**

*   **bin**：包含启动、配置等相关命令
*   docs：官方接口文档，二次开发需要
*   extras：辅助库，持续集成会用到
*   **lib**：核心库，包含 JMeter 用到的各种基础库和插件依赖的 jar 包或存放自己二次开发的 jar 包；lib/ext 文件夹：第三方插件、Jmeter 二进制文件，
*   license：包含 non-ASF 软件的许可证
*   printable_docs：离线的帮助文档，可以查看函数等内容
*   LICENSE：JMeter 许可说明
*   NOTICE：JMeter 简单信息说明
*   README.md：JMeter 官方基本介绍

 **printable_docs**

函数和变量见：functions.html

取样器见：component_reference.html，包括取样器、逻辑控制器、监听器、环境变量、断言、计时器、前置处理器、后置处理器、其他功能

远程测试：remote-test.html

**关于 bin 目录**

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230814151947944-825306536.png)

*   **jmeter.properties：**JMeter 核心配置文件，各种配置基本在这完成
*   **log4j.conf：**JMeter 日志配置管理
*   **jmeter.log：**JMeter 运行日志记录，什么输出信息、警告、报错都在这里进行了记录
*   **jmeter.bat：**windows 下 jmeter 启动文件
*   **jmeterw.cmd：**windows 下 jmeter 的启动文件，不带 cmd 窗口
*   **shutdown.cmd：**windows 下 jmeter 关闭文件
*   **stoptest.cmd**：windows 下 jmeter 测试停止文件
*   **jmeter-server.bat**：windows 下 jmeter 服务器模式启动文件
*   **jmeter-server：**mac 或者 Liunx 分布式压测使用的启动文件

**jmeter.properties**

配置项的说明在目录下：

printable_docs/usermanual/properties_reference.html

![](https://img2023.cnblogs.com/blog/2565457/202308/2565457-20230814180507590-222643719.png)

 _**最佳实践（注意事项）：**_

_网址：https://jmeter.apache.org/usermanual/best-practices.html_

（1）最好使用最新的 JMeter 版本

（2）使用正确的线程数

影响线程数的因素：

*   硬件（自己的电脑等用于压测的机器）
*   测试用例（test plan）的设计
*   服务器响应速度（服务器响应速度快，返回给 JMeter 的响应速度就快，JMeter 就需要花时间处理）

a）与任何负载测试工具一样，如果您没有正确地确定线程的数量，您将面临 “协调遗漏” 问题，这可能会给您错误或不准确的结果。

b）如果需要大规模的负载测试，可以考虑使用 (或不使用) 分布式模式在多台机器上运行多个 CLIJMeter 实例。当使用分布式模式时，result 文件在 Controller 节点上组合。如果使用多自治实例，可以将示例结果文件组合起来进行后续的定量分析。

（3）添加 Cookie 管理器

要添加 cookie 支持，只需 在测试计划中的每个线程组中添加一个 HTTP Cookie 管理器。这将确保每个线程都有自己的 cookie，但在所有 HTTP 请求对象之间共享

（4）请求头管理

HTTP Header Manager 允许您自定义 JMeter 在 HTTP 请求标头中发送的信息。此标头包括 “User-Agent”、“Pragma”、“Referer” 等属性。

HTTP Header Manager 和 HTTP Cookie Manager 一样，应该在线程组级别添加，除非出于某种原因，您希望在测试中为不同的 HTTP Reques[t](https://jmeter.net/usermanual/component_reference.html#HTTP_Request) 对象指定不同的标头。

（5）减少资源使用

关于减少资源使用的一些建议：

*   使用 CLI 模式：jmeter -n -t test.jmx -l test.jtl
*   使用尽可能少的 Listeners；如果使用上面的 -l 标志，它们都可以被删除或禁用。
*   不要在负载测试期间使用 “查看结果树” 或“在表中查看结果”侦听器，仅在脚本编写阶段使用它们来调试脚本。
*   与其使用大量相似的采样器，不如在循环中使用相同的采样器，并使用变量（CSV 数据集）来改变样本。[包含控制器在这里没有帮助，因为它将文件中的所有测试元素添加到测试计划中]
*   不要使用功能模式（functional mode）
*   使用 CSV 输出而不是 XML
*   只保存您需要的数据
*   使用尽可能少的断言
*   使用性能最高的脚本语言（参见 JSR223 部分）
*   不要忘了删除的本地路径设置配置如果使用 CSV 数据。
*   每次测试运行前清理 “文件” 选项卡。

如果您的测试需要大量数据 - 特别是如果需要随机化 - 在可以使用 CSV 数据集读取的文件中创建测试数据。这避免了在运行时浪费资源。