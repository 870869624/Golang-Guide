# 使用wrk进行性能测试

## 1 wrk介绍

wrk是一款现代化的HTTP性能测试工具，即使运行在单核CPU上也能产生显著的压力。它融合了一种多线程设计，并使用了一些可扩展事件通知机制，例如epoll and kqueue。一个可选的LuaJIT脚本能产生HTTP请求，响应处理和自定义报告，更详细的脚本内容可以参考scripts目录下的一些例子。

## 2 wrk下载和安装

先安装git

ubnutu:

```bash
$ sudo apt install git -y
```

centos:

```bash
$ sudo yum install git -y
```

下载wrk文件

```bash
$ git clone https://github.com/wg/wrk.git  
$ cd wrk  
$ make
```

编译成功后，目录下就会有一个wrk文件。如果编译过程中，出现如下错误：

若报错gcc: Command not found，则需安装[gcc](https://gitcode.com/Resource-Bundle-Collection/1e528?utm_source=highlight_word_gitcode&word=%E5%AE%89%E8%A3%85gcc&isLogin=1)，[参考wrk编译报错gcc: Command not found](https://www.cnblogs.com/ycyzharry/p/10383845.html)

若报错fatal error: openssl/ssl.h: No such file or directory，则需要安装openssl的库。

```bash
$ sudo apt install libssl-dev
```

或者

```bash
$ sudo yum install  openssl-devel
```

## 3 一个简单的例子

```bash
$ ./wrk -t12 -c400 -d30s http://localhost:8080/
```

它将会产生如下测试，12个线程（threads），保持400个HTTP连接（connections）开启，测试时间30秒(seconds)。

详细的命令行参数如下：

```bash
-c,    --connections（连接数）:      total number of HTTP connections to keep open with each thread handling N = connections/threads
 
-d,    --duration（测试持续时间）:     duration of the test, e.g. 2s, 2m, 2h
 
-t,    --threads（线程）:            total number of threads to use
 
-s,    --script（脚本）:             LuaJIT script, see SCRIPTING
 
-H,    --header（头信息）:           HTTP header to add to request, e.g. "User-Agent: wrk"
 
       --latency（响应信息）:         print detailed latency statistics
 
       --timeout（超时时间）:         record a timeout if a response is not received within this amount of time.
```

下面是输出结果：

```bash
Running 30s test @ http://localhost:8000/
  12 threads and 400 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   635.91us    0.89ms  12.92ms   93.69%
    Req/Sec    56.20k     8.07k   62.00k    86.54%
  22464657 requests in 30.00s, 17.76GB read
Requests/sec: 748868.53
Transfer/sec:    606.33MB
```

结果解读如下：

- Latency: 响应信息, 包括平均值, 标准偏差, 最大值, 正负一个标准差占比。 

- Req/Sec: 每个线程每秒钟的完成的请求数，同样有平均值，标准偏差，最大值，正负一个标准差占比。

- 30秒钟总共完成请求数为22464657，读取数据量为17.76GB。

- 线程总共平均每秒钟处理748868.53个请求（QPS）,每秒钟读取606.33MB数据量。

## 4 发送post请求例子

首先需要创建一个post.lua文件，内容如下：

```lua
wrk.method = "POST"
wrk.headers["uid"] = "127.0.0.1"
wrk.headers["Content-Type"] = "application/json"
wrk.body     ='{"uid":"127.0.0.1","Version":"1.0","devicetype":"web","port":"8080"}'
```

测试执行命令如下：

```bash
./wrk --latency -t100 -c1500  -d120s --timeout=15s -s post.lua http://127.0.0.1:8080/index.html
```

这个脚本加入了--lantency：输出结果里可以看到响应时间分布情况，

```bash
Running 2m test @ http://127.0.0.1:8080/index.html
  100 threads and 1500 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency   127.26ms  157.88ms   2.44s    87.94%
    Req/Sec   177.91     45.09     1.10k    69.97%
  Latency Distribution
     50%   66.05ms
     75%  143.57ms
     90%  325.41ms
     99%  760.20ms
  2138290 requests in 2.00m, 3.17GB read
Requests/sec:  17804.57
Transfer/sec:  27.05MB
```

## 5 带随机参数的get请求例子

如果想构造不同的get请求，请求带随机参数，则lua脚本如下：

```lua
request = function()
num = math.random(1000,9999)
   path = "/test.html?t=" .. num
   return wrk.format("GET", path)
end
```

## 6 添加参数txt文件的get请求例子

如果要测试的url需要参数化，uids.txt文件内容如下：

```txt
100
101
102
```

lua脚本如下：

```lua
urimap = {}
counter = 0
function init(args)
    for line in io.lines("uids.txt") do
       print(line)
       urimap[counter] = line
       counter = counter + 1
   end
   counter = 0
end
 
request = function()
   local path ="/GetInfo.aspx?u=%s&m=1"
   parms = urimap[counter%(table.getn(urimap) + 1)]
   path = string.format(path,parms)
   counter = counter + 1
   return wrk.format(nil, path) 
end
```

## 7 添加参数txt文件的post请求例子

lua脚本如下：

```lua
urimap = {}
counter = 0
function init(args)
    for line in io.lines("uids.txt") do
       urimap[counter] = line
       counter = counter + 1
   end
   counter = 0
end
 
request = function()
   local body1 = '{"uid":"100%s'
   local body2 = '","name":"1"}'
   parms = urimap[counter%(table.getn(urimap) + 1)]
   path = "/getinfo"
   method = "POST"
   wrk.headers["Content-Type"] = "application/json"
   body = string.format(body1,parms)..body2
   counter = counter + 1
   return wrk.format(method, path, wrk.headers, body)
end
```

若参数txt中有转义字符，可用如下方法处理：

```lua
parms = string.gsub(urimap[counter%(table.getn(urimap) + 1)],'\r','')
```

如果要打印返回数据，可添加如下脚本：

```lua 
a=1
function response(status, headers, body)
   if(a==1)
   then
           a=2
          print(body)
      end
end
```

## 8 提交不同表单内容例子

lua脚本如下：

```lua 
wrk.method = "POST"
wrk.body = ""
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
 
-- 提交不同表单内容
local queries = {
    "version=1.0",
    "version=2.0",
    "version=3.0"
}
local i = 0
request = function()
    local body = wrk.format(nil, nil, nil, queries[i % #queries + 1])
    i = i + 1
    return body
end
```

## 9 访问多个url例子

如果要随机测试多个url，可参考[wrk-scripts](https://github.com/timotta/wrk-scripts)这个项目。

需要创建一个文件名为paths.txt，里面每行是一个要测试的url网址。lua脚本如下:

```lua 
counter = 0
 
-- Initialize the pseudo random number generator - http://lua-users.org/wiki/MathLibraryTutorial
math.randomseed(os.time())
math.random(); math.random(); math.random()
 
function file_exists(file)
  local f = io.open(file, "rb")
  if f then f:close() end
  return f ~= nil
end
 
function shuffle(paths)
  local j, k
  local n = #paths
  for i = 1, n do
    j, k = math.random(n), math.random(n)
    paths[j], paths[k] = paths[k], paths[j]
  end
  return paths
end
 
function non_empty_lines_from(file)
  if not file_exists(file) then return {} end
  lines = {}
  for line in io.lines(file) do
    if not (line == '') then
      lines[#lines + 1] = line
    end
  end
  return shuffle(lines)
end
 
paths = non_empty_lines_from("paths.txt")
 
if #paths <= 0 then
  print("multiplepaths: No paths found. You have to create a file paths.txt with one path per line")
  os.exit()
end
 
print("multiplepaths: Found " .. #paths .. " paths")
 
request = function()
    path = paths[counter]
    counter = counter + 1
    if counter > #paths then
      counter = 0
    end
    return wrk.format(nil, path)
end
```