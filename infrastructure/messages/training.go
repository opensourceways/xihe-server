package messages

import (
	"encoding/json"
	"fmt"

	commsg "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	projectOwner = "project_owner"
	projectId    = "project_id"
	trainingId   = "training_id"
	input        = "input"
)

func NewTrainingMessageAdapter(cfg *TrainingConfig, p commsg.Publisher) *trainingMessageAdapter {
	return &trainingMessageAdapter{cfg: *cfg, publisher: p}
}

type trainingMessageAdapter struct {
	cfg       TrainingConfig
	publisher commsg.Publisher
}

func (impl *trainingMessageAdapter) SendTrainingCreated(v *domain.TrainingCreatedEvent) error {
	cfg := &impl.cfg.TrainingCreated

	bytes, err := json.Marshal(v.TrainingInputs)
	if err != nil {
		return err
	}

	msg := commsg.MsgNormal{
		Type: cfg.Name,
		User: v.Account.Account(),
		Desc: fmt.Sprintf("create training, id: %s", v.TrainingIndex.TrainingId),
		Details: map[string]string{
			projectOwner: v.TrainingIndex.Project.Owner.Account(),
			projectId:    v.TrainingIndex.Project.Id,
			trainingId:   v.TrainingIndex.TrainingId,
			input:        string(bytes),
		},
		CreatedAt: utils.Now(),
	}

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

type TrainingConfig struct {
	TrainingCreated commsg.TopicConfig `json:"training_created" required:"true"`
}
