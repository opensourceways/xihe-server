package kafka

import (
	"encoding/json"

	kfklib "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/kafka-lib/mq"

	"github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	queueName     = "xihe-kafka-queue"
	deaultVersion = "2.1.0"
)

type Config struct {
	kfklib.Config
}

func (cfg *Config) SetDefault() {
	if cfg.Version == "" {
		cfg.Version = deaultVersion
	}
}

func Init(cfg *Config, log mq.Logger, redis kfklib.Redis) error {
	return kfklib.Init(&cfg.Config, log, redis, queueName)
}

var Exit = kfklib.Exit

func PublisherAdapter() publisherAdapter {
	return publisherAdapter{}
}

func SubscriberAdapter() subscriberAdapter {
	return subscriberAdapter{}
}

func OperateLogPublisherAdapter(topic string, publisher publisherAdapter) operatePublisherAdapter {
	return operatePublisherAdapter{
		topic:     topic,
		publisher: publisher,
	}
}

// publisherAdapter
type publisherAdapter struct{}

func (publisherAdapter) Publish(topic string, v interface{}, header map[string]string) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return kfklib.Publish(topic, header, body)
}

type operatePublisherAdapter struct {
	topic     string
	publisher publisherAdapter
}

func (o operatePublisherAdapter) SendOperateLog(u string, t string, info map[string]string) error {
	return o.publisher.Publish(o.topic, &message.MsgOperateLog{
		When: utils.Now(),
		User: u,
		Type: t,
		Info: info,
	}, nil)
}

// subscriberAdapter
type subscriberAdapter struct{}

func (subscriberAdapter) SubscribeWithStrategyOfRetry(
	group string, h kfklib.Handler, topics []string, retryNum int,
) error {
	return kfklib.SubscribeWithStrategyOfRetry(group, h, topics, retryNum)
}
