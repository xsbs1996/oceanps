package oceanps

import (
	"time"

	"github.com/gorilla/websocket"
)

// Subscribe 订阅
func (e *EventTopic) Subscribe(conn *websocket.Conn, ch chan<- []byte) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.ConnList[conn] = ch
}

// Unsubscribe 取消订阅 isCancel-是否取消队列
func (e *EventTopic) Unsubscribe(conn *websocket.Conn, isCancel bool) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, ok := e.ConnList[conn]
	if ok {
		delete(e.ConnList, conn)
		if isCancel {
			DelEventTopicQueue(e.queueName)
			e.cancel()
		}
	}
	return len(e.ConnList)
}

// GetConnList 获取订阅者
func (e *EventTopic) GetConnList(conn *websocket.Conn) (chan<- []byte, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	ch, ok := e.ConnList[conn]
	return ch, ok
}

// ForeachConnList 遍历订阅者并发送
func (e *EventTopic) ForeachConnList(msgJson []byte) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	for _, writeMsg := range e.ConnList {
		go e.sendMsg(writeMsg, msgJson)
	}
}

// 发送消息,x秒延迟
func (e *EventTopic) sendMsg(writeMsg chan<- []byte, msgJson []byte) {
	if e.timeout > 0 {
		ticker := time.NewTicker(time.Second * 5)
		defer ticker.Stop()
		for {
			select {
			case writeMsg <- msgJson:
				return
			case <-ticker.C:
				return
			}
		}
	} else {
		for {
			select {
			case writeMsg <- msgJson:
				return
			}
		}
	}
}

// LenConnList 查看订阅者数量
func (e *EventTopic) LenConnList() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.ConnList)
}
