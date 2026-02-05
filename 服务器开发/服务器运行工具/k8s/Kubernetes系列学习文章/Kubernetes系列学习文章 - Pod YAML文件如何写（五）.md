# 一、什么是YAML文件

前面我们了解到K8S配置文件都是走YAML文件格式的，那么什么是YAML？它的编写语法是什么？

## 1. YAML特点

YAML 的英文全称是：”Yet Another Markup Language”（仍是一种置标语言）。我们之前接触过properties、XML、json等数据格式，越到后面越好用。其实YAML就是结合了这些标记语言的特性，整合新开发的。

大体来讲，YAML有下面特点：

层次分明、结构清晰
使用简单、上手容易
表达强大、语义丰富
但是要注意的是，下面几点：

大小写敏感
禁止使用tab键缩进，只能空格键

## 2. YAML语法

我们举个模版例子就能快速了解YAML的语法：

```yaml
# 前面是key，后面是value，表示如下：
name: nginx

# 表示metadata.name=nginx:
metadata:
  name: nginx
# 表达数组，即表示containers为[name,image,port] 那么用 - 就可以表示了

containers:
  - nginx
  - nginx
  - 80

# 常量、布尔、字符串定义
version: 1.1 # 定义一个数值1.1
rich: true # 定义一个boolean值
say: "hello world" # 定义一个字符串
```

掌握了上面的语法，基本上K8S的 YAML你就能看懂和编写了。K8S里的YAML几乎用不到什么高级的其他语法格式。

详细想了解全的话，可以看看这个文章：[https://learnxinyminutes.com/docs/yaml/](https://learnxinyminutes.com/docs/yaml/)

# 二、Pod YAML参数定义

Pod是K8S的最小单元，它的信息都记录在了一个YAML文件里。那么这个YAML文件到底怎么写呢？里面有哪些参数？如何去修改YAML文件？带着这几个问题我们来了解下。

Pod YAML有哪些参数？

K8S的YAML配置文件我们初学者看了后都觉得很长，然后也觉得没什么规律。其实，我们可以梳理下从两个方面去了解。第一个是哪些是必写项，第二个是YAML包含哪些主要参数对象。

## 1. 哪些是必写项

注意，一个YAML文件，下面几个参数是必须要声明，不然绝对会出错：

| 参数名                  | 字段类型 | 说明                                                                            |
| ----------------------- | -------- | ------------------------------------------------------------------------------- | 
| version                 | String   | 这里是指的是K8S API的版本，目前基本上是v1，可以用`kubectl api-versions`命令查询 |
| kind                    | String   | 这里指的是yaml文件定义的资源类型和角色，比如：Pod                               |
| metadata                | Object   | 元数据对象，固定值就写metadata                                                  |     
| metadata.name           | String   | 元数据对象的名字，这里由我们编写，比如命名Pod的名字                             |
| metadata.namespace      | String   | 元数据对象的命名空间，由我们自身定义                                            |
| Spec                    | Object   | 详细定义对象，固定值就写Spec                                                    |
| spec.containers[]       | list     | 这里是Spec对象的容器列表定义，是个列表                                          |
| spec.containers[].name  | String   | 这里定义容器的名字                                                              |
| spec.containers[].image | String   | 这里定义要用到的镜像名称                                                        |

以上这些都是编写一个YAML文件的必写项，一个最基本的YAML文件就包含它们。

## 2. 主要参数对象
第一小点里讲的都是必选参数，那么还是否有其他参数呢？其他功能的参数，虽然不是必选项，但是为了让YAML定义得更详细、功能更丰富，这里其他参数也需要了解下。接下来的参数都是Spec对象下面的，主要分了两大块：`spec.containers` 和 `spec.volumes`。

`spec.containers`

`spec.containers` 是个list数组，很明显，它代表的是描述container容器方面的参数。所以它下面的参数是非常多的，具体参数看如下表格：

[原文](https://cloud.tencent.com/developer/article/1478634)