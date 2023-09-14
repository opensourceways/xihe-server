package messageadapter

import (
	"fmt"

	"github.com/opensourceways/xihe-server/cloud/domain/message"
	common "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/utils"
)

type CloudCreateMsg struct {
	common.MsgNormal
	PodId     string `json:"pod_id"`
	CloudId   string `json:"cloud_id"`
	CloudName string `json:"cloud_name"`
}

type CloudReleaseMsg struct {
	common.MsgNormal
	message.MsgPod
}

func NewPublisher(cfg *Config, p common.Publisher) *publisher {
	return &publisher{*cfg, p}
}

type publisher struct {
	cfg       Config
	publisher common.Publisher
}

func (s publisher) SubscribeCloud(m *message.MsgCloudConf) error {
	msg := CloudCreateMsg{
		MsgNormal: common.MsgNormal{
			Type:      s.cfg.JupyterCreated.Name,
			User:      m.User,
			CreatedAt: utils.Now(),
			Desc:      fmt.Sprintf("start a jupyter notebook on %s", m.CloudName),
		},
		PodId:     m.PodId,
		CloudId:   m.CloudId,
		CloudName: m.CloudName,
	}

	return s.publisher.Publish(s.cfg.JupyterCreated.Topic, msg, nil)
}

func (s publisher) ReleasePod(m *message.MsgPod) error {
	// no impl for now
	return nil
}

// Config
type Config struct {
	JupyterCreated common.TopicConfig `json:"jupyter_created"`
}
