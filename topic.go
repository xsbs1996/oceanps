package oceanps

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/xsbs1996/oceanps/oceanpsfuncs"
)

var (
	topicMap   = make(map[string]*EventTopic, 0) // map[string]*EventTopic
	topicMapMu = sync.RWMutex{}
)

// DeleteTopicMapConn 根据连接删除相关的订阅者
func DeleteTopicMapConn(ctx context.Context, conn *websocket.Conn, isCancel bool) {
	topicMapMu.Lock()
	defer topicMapMu.Unlock()
	for _, e := range topicMap {
		connListLen := e.Unsubscribe(ctx, conn, isCancel)
		if connListLen <= 0 {
			delete(topicMap, e.queueName)
			e.cancel()
		}
	}
}

type TopicInterface interface {
	Publish(ctx context.Context, msg []byte, exp time.Duration, extra ...string) error // 发布事件
	Subscribe(ctx context.Context, conn *websocket.Conn, ch chan<- []byte)             // 订阅主题
	Unsubscribe(ctx context.Context, conn *websocket.Conn, isCancel bool) int          // 取消订阅主题
}

// EventTopic 主题结构体
type EventTopic struct {
	ctx       context.Context
	cancel    context.CancelFunc
	timeout   time.Duration                     // 超时时间
	topic     string                            // 主题名称
	queueName string                            // 队列名称
	user      string                            // 用户
	mu        sync.RWMutex                      // 读写锁
	ConnList  map[*websocket.Conn]chan<- []byte // 普通订阅者
	MsgChan   chan []byte                       // 消息通道
	pPManage  oceanpsfuncs.PushPullManage       // 消息管理者
	Error     error                             // 错误
}

// NewEventTopic 新建主题
// topic-主题名 user-用户名 timeout-消息传递超时时间,0为无超时
func NewEventTopic(topic string, user string, timeout time.Duration, pPM oceanpsfuncs.PushPullManage) *EventTopic {
	queueName := fmt.Sprintf("%s-%s", topic, user)
	event := GetEventTopicQueue(queueName)
	if event != nil {
		return event
	}

	ctx, cancel := context.WithCancel(context.Background())
	event = &EventTopic{
		ctx:       ctx,
		cancel:    cancel,
		timeout:   timeout,
		topic:     topic,
		queueName: queueName,
		user:      user,
		mu:        sync.RWMutex{},
		ConnList:  make(map[*websocket.Conn]chan<- []byte, 0),
		MsgChan:   make(chan []byte, 1),
		pPManage:  pPM,
	}
	_, ok := SetEventTopicQueue(queueName, event)

	// 拉取消息
	if !ok {
		err := event.pPManage.CheckClient()
		if err != nil {
			event.Error = err
			return event
		}
		go func() {
			if err := event.pPManage.PullMsg(ctx, event.queueName, event.MsgChan); err != nil {
				panic(err)
			}
		}()
	}

	return event
}

// GetEventTopicQueue 获取主题队列
func GetEventTopicQueue(queueName string) *EventTopic {
	topicMapMu.RLock()
	defer topicMapMu.RUnlock()
	event, ok := topicMap[queueName]
	if ok {
		return event
	}
	return nil
}

// SetEventTopicQueue 设置主题队列
func SetEventTopicQueue(queueName string, event *EventTopic) (*EventTopic, bool) {
	topicMapMu.Lock()
	defer topicMapMu.Unlock()
	eventOld, ok := topicMap[queueName]
	if ok {
		return eventOld, ok
	}
	topicMap[queueName] = event
	return event, ok
}

// DelEventTopicQueue 删除主题队列
func DelEventTopicQueue(queueName string) {
	topicMapMu.Lock()
	defer topicMapMu.Unlock()
	_, ok := topicMap[queueName]
	if !ok {
		return
	}
	delete(topicMap, queueName)
}
