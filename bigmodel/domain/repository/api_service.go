package repository

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type ApiService interface {
	ApplyApi(*domain.UserApiRecord) error
	GetApiByUserModel(types.Account, domain.ModelName) (domain.UserApiRecord, error)
	GetApiByUser(types.Account) ([]domain.UserApiRecord, error)
	AddApiCallCount(types.Account, domain.ModelName, int) error
	UpdateToken(types.Account, domain.ModelName, string, string, int) error
}
