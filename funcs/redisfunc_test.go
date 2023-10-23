package funcs

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestRedisPushMsg(t *testing.T) {
	redisC := &RedisPushPull{
		Ip:       "127.0.0.1",
		Port:     "6379",
		DB:       0,
		Password: "",
	}
	for i := 0; i <= 100; i++ {
		err := redisC.PushMsgFn(context.Background(), "TestRedisPushMsg", []byte("hello world"))
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(time.Second * 2)
	}
	return
}

func TestRedisPullMsg(t *testing.T) {
	redisC := &RedisPushPull{
		Ip:       "127.0.0.1",
		Port:     "6379",
		DB:       0,
		Password: "",
	}
	msgChan := make(chan []byte, 10)
	go func() {
		err := redisC.PullMsgFn(context.Background(), "TestRedisPushMsg", msgChan)
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
