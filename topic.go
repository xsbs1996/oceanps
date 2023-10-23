package oceanps

import (
	"context"
	"fmt"
	"github.com/xsbs1996/oceanps/funcs"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	topicMap   = make(map[string]*EventTopic, 0) // map[string]*EventTopic
	topicMapMu = sync.RWMutex{}
)

// DeleteTopicMapConn 根据连接删除相关的订阅者
func DeleteTopicMapConn(conn *websocket.Conn, isCancel bool) {
	topicMapMu.Lock()
	defer topicMapMu.Unlock()
	for _, e := range topicMap {
		connListLen := e.Unsubscribe(conn, isCancel)
		if connListLen <= 0 {
			delete(topicMap, e.queueName)
			e.cancel()
		}
	}
}

type TopicInterface interface {
	Publish(msg []byte) error                            // 发布事件
	Subscribe(conn *websocket.Conn, ch chan<- []byte)    // 订阅主题
	Unsubscribe(conn *websocket.Conn, isCancel bool) int // 取消订阅主题
}

// EventTopic 主题结构体
type EventTopic struct {
	ctx       context.Context
	cancel    context.CancelFunc
	topic     string                            // 主题名称
	queueName string                            // 队列名称
	user      string                            // 用户
	mu        sync.RWMutex                      // 读写锁
	ConnList  map[*websocket.Conn]chan<- []byte // 普通订阅者
	MsgChan   chan []byte                       // 消息通道
	pPManage  funcs.PushPullManage              // 消息管理者
}

// NewEventTopic 新建主题
func NewEventTopic(topic string, user string, pPM funcs.PushPullManage) *EventTopic {
	queueName := fmt.Sprintf("%s-%s", topic, user)
	event := GetEventTopicQueue(queueName)
	if event != nil {
		return event
	}

	ctx, cancel := context.WithCancel(context.Background())
	event = &EventTopic{
		ctx:       ctx,
		cancel:    cancel,
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
		go func() {
			err := event.pPManage.PullMsgFn(ctx, event.queueName, event.MsgChan)
			if err != nil {
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
