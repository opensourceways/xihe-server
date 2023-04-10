package message

import "github.com/opensourceways/xihe-server/domain/message"

type MsgWuKongLinks struct {
	Type   string            `json:"type"`
	TaskId uint64            `json:"task_id"`
	User   string            `json:"user"`
	Links  map[string]string `json:"links"`
}

type AsyncMessageProducer interface {
	message.Sender
	UpdateWuKongTask(*MsgWuKongLinks) error
}
