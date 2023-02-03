package repositories

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type WuKongPictureMapper interface {
	GetVersion(string) (int, error)
	ListLikesByUserName(string) ([]WuKongPictureDO, int, error)
	InsertIntoLikes(string, *WuKongPictureDO, int) (string, error)
	DeleteLike(string, string) error
	GetLikeByUserName(string, string) (WuKongPictureDO, error)
	GetPublicByUserName(string, string) (WuKongPictureDO, error)
}

func NewWuKongPictureRepository(mapper WuKongPictureMapper) repository.WuKongPicture {
	return wukongPicture{mapper}
}

type wukongPicture struct {
	mapper WuKongPictureMapper
}

func (impl wukongPicture) GetVersion(user domain.Account) (version int, err error) {
	return impl.mapper.GetVersion(user.Account())
}

func (impl wukongPicture) ListLikesByUserName(user domain.Account) (
	[]domain.WuKongPicture, int, error,
) {
	v, version, err := impl.mapper.ListLikesByUserName(user.Account())
	if err != nil {
		return nil, 0, err
	}

	r := make([]domain.WuKongPicture, len(v))
	for i := range v {
		if err = v[i].toWuKongPicture(&r[i]); err != nil {
			return nil, 0, err
		}
	}

	return r, version, nil
}

func (impl wukongPicture) SaveLike(p *domain.WuKongPicture, version int) (string, error) {
	if p.Id != "" {
		return "", errors.New("must be a new picture")
	}

	do := WuKongPictureDO{
		OBSPath:   p.OBSPath,
		CreatedAt: p.CreatedAt,
	}
	do.Style = p.Style
	do.Desc = p.Desc.WuKongPictureDesc()

	v, err := impl.mapper.InsertIntoLikes(p.Owner.Account(), &do, version)
	if err != nil {
		return "", convertError(err)
	}

	return v, nil
}

func (impl wukongPicture) DeleteLike(user domain.Account, pid string) error {
	if err := impl.mapper.DeleteLike(user.Account(), pid); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl wukongPicture) GetLikeByUserName(user domain.Account, pid string) (
	p domain.WuKongPicture, err error,
) {
	if v, err := impl.mapper.GetLikeByUserName(user.Account(), pid); err != nil {
		err = convertError(err)
	} else {
		err = v.toWuKongPicture(&p)
	}

	return
}

func (impl wukongPicture) GetPublicByUserName(user domain.Account, pid string) (
	p domain.WuKongPicture, err error,
) {
	if v, err := impl.mapper.GetPublicByUserName(user.Account(), pid); err != nil {
		err = convertError(err)
	} else {
		err = v.toWuKongPicture(&p)
	}

	return
}

type WuKongPictureDO struct {
	Id        string
	OBSPath   string
	Diggs     []string
	DiggCount int
	CreatedAt string

	WuKongPictureMetaDO
}

func (do *WuKongPictureDO) toWuKongPicture(r *domain.WuKongPicture) (err error) {
	r.Id = do.Id
	r.OBSPath = do.OBSPath
	r.Diggs = do.Diggs
	r.DiggCount = do.DiggCount
	r.CreatedAt = do.CreatedAt

	r.WuKongPictureMeta, err = do.toWuKongPictureMeta()

	return
}
