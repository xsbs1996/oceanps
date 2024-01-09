package oceanpsfuncs

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	rabbitMqConn *amqp.Connection
	mx           = new(sync.Mutex)
	once         = new(sync.Once)
	c            *RabbitMqPushPull
)

type RabbitMqPushPull struct {
	Method   string `json:"method" required:"true"`
	Ip       string `json:"ip" required:"true"`
	Port     string `json:"port" required:"true"`
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
}

// GetRabbitMqConn 获取连接
func GetRabbitMqConn(c *RabbitMqPushPull) (*amqp.Connection, error) {
	var err error
	if rabbitMqConn == nil || rabbitMqConn.IsClosed() {
		rabbitMqConn, err = initRabbitMq(c)
		if err != nil {
			return nil, err
		}
		go once.Do(heartbeat)
	}
	return rabbitMqConn, nil
}

// initRabbitMq 初始化连接
func initRabbitMq(rabbitConf *RabbitMqPushPull) (*amqp.Connection, error) {
	mx.Lock()
	defer mx.Unlock()
	if rabbitMqConn != nil && !rabbitMqConn.IsClosed() {
		return rabbitMqConn, nil
	}

	amqUrl := fmt.Sprintf("%s://%s:%s@%s:%s", rabbitConf.Method, rabbitConf.Username, rabbitConf.Password, rabbitConf.Ip, rabbitConf.Port)

	var err error
	if rabbitConf.Method == "amqps" {
		rabbitMqConn, err = amqp.DialTLS(amqUrl, &tls.Config{
			InsecureSkipVerify: true,
		})
	} else {
		rabbitMqConn, err = amqp.Dial(amqUrl)
	}

	if err != nil {
		return nil, err
	}
	c = rabbitConf
	return rabbitMqConn, nil
}

// 心跳断线重连
func heartbeat() {
	tick := time.NewTicker(time.Second * 3)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			if rabbitMqConn == nil || rabbitMqConn.IsClosed() {
				rabbitMqConn, _ = initRabbitMq(c)
			}
		}
	}
}

// CheckClient 检测链接
func (c *RabbitMqPushPull) CheckClient() error {
	_, err := GetRabbitMqConn(c)
	if err != nil {
		return err
	}
	return nil
}

// PushMsgFn redis发送订阅消息
func (c *RabbitMqPushPull) PushMsgFn(ctx context.Context, queueName string, msg []byte) error {
	conn, err := GetRabbitMqConn(c)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer func(ch *amqp.Channel) { _ = ch.Close() }(ch)

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	//消息体
	return ch.PublishWithContext(ctx, "", queueName, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         msg,
	})
}

// PullMsgFn rabbitmq拉取订阅消息并发送到管道
func (c *RabbitMqPushPull) PullMsgFn(ctx context.Context, queueName string, msgChan chan<- []byte) error {
	conn, err := GetRabbitMqConn(c)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer func(ch *amqp.Channel) { _ = ch.Close() }(ch)

	err = ch.Qos(1, 0, false)
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	msgList, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	notifyErr := make(chan *amqp.Error, 1)

	for {
		select {
		case <-ctx.Done():
			return nil
		case closeErr := <-ch.NotifyClose(notifyErr):
			if closeErr != nil {
				return err
			}
		case msg := <-msgList:
			if msg.Acknowledger == nil {
				return errors.New("seems to have encountered an unknown problem")
			}
			msgChan <- msg.Body
			_ = msg.Ack(false)
		}
	}

}
