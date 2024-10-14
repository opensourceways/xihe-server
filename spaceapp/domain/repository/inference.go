package repository

import (
	"github.com/opensourceways/xihe-server/domain"
	spaceappdomain "github.com/opensourceways/xihe-server/spaceapp/domain"
)

type InferenceSummary struct {
	Id string

	spaceappdomain.InferenceDetail
}

type Inference interface {
	Save(*spaceappdomain.Inference, int) (string, error)
	UpdateDetail(*spaceappdomain.InferenceIndex, *spaceappdomain.InferenceDetail) error
	FindInstance(*spaceappdomain.InferenceIndex) (InferenceSummary, error)
	FindInstances(index *domain.ResourceIndex, lastCommit string) ([]InferenceSummary, int, error)
}
