## 410 到底对不队 抖音项目



注意事项：

> 需要安装[FFmpeg](http://ffmpeg.org/download.html),   并修改controller/publish.go  row: 47 的ffmpeg路径
>
<<<<<<< HEAD
> 需要安装MySQL，并创建golang_mysql 数据库， 并修改dbmodel.go 文件中连接mysql的账户和密码

=======
> 需要安装MySQL，并创建golang_mysql 数据库
> 
> 需要修改BaseURL.
>>>>>>> 95911ffeae214caea9be587255c3079ae83f5fbc
执行以下命令运行服务端程序

```shell
go build main.go router.go
./main.exe
```

相关资料：

[极简抖音App使用说明 - 青训营版 - 成电飞书云文档 (feishu.cn)](https://bytedance.feishu.cn/docs/doccnM9KkBAdyDhg8qaeGlIz7S7)

[视频流接口 - 抖音极简版 (apifox.cn)](https://www.apifox.cn/apidoc/shared-8cc50618-0da6-4d5e-a398-76f3b8f766c5/api-18345145)