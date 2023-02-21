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
	GetVersion(domain.Account) (int, error)
	ListLikesByUserName(domain.Account) ([]domain.WuKongPicture, int, error)
	ListPublicsByUserName(domain.Account) ([]domain.WuKongPicture, int, error)
	SaveLike(domain.Account, *domain.WuKongPicture, int) (string, error)
	SavePublic(*domain.WuKongPicture, int) (string, error)
	DeleteLike(domain.Account, string) error
	DeletePublic(domain.Account, string) error
	GetLikeByUserName(domain.Account, string) (domain.WuKongPicture, error)
	GetPublicByUserName(domain.Account, string) (domain.WuKongPicture, error)
	GetPublicsGlobal() ([]domain.WuKongPicture, error)
	UpdatePublicPicture(domain.Account, string, int, *domain.WuKongPicture) error
}
