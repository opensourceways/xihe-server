package repository

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain"
)

type ApiInfo interface {
	GetApiInfo(domain.ModelName) (domain.ApiInfo, error)
}
