package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type WuKongPictureListOption struct {
	CountPerPage int
	PageNum      int
}

type WuKong interface {
	ListSamples(string, []int) ([]string, error)
}

type WuKongPicture interface {
	GetVersion(user domain.Account) (int, error)
	ListLikesByUserName(user domain.Account) ([]domain.WuKongPicture, int, error)
	ListPublicsByUserName(user domain.Account) ([]domain.WuKongPicture, int, error)
	SaveLike(*domain.WuKongPicture, int) (string, error)
	SavePublic(*domain.WuKongPicture, int) (string, error)
	DeleteLike(user domain.Account, pid string) error
	GetLikeByUserName(user domain.Account, pid string) (domain.WuKongPicture, error)
	GetPublicByUserName(user domain.Account, pid string) (domain.WuKongPicture, error)
	GetPublicsGlobal()([]domain.WuKongPicture, error)
}
