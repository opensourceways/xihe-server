package repository

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type WuKong interface {
	ListSamples(string, []int) ([]string, error)
}

type WuKongPicture interface {
	GetVersion(types.Account) (int, error)
	ListLikesByUserName(types.Account) ([]domain.WuKongPicture, int, error)
	ListPublicsByUserName(types.Account) ([]domain.WuKongPicture, int, error)
	SaveLike(types.Account, *domain.WuKongPicture, int) (string, error)
	SavePublic(*domain.WuKongPicture, int) (string, error)
	DeleteLike(types.Account, string) error
	DeletePublic(types.Account, string) error
	GetLikeByUserName(types.Account, string) (domain.WuKongPicture, error)
	GetPublicByUserName(types.Account, string) (domain.WuKongPicture, error)
	GetPublicsGlobal() ([]domain.WuKongPicture, error)
	GetOfficialPublicsGlobal() ([]domain.WuKongPicture, error)
	UpdatePublicPicture(types.Account, string, int, *domain.WuKongPicture) error
}
