package message

import (
	"github.com/opensourceways/xihe-server/cloud/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type MsgCloudConf struct {
	User      string `json:"user"`
	PodId     string `json:"pod_id"`
	CloudId   string `json:"cloud_id"`
	CloudName string `json:"cloud_name"`
}

type MsgPod struct {
	PodId   string `json:"pod_id"`
	CloudId string `json:"cloud_id"`
	Owner   string `json:"owner"`
}

type CloudMessageProducer interface {
	SubscribeCloud(*MsgCloudConf) error
}

func (r *MsgCloudConf) ToMsgCloudConf(c *domain.CloudConf, u types.Account, pid string) {
	*r = MsgCloudConf{
		User:      u.Account(),
		PodId:     pid,
		CloudId:   c.Id,
		CloudName: c.Name.CloudName(),
	}
}

func (r *MsgPod) ToMsgPod(p *domain.Pod) {
	*r = MsgPod{
		PodId:   p.Id,
		CloudId: p.CloudId,
		Owner:   p.Owner.Account(),
	}
}

type CloudMessageHandler interface {
	HandleEventPodSubscribe(info *domain.PodInfo) error
}

type CloudRecordEvent struct {
	Owner  types.Account
	ClouId string
}

type CloudRecordEventPublisher interface {
	Publish(*CloudRecordEvent) error
}
