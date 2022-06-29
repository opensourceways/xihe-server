package app

import (
	"github.com/opensourceways/xihe-server/domain"
)

type ProjectRepository interface {
	Save(domain.Project) error
}
