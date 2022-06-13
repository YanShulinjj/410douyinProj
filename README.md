## 410 到底对不队 抖音项目



注意事项：

> 需要安装[FFmpeg](http://ffmpeg.org/download.html),   并修改controller/publish.go  row: 47 的ffmpeg路径
>
> 需要安装MySQL，并创建golang_mysql 数据库， 并修改dbmodel.go 文件中连接mysql的账户和密码
>
> 需要修改BaseURL.
>
> 在config.json文件中设置你的相关配置！！
>
> 若没有 public 文件夹，需要新建
>
> *新增 需要安装[NSQ](http://nsq.io/deployment/installing.html)



执行以下命令开启消息队列：[*注意每条命令开启一个终端，或让命令**后台运行**]

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

文件说明：

```bash
│  go.mod      
│  go.sum
│  main.go              # 程序入口
│  README.md
│  router.go            # gin路由
│  
├─config
│      config.json      # 相关配置文件
│      
├─controller
│      comment.go       # 评论业务逻辑
│      common.go        # 通用数据结构的定义
│      config.go        # 解析配置脚本
│      consumer.go      # [消息队列] 消费者，定义处理消息的业务逻辑
│      dbmodel.go       # [MySQL] 对mysql的crud
│      demo_data.go     # 得到视频流
│      favorite.go      # 点赞/取消点赞
│      feed.go          # 将视频流发给前端
│      producer.go      # [消息队列] 生产者，定义发送消息的类型
│      publish.go       # 发布视频业务逻辑
│      redis.go         # [Redis] 缓存一些map数据
│      relation.go      # 关注/取消关注
│      user.go          # 用户注册/登陆
│      
├─mylog
│      a.log            # 日志文件
│      logger.go        # 配置一个日志对象
│      
└─public                # 保存上传的视频以及封面

```



相关资料：

[极简抖音App使用说明 - 青训营版 - 成电飞书云文档 (feishu.cn)](https://bytedance.feishu.cn/docs/doccnM9KkBAdyDhg8qaeGlIz7S7)

[视频流接口 - 抖音极简版 (apifox.cn)](https://www.apifox.cn/apidoc/shared-8cc50618-0da6-4d5e-a398-76f3b8f766c5/api-18345145)
