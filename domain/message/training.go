package message

import (
	"github.com/opensourceways/xihe-server/domain"
)

type MessageProducer interface {
	SendTrainingCreated(*domain.TrainingCreatedEvent) error
}
