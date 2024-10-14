package inference

import (
	"github.com/opensourceways/xihe-server/spaceapp/domain"
)

type InferenceInfo struct {
	*domain.InferenceInfo
	UserToken string
}

type Inference interface {
	Create(*InferenceInfo) (int, error)
	GetSurvivalTime(*domain.InferenceInfo) int
	ExtendSurvivalTime(index *domain.InferenceIndex, timeToExtend int) error
}
