## 410 到底对不队 抖音项目

具体功能内容参考飞书说明文档



注意事项：

> 需要安装[FFmpeg](http://ffmpeg.org/download.html),   并修改controller/publish.go  row: 47 的ffmpeg路径
>
> 需要安装MySQL，并创建golang_mysql 数据库
> 
> 需要修改BaseURL.
执行以下命令运行服务端程序

```shell
go build main.go router.go
./main.exe
```

