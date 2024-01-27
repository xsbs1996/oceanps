package oceanpsfuncs

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisPushPull struct {
	Ip       string `json:"ip" required:"true"`
	Port     string `json:"port" required:"true"`
	DB       int    `json:"db" default:"0"`
	Password string `json:"password"`
}

// GetRedisClient 获取redis链接
func GetRedisClient(c *RedisPushPull) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.Ip, c.Port),
		Password: c.Password,
		DB:       c.DB,
	})
	err := redisClient.Ping(context.Background()).Err()
	return redisClient, err
}

// CheckClient 检测链接
func (c *RedisPushPull) CheckClient() error {
	rdb, err := GetRedisClient(c)
	if err != nil {
		return err
	}
	defer func(rdb *redis.Client) { _ = rdb.Close() }(rdb)

	return nil
}

// PushMsgFn redis发送订阅消息
func (c *RedisPushPull) PushMsgFn(ctx context.Context, queueName string, msg []byte) error {
	rdb, err := GetRedisClient(c)
	if err != nil {
		return err
	}
	defer func(rdb *redis.Client) { _ = rdb.Close() }(rdb)

	err = rdb.LPush(ctx, queueName, string(msg)).Err()
	if err != nil {
		return err
	}
	return nil
}

// PullMsgFn redis拉取订阅消息并发送到管道
func (c *RedisPushPull) PullMsgFn(ctx context.Context, queueName string, msgChan chan<- []byte) error {
	defer close(msgChan)
	rdb, err := GetRedisClient(c)
	if err != nil {
		return err
	}
	defer func(rdb *redis.Client) { _ = rdb.Close() }(rdb)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			result, err := rdb.BRPop(ctx, time.Second*1, queueName).Result()
			if err != nil {
				continue
			}
			// 第一个元素是列表的 key，第二个元素是弹出的消息
			msg := result[1]
			msgChan <- []byte(msg)
		}
	}
}
