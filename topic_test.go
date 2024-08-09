package oceanps

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/xsbs1996/oceanps/oceanpsfuncs"
	"testing"
	"time"
)

func TestGetEventTopicQueue(t *testing.T) {
	e := NewEventTopic("TestTopic", "User", time.Second*5, &oceanpsfuncs.RedisPushPull{
		Ip:       "127.0.0.1",
		Port:     "6379",
		DB:       0,
		Password: "",
	})
	if e.Error != nil {
		fmt.Println(e.Error)
		return
	}

	conn := &websocket.Conn{}
	wMsg := make(chan []byte, 10)
	e.Subscribe(context.Background(), conn, wMsg)

	go func() {
		for i := 0; i <= 1; i++ {
			err := e.Publish(context.Background(), []byte("hello world"), 0)
			if err != nil {
				fmt.Println("err:", err)
				return
			}
		}
	}()

	go func() {
		time.Sleep(3 * time.Second)
		e.Unsubscribe(context.Background(), conn, true)
	}()

	go func() {
		for v := range e.MsgChan {
			e.ForeachConnList(v)
		}
	}()

	go func() {
		for v := range wMsg {
			fmt.Println(string(v))
		}
	}()

	<-e.ctx.Done()

}
