package taskdoc

import (
	common "github.com/opensourceways/xihe-server/common/domain"
	"github.com/opensourceways/xihe-server/points/domain"
)

type TaskDoc interface {
	Doc([]domain.Task, common.Language) ([]byte, error)
}
