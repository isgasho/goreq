---
title: README
---

类似java的OKHttp的简单易用的HTTP客户端。

## 设计

简单的三层：

- Client
  - 可复用的通用逻辑的封装，比如api的基础路径，添加api认证参数，超时，代理，tls，body的编解码等
  - 使用codec机制实现请求和响应内容的编解码，默认支持标准库的json、xml编解码，也可以自定义自己的视线方式
  - 使用插件机制实现在请求前后对request和response做二次处理，比如请求前携带auth，请求日志打印，错误日志发送到日志系统，熔断降级等，支持自定义插件开发
- Req：
  - request的封装，更方便的设置
  - 支持链式操作
- Resp：
  - response的封装，提供各种方便的获取响应内容的方法

## 应用场景

可以简单高效的解决如下常用使用场景：

- 简单的调用一个接口地址并获取返回值
- 调用api接口时需要补充认证相关的参数
- 上传/下载文件
- 调用api接口时，需要判断接口是否返回错误，并解析返回数据到struct
- 调用api接口时，只是想获取返回内容的某一个字段的值，并不想构建struct
- 调用时自动打印请求和响应内容
- 调用时自动将请求和响应内容发送到消息队列，保存调用日志
- 调用api失败时，使用默认内容作为响应内容。类似熔断降级的callback
- 请求的json可以设置是否EscapeHTML

## 安装

```shell
go get -u github.com/aiscrm/goreq
```

## 例子

### 简单例子

```go
goreq.NewClient().Use(log.Dump()).Get(ts.URL).Do().AsString()
```

## TODO
