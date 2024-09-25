package oceanpsfuncs

import (
	"context"
	"github.com/IBM/sarama"
	"time"
)

var (
	kafkaAsyncProduceConn sarama.AsyncProducer // 异步生产者
	kafkaSyncProduceConn  sarama.SyncProducer  // 同步生产者
)

type KafkaMqPushPull struct {
	Host     []string `json:"host" yaml:"host"`
	Username string   `json:"username" required:"true"`
	Password string   `json:"password" required:"true"`
}

// 初始化kafka配置
func initKafkaCfg(kafkaConf *KafkaMqPushPull) (brokers []string, config *sarama.Config) {
	// 创建Sarama配置
	config = sarama.NewConfig()
	// 启用SASL认证
	if kafkaConf.Username != "" && kafkaConf.Password != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = kafkaConf.Username
		config.Net.SASL.Password = kafkaConf.Password
		config.Net.SASL.Handshake = true
	} else {
		config.Net.SASL.Enable = false
	}
	// 构造Broker地址
	brokers = kafkaConf.Host

	return brokers, config
}

// 获取消费者
func getKafkaConsumer(kafkaConf *KafkaMqPushPull) (sarama.Consumer, error) {
	brokers, config := initKafkaCfg(kafkaConf)
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

// 获取异步生产者
func getKafkaAsyncProducer(kafkaConf *KafkaMqPushPull) (sarama.AsyncProducer, error) {
	if kafkaAsyncProduceConn != nil {
		return kafkaAsyncProduceConn, nil
	}

	brokers, config := initKafkaCfg(kafkaConf)
	// 创建异步生产者
	asyncProducer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	kafkaAsyncProduceConn = asyncProducer
	return kafkaAsyncProduceConn, nil
}

// 获取同步生产者
func getKafkaSyncProducer(kafkaConf *KafkaMqPushPull) (sarama.SyncProducer, error) {
	if kafkaSyncProduceConn != nil {
		return kafkaSyncProduceConn, nil
	}
	brokers, config := initKafkaCfg(kafkaConf)
	config.Producer.Return.Successes = true
	// 创建同步生产者
	syncProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	kafkaSyncProduceConn = syncProducer

	return kafkaSyncProduceConn, nil
}

// CheckClient 检测链接
func (c *KafkaMqPushPull) CheckClient() error {
	return nil
}

// PushMsg kafka发送订阅消息
func (c *KafkaMqPushPull) PushMsg(ctx context.Context, topic string, body *PushMsgBody) error {
	if body.IsAsync == true { // 异步
		asyncProducer, err := getKafkaAsyncProducer(c)
		if err != nil {
			return err
		}
		message := &sarama.ProducerMessage{
			Key:       sarama.StringEncoder(body.Key),
			Topic:     topic,
			Value:     sarama.StringEncoder(body.Msg),
			Timestamp: time.Now(),
		}
		asyncProducer.Input() <- message
		return nil
	} else { // 同步
		syncProducer, err := getKafkaSyncProducer(c)
		if err != nil {
			return err
		}
		message := &sarama.ProducerMessage{
			Key:       sarama.StringEncoder(body.Key),
			Topic:     topic,
			Value:     sarama.StringEncoder(body.Msg),
			Timestamp: time.Now(),
		}
		_, _, err = syncProducer.SendMessage(message)
		return err
	}
}

// PullMsg kafka拉取订阅消息并发送到管道 topic={topic}{partition}
func (c *KafkaMqPushPull) PullMsg(ctx context.Context, topic string, msgChan chan<- []byte) error {
	consumer, err := getKafkaConsumer(c)
	if err != nil {
		return err
	}

	// 创建分区消费者
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer func() {
		_ = partitionConsumer.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-partitionConsumer.Errors():
			if err != nil {
				return err
			}
		case msg := <-partitionConsumer.Messages():
			if msg == nil {
				continue
			}
			msgChan <- msg.Value
		}
	}

}
