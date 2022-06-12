# simple-demo

## 抖音项目服务端简单示例

具体功能内容参考飞书说明文档

工程无其他依赖，直接编译运行即可

```shell
go build && ./simple-demo
```

## 410 到底对不队 抖音项目



注意事项：

> 需要安装[FFmpeg](http://ffmpeg.org/download.html),   并修改controller/publish.go  row: 47 的ffmpeg路径
>
> 需要安装MySQL，并创建golang_mysql 数据库， 并修改dbmodel.go 文件中连接mysql的账户和密码
>
>
> 需要修改BaseURL.
>
> 需要自己新建 public 文件夹
>
> *新增 需要安装[NSQ](http://nsq.io/deployment/installing.html)

在config.json文件中设置你的相关配置！！

执行以下命令开启消息队列：

```bash
nsqlookupd
nsqd --lookupd-tcp-address=127.0.0.1:4160
nsqd --lookupd-tcp-address=127.0.0.1:4160 --tcp-address=0.0.0.0:4152 --http-address=0.0.0.0:4153
nsqadmin --lookupd-http-address=127.0.0.1:4161
```



执行以下命令运行服务端程序

```shell
go build main.go router.go
./main.exe
```

相关资料：

[极简抖音App使用说明 - 青训营版 - 成电飞书云文档 (feishu.cn)](https://bytedance.feishu.cn/docs/doccnM9KkBAdyDhg8qaeGlIz7S7)

[视频流接口 - 抖音极简版 (apifox.cn)](https://www.apifox.cn/apidoc/shared-8cc50618-0da6-4d5e-a398-76f3b8f766c5/api-18345145)