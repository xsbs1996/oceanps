package oceanpsfuncs

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestRedisPushMsg(t *testing.T) {
	redisC := &RedisPushPull{
		Ip:       "127.0.0.1",
		Port:     "6379",
		DB:       0,
		Password: "exchange",
	}
	for i := 0; i <= 100; i++ {
		err := redisC.PushMsg(context.Background(), "TestRedisPushMsg"+strconv.Itoa(i), &PushMsgBody{Msg: []byte("hello world" + strconv.Itoa(i)), Exp: time.Second * 30})
		if err != nil {
			fmt.Println(err)
			return
		}

	}
	return
}

func TestRedisPullMsg(t *testing.T) {
	redisC := &RedisPushPull{
		Ip:       "127.0.0.1",
		Port:     "6379",
		DB:       0,
		Password: "exchange",
	}
	msgChan := make(chan []byte, 10)
	go func() {
		err := redisC.PullMsg(context.Background(), "TestRabbitMqPushMsg", msgChan)
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
