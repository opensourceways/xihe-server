package message

import (
	"encoding/json"

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

func Publish(topic string, v interface{}, header map[string]string) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return kfklib.Publish(topic, header, body)
}
