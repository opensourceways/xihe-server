package message

import "github.com/opensourceways/xihe-server/competition/domain"

type MessageProducer interface {
	SendWorkSubmittedEvent(*domain.WorkSubmittedEvent) error
	SendCompetitorAppliedEvent(*domain.CompetitorAppliedEvent) error
}
