# oceanps
golang消息队列，可灵活支持多种方式，内置redis方式

# 使用方法
1.实现 oceanpsfuncs 下的 PushPullManage 接口

2.使用 oceanpsfuncs 下的RegisterPushPull函数将实现此 PushPullManage 接口的结构体注册到PushPullMap

3.使用 NewEventTopic 方法New一个主题，将 PushPullManage 接口传入,可用 oceanpsfuncs 下的 GetPushPull 函数获取你已经注册的接口

4.NewEventTopic 方法返回一个结构体,通过此结构体的 MsgChan channel传递消息

5.详细方法可见 oceanpsfuncs->redisfunc_test.go 与 topic_test.go

6. 使用damo可看https://github.com/xsbs1996/oceanps_demo项目
