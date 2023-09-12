package message

import "github.com/opensourceways/xihe-server/competition/domain"

type MessageProducer interface {
	SendWorkSubmittedEvent(*domain.WorkSubmittedEvent) error
}
