# 简单的高性能 Golang MVC 框架

这个框架是我在用golang重构一个php项目时写的，因为被重构的项目不大，感觉用别的框架显得有点臃肿，所以干脆自己写一个。

这个框架之所以称之为高性能，一是因为结构简单，没有多余的东西，不臃肿；二是引用到的第三方包，都是精选同类型里面性能最高的那些。

## 主要功能

简单的路由调度

HTTP服务支持优雅退出、平滑重启

TLS支持绑定多个域名证书

支持HTTP BasicAuth

全局异常处理

支持MySQL、Redis

Redis消息队列处理

JSON配置文件

CLI方法封装

日志封装

控制器封装

数据模型封装

## 主要用到的第三方包

fasthttp（高性能HTTP服务，比原生快10倍）

jsoniter（高性能 JSON解析包，比原生快6倍）

hero（高性能模板引擎，支持模板预编译，所有模板引擎里面性能最高）

gorm（数据库ORM包，支持连接池、缓存预编译SQL语句，性能很好）

redigo（Redis客户端，支持连接池，高并发下性能很好）

## 使用说明

### 安装依赖
```
go mod init main
go get
```

### 运行

先复制`config.bak.json`配置文件，并改名成`config.json`，里面只保留需要的内容。

然后执行以下命令

```
go run .
```

可以访问 127.0.0.1:8888，如果启用TLS的话，端口是443

### 控制器

在`controller`目录下面创建文件

```go
package controller

type IndexController struct{}

var Index IndexController

func (t *IndexController) Default(ctx *fasthttp.RequestCtx) {
	//方法内容
}
```

### 路由

在`main.go`文件里面的`index`方法中添加

```go
var c = map[string]func(ctx *fasthttp.RequestCtx){
	"/": controller.Index.Default,
}
```

### 模型

在`model`目录下创建文件

```go
package model

type SampleModel struct{}

var Sample SampleModel

func (t *SampleModel) GetList() (r string) {
  r = "Hello Go ~"
  return
}
```

调用模型方法

```go
r := model.Sample.GetList()
```

全局可使用`base.DB`操作数据库，具体可以参考 [GORM文档](https://gorm.io/zh_CN/docs/)

### 视图

在`veiw`目录下对应的控制器和方法下面创建文件，模板引擎使用方法可以参考 [hero文档](https://github.com/shiyanhui/hero/blob/master/README_CN.md)

### HTTP服务优雅退出与平滑重启

优雅退出可以使用`Ctrl+C`或者执行以下命令

```
kill 进程号
```

平滑重启执行以下命令

```
kill -USR2 进程号
```

### 日志

全局可使用`base.Log()`方法记录日志，文件存放在`log`目录下面。

### CLI操作

在`cli.go`文件里面添加方法，然后在`main.go`文件的入口方法添加CLI的路由

```go
switch cmd[1] {
case "cli":
	var f = map[string]func(){
		"Test": Test,//CLI方法
	}
	f[cmd[2]]()
	break
}
```

通过以下命令执行

```
go run . cli 方法名
```

### Redis 操作

全局可使用`base.Redis`操作连接池，简单读取字符串的话，可以使用`base.RedisGet(key)`

```go
re := base.Redis.Get()//在连接池中取出一个链接
defer re.Close()//放回链接池
re.Do("SET", "key", "value", "EX", 3600)
```

### Redis消息队列处理

可以使用以下方法推数据到队列，两个参数分别是队列名称和数据

```go
base.PushQueue("test", "data")
```

创建队列相应的处理方法，在`queue`目录下创建文件

```go
package queue

type TestQueue struct{}

var Test TestQueue

func (t *TestQueue) Exec(d string) {
	//方法内容
}
```

加载队列处理方法，在`base.go`文件中添加配置，队列名称和相应的处理方法

```go
queueFunc = map[string]func(string){
	"test": queue.Test.Exec,
}
```

