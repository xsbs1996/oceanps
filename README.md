# oceanps
golang消息队列，可灵活支持多种方式，目前内置方式如下
1. redis
2. rabbitMQ 

# 使用方法
一. 使用内置方式 redis/rabbitMQ
1. 将 oceanpsfuncs 目录下的 redis/rabbitMQ 配置文件 RedisPushPull/RabbitMqPushPull 填充完毕
   ```
   type RedisPushPull struct {
        Ip       string // ip 地址
        Port     string // 端口
        DB       int    // 库
        Password string // 密码
   }
   
   type RabbitMqPushPull struct {
	    Method   string
	    Ip       string
	    Port     string
	    Username string
	    Password string
    }
   ```
2. 订阅消息
   ```
   // 全员主题第二个参数传递空字符串，如果每个用户是一个主题第二个参数则传值
   oceanps.NewEventTopic("队列名称", "会员主题名称", config.OceanpsRedisConf)
    ```

3. 发布消息 
    ```
    直接调用 RedisPushPull/RabbitMqPushPull 的 PushMsgFn 方法即可
    ```
   
二. 灵活扩展方式

1. 实现 oceanpsfuncs 下的 PushPullManage 接口

2. 使用 oceanpsfuncs 下的RegisterPushPull函数将实现此 PushPullManage 接口的结构体注册到PushPullMap

3. 使用 NewEventTopic 方法New一个主题，将 PushPullManage 接口传入,可用 oceanpsfuncs 下的 GetPushPull 函数获取你已经注册的接口

4. NewEventTopic 方法返回一个结构体,通过此结构体的 MsgChan channel传递消息

5. 详细方法可见 oceanpsfuncs->redisfunc_test.go 与 topic_test.go

三. 使用demo可看 https://github.com/xsbs1996/oceanps_demo 项目
