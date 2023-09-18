package kafka

import (
	"encoding/json"

	kfklib "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/kafka-lib/mq"
)

const (
	queueName     = "xihe-kafka-queue"
	deaultVersion = "2.1.0"
)

type Config struct {
	kfklib.Config
}

func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Config,
	}
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

// publisherAdapter
type publisherAdapter struct{}

func (publisherAdapter) Publish(topic string, v interface{}, header map[string]string) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return kfklib.Publish(topic, header, body)
}

// subscriberAdapter
type subscriberAdapter struct{}

func (subscriberAdapter) SubscribeWithStrategyOfRetry(
	group string, h kfklib.Handler, topics []string, retryNum int,
) error {
	return kfklib.SubscribeWithStrategyOfRetry(group, h, topics, retryNum)
}
