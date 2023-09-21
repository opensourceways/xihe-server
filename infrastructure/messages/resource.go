package messages

import (
	commsg "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/utils"
)

func NewResourceMessageAdapter(cfg *ResourceConfig, p commsg.Publisher, o commsg.OperateLogPublisher) *resourceMessageAdapter {
	return &resourceMessageAdapter{cfg: *cfg, publisher: p, operateLog: o}
}

type resourceMessageAdapter struct {
	cfg        ResourceConfig
	publisher  commsg.Publisher
	operateLog commsg.OperateLogPublisher
}

func (s *resourceMessageAdapter) AddOperateLogForCreateResource(
	obj domain.ResourceObject, name domain.ResourceName,
) error {
	return s.operateLog.SendOperateLog(obj.Owner.Account(), "resource", map[string]string{
		"id":   obj.Id,
		"name": name.ResourceName(),
		"type": obj.Type.ResourceType(),
	})
}

func (s *resourceMessageAdapter) CreateDataset(e message.DatasetCreatedEvent) error {
	return s.publisher.Publish(
		s.cfg.DatasetCreated.Topic,
		&commsg.MsgNormal{
			Type:      s.cfg.DatasetCreated.Name,
			CreatedAt: utils.Now(),
			Desc:      "Created a dataset",
			User:      e.Account.Account(),
		},
		nil,
	)
}

func (s *resourceMessageAdapter) CreateModel(e message.ModelCreatedEvent) error {
	return s.publisher.Publish(
		s.cfg.ModelCreated.Topic,
		&commsg.MsgNormal{
			Type:      s.cfg.ModelCreated.Name,
			CreatedAt: utils.Now(),
			Desc:      "Created a model",
			User:      e.Account.Account(),
		},
		nil,
	)
}

func (s *resourceMessageAdapter) CreateProject(e message.ProjectCreatedEvent) error {
	return s.publisher.Publish(
		s.cfg.ProjectCreated.Topic,
		&commsg.MsgNormal{
			Type:      s.cfg.ProjectCreated.Name,
			CreatedAt: utils.Now(),
			Desc:      "Created a project",
			User:      e.Account.Account(),
		},
		nil,
	)
}

// RelatedResource
func (s *resourceMessageAdapter) AddRelatedResource(msg *message.RelatedResource) error {
	return s.sendRelatedResource(msg, actionAdd)
}

func (s *resourceMessageAdapter) RemoveRelatedResource(msg *message.RelatedResource) error {
	return s.sendRelatedResource(msg, actionRemove)
}

func (s *resourceMessageAdapter) RemoveRelatedResources(msg *message.RelatedResources) error {
	v := msgRelatedResources{Action: actionRemove}

	toMsgResourceObject(&msg.Promoter, &v.Promoter)

	v.Resources = make([]resourceObjects, len(msg.Resources))
	for i := range msg.Resources {
		toMsgResourceObjects(&msg.Resources[i], &v.Resources[i])
	}

	return s.publisher.Publish(s.cfg.RelatedResource, &v, nil)
}

// Fork
func (s *resourceMessageAdapter) IncreaseFork(msg *domain.ResourceIndex) error {
	v := new(resourceIndex)
	toMsgResourceIndex(msg, v)

	return s.publisher.Publish(s.cfg.Fork, v, nil)
}

func (s *resourceMessageAdapter) sendRelatedResource(msg *message.RelatedResource, action string) error {
	v := msgRelatedResources{Action: action}

	toMsgResourceObject(msg.Promoter, &v.Promoter)

	v.Resources = []resourceObjects{
		{
			Type: msg.Resource.Type.ResourceType(),
			Objects: []resourceIndex{
				{
					Owner: msg.Resource.Owner.Account(),
					Id:    msg.Resource.Id,
				},
			},
		},
	}

	return s.publisher.Publish(s.cfg.RelatedResource, &v, nil)
}

type ResourceConfig struct {
	RelatedResource string             `json:"related_resource" required:"true"`
	Fork            string             `json:"fork"             required:"true"`
	ProjectCreated  commsg.TopicConfig `json:"project_created" required:"true"`
	ModelCreated    commsg.TopicConfig `json:"model_created"           required:"true"`
	DatasetCreated  commsg.TopicConfig `json:"dataset_created"           required:"true"`
}
