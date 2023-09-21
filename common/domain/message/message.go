package message

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
)

type MsgNormal struct {
	Type      string            `json:"type"`
	User      string            `json:"user"`
	Desc      string            `json:"desc"`
	Details   map[string]string `json:"details"`
	CreatedAt int64             `json:"created_at"`
}

type TopicConfig struct {
	// Name is the event name
	Name  string `json:"name"   required:"true"`
	Topic string `json:"topic"  required:"true"`
}

type Publisher interface {
	Publish(topic string, v interface{}, header map[string]string) error
}

type OperateLogPublisher interface {
	SendOperateLog(user string, t string, info map[string]string) error
}

type Subscriber interface {
	SubscribeWithStrategyOfRetry(group string, h kfklib.Handler, topics []string, retryNum int) error
}
