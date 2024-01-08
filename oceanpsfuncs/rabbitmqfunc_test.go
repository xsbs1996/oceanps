package oceanpsfuncs

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestRabbitMqPushMsg(t *testing.T) {
	rabbitC := &RabbitMqPushPull{
		Method:   "amqp",
		Ip:       "127.0.0.1",
		Port:     "5672",
		Username: "root",
		Password: "123456",
	}
	for i := 0; i <= 100; i++ {
		err := rabbitC.PushMsgFn(context.Background(), "TestRabbitMqPushMsg", []byte("hello world"+strconv.Itoa(i)))
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(time.Second * 1)
	}
	return
}

func TestRabbitPullMsg(t *testing.T) {
	rabbitC := &RabbitMqPushPull{
		Method:   "amqp",
		Ip:       "127.0.0.1",
		Port:     "5672",
		Username: "root",
		Password: "123456",
	}
	msgChan := make(chan []byte, 10)
	go func() {
		err := rabbitC.PullMsgFn(context.Background(), "TestRabbitMqPushMsg", msgChan)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()
	for {
		select {
		case msg := <-msgChan:
			fmt.Println(string(msg))
		}
	}
}
