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
		err := redisC.PushMsgFnExp(context.Background(), "TestRedisPushMsg"+strconv.Itoa(i), []byte("hello world"), time.Second*30)
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
