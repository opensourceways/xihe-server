package messages

import (
	commsg "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

type msgLike struct {
	Action   string         `json:"action"`
	Resource resourceObject `json:"resource"`

	commsg.MsgNormal
}

type LikeConfig struct {
	ModelLiked   commsg.TopicConfig `json:"model_liked"   required:"true"`
	ProjectLiked commsg.TopicConfig `json:"project_liked" required:"true"`
	DatasetLiked commsg.TopicConfig `json:"dataset_liked" required:"true"`
}

func NewLikeMessageAdapter(topic string, cfg *LikeConfig, p commsg.Publisher) *likeMessageAdapter {
	return &likeMessageAdapter{topic: topic, cfg: *cfg, publisher: p}
}

type likeMessageAdapter struct {
	cfg       LikeConfig
	topic     string
	publisher commsg.Publisher
}

func (s *likeMessageAdapter) toLikePointsMsg(t domain.ResourceType, u string) commsg.MsgNormal {
	m := commsg.MsgNormal{
		CreatedAt: utils.Now(),
		User:      u,
	}

	switch t {
	case domain.ResourceTypeDataset:
		m.Type = s.cfg.DatasetLiked.Name
		m.Desc = "Liked a dataset"

	case domain.ResourceTypeProject:
		m.Type = s.cfg.ProjectLiked.Name
		m.Desc = "Liked a project"

	case domain.ResourceTypeModel:
		m.Type = s.cfg.ModelLiked.Name
		m.Desc = "Liked a model"

	default:
		m = commsg.MsgNormal{}
	}

	return m
}

func (s *likeMessageAdapter) toLikeMsg(msg *domain.ResourceObject, action string) msgLike {
	v := msgLike{
		Action: action,
	}

	toMsgResourceObject(msg, &v.Resource)

	if action == actionAdd {
		v.MsgNormal = s.toLikePointsMsg(msg.Type, msg.Owner.Account())
	}

	return v
}

// Like
func (s *likeMessageAdapter) AddLike(msg *domain.ResourceObject) error {
	return s.sendLike(msg, actionAdd)
}

func (s *likeMessageAdapter) RemoveLike(msg *domain.ResourceObject) error {
	return s.sendLike(msg, actionRemove)
}

// we send all the projectLiked/modelLiked/datasetLikded msg to like topic
// but with different Type in MsgNormal
func (s *likeMessageAdapter) sendLike(msg *domain.ResourceObject, action string) error {
	return s.publisher.Publish(s.topic, s.toLikeMsg(msg, action), nil)
}
