package message

import (
	"github.com/opensourceways/xihe-server/cloud/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type MsgCloudConf struct {
	User          string `json:"user"`
	PodId         string `json:"pod_id"`
	CloudId       string `json:"cloud_id"`
	CloudName     string `json:"cloud_name"`
	CloudImage    string `json:"cloud_image"`
	CloudCardsNum int    `json:"cloud_cards_num"`
}

type MsgPod struct {
	PodId   string `json:"pod_id"`
	CloudId string `json:"cloud_id"`
	Owner   string `json:"owner"`
}

type CloudMessageProducer interface {
	SubscribeCloud(*MsgCloudConf) error
	ReleaseCloud(*ReleaseCloudEvent) error
}

func (r *MsgCloudConf) ToMsgCloudConf(
	c *domain.CloudConf, u types.Account, pid string, cloudImage domain.ICloudImage, cardsNum domain.CloudSpecCardsNum,
) {
	*r = MsgCloudConf{
		User:          u.Account(),
		PodId:         pid,
		CloudId:       c.Id,
		CloudName:     c.Name.CloudName(),
		CloudImage:    cloudImage.Image(),
		CloudCardsNum: cardsNum.CloudSpecCardsNum(),
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
	HandleEventPodRelease(podId, cloudType string) error
}

type CloudRecordEvent struct {
	Owner  types.Account
	ClouId string
}

type CloudRecordEventPublisher interface {
	Publish(*CloudRecordEvent) error
}

type ReleaseCloudEvent struct {
	PodId     string `json:"pod_id"`
	CloudType string `json:"cloud_type"`
}
