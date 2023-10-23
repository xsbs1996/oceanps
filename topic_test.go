package oceanps

import (
	"fmt"
	"github.com/xsbs1996/oceanps/oceanpsfuncs"
	"testing"
)

func TestGetEventTopicQueue(t *testing.T) {
	e := NewEventTopic("TestTopic", "User", &oceanpsfuncs.RedisPushPull{
		Ip:       "127.0.0.1",
		Port:     "6379",
		DB:       0,
		Password: "",
	})

	go func() {
		for i := 0; i <= 100; i++ {
			err := e.Publish([]byte("hello world"))
			if err != nil {
				fmt.Println("err:", err)
				return
			}
		}
	}()

	for v := range e.MsgChan {
		fmt.Println(string(v))
	}
}
