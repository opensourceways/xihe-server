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

type configInterface interface {
	Validate() error
	SetDefault()
}

var _ configInterface = (*Config)(nil)

type Config struct {
	*kfklib.Config
}

func (cfg *Config) SetDefault() {
	if cfg.Version == "" {
		cfg.Version = deaultVersion
	}
}

func Init(cfg *Config, log mq.Logger, redis kfklib.Redis) error {
	return kfklib.Init(cfg.Config, log, redis, queueName)
}

var Exit = kfklib.Exit

func PlublisherAdapter() plublisherAdapter {
	return plublisherAdapter{}
}

func SubscriberAdapter() subscriberAdapter {
	return subscriberAdapter{}
}

// plublisherAdapter
type plublisherAdapter struct{}

func (plublisherAdapter) Publish(topic string, v interface{}, header map[string]string) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return kfklib.Publish(topic, header, body)
}

// subscriberAdapter
type subscriberAdapter struct{}

func (subscriberAdapter) Subscribe(group string, h kfklib.Handler, topics []string) error {
	return kfklib.Subscribe(group, h, topics)
}

func (subscriberAdapter) SubscribeWithStrategyOfRetry(
	group string, h kfklib.Handler, topics []string, retryNum int,
) error {
	return kfklib.SubscribeWithStrategyOfRetry(group, h, topics, retryNum)
}
