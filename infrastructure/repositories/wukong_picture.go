package repositories

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type WuKongPictureMapper interface {
	GetVersion(string) (int, error)
	ListLikesByUserName(string) ([]WuKongPictureDO, int, error)
	ListPublicsByUserName(string) ([]WuKongPictureDO, int, error)
	InsertIntoLikes(string, *WuKongPictureDO, int) (string, error)
	InsertIntoPublics(string, *WuKongPictureDO, int) (string, error)
	DeleteLike(string, string) error
	DeletePublic(string, string) error
	GetLikeByUserName(string, string) (WuKongPictureDO, error)
	GetPublicByUserName(string, string) (WuKongPictureDO, error)
	GetPublicsGlobal() ([]WuKongPictureDO, error)
	UpdatePublicPicture(string, string, int, *WuKongPictureDO) error
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

func (impl wukongPicture) ListPublicsByUserName(user domain.Account) (
	[]domain.WuKongPicture, int, error,
) {
	v, version, err := impl.mapper.ListPublicsByUserName(user.Account())
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

func (impl wukongPicture) SaveLike(user domain.Account, p *domain.WuKongPicture, version int) (string, error) {
	if p.Id != "" {
		return "", errors.New("must be a new picture")
	}

	do := new(WuKongPictureDO)
	do.toWuKongPictureDO(p)
	do.setDefault()

	v, err := impl.mapper.InsertIntoLikes(user.Account(), do, version)
	if err != nil {
		return "", convertError(err)
	}

	return v, nil
}

func (impl wukongPicture) SavePublic(p *domain.WuKongPicture, version int) (string, error) {
	if p.Id != "" {
		return "", errors.New("must be a new picture")
	}

	do := new(WuKongPictureDO)
	do.toWuKongPictureDO(p)
	do.setDefault()

	v, err := impl.mapper.InsertIntoPublics(p.Owner.Account(), do, version)
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

func (impl wukongPicture) DeletePublic(user domain.Account, pid string) error {
	if err := impl.mapper.DeletePublic(user.Account(), pid); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl wukongPicture) GetLikeByUserName(user domain.Account, pid string) (
	p domain.WuKongPicture, err error,
) {
	if v, err := impl.mapper.GetLikeByUserName(user.Account(), pid); err != nil {
		err = convertError(err)

		return p, err
	} else {
		err = v.toWuKongPicture(&p)

		return p, err
	}
}

func (impl wukongPicture) GetPublicByUserName(user domain.Account, pid string) (
	p domain.WuKongPicture, err error,
) {
	if v, err := impl.mapper.GetPublicByUserName(user.Account(), pid); err != nil {
		err = convertError(err)

		return p, err
	} else {
		err = v.toWuKongPicture(&p)

		return p, err
	}
}

func (impl wukongPicture) GetPublicsGlobal() (r []domain.WuKongPicture, err error) {
	do, err := impl.mapper.GetPublicsGlobal()
	if err != nil {
		return
	}

	r = make([]domain.WuKongPicture, len(do))

	for i := range do {
		err = do[i].toWuKongPicture(&r[i])
		if err != nil {
			return
		}
	}

	return
}

func (impl wukongPicture) GetOfficialPublicsGlobal() (r []domain.WuKongPicture, err error) {
	do, err := impl.mapper.GetPublicsGlobal()
	if err != nil {
		return
	}

	var offDo []WuKongPictureDO
	for i := range do {
		if do[i].Level == 2 {
			offDo = append(offDo, do[i])
		}
	}

	r = make([]domain.WuKongPicture, len(offDo))

	for i := range offDo {
		err = offDo[i].toWuKongPicture(&r[i])
		if err != nil {
			return
		}
	}

	return
}

func (impl wukongPicture) UpdatePublicPicture(
	user domain.Account, pid string, version int,
	update *domain.WuKongPicture,
) error {
	do := new(WuKongPictureDO)
	do.toWuKongPictureDO(update)

	return impl.mapper.UpdatePublicPicture(user.Account(), pid, version, do)
}

type WuKongPictureDO struct {
	Id        string
	Owner     string
	OBSPath   string
	Level     int
	Diggs     []string
	DiggCount int
	Version   int
	CreatedAt string

	WuKongPictureMetaDO
}

func (do *WuKongPictureDO) setDefault() {
	do.Id = ""
	do.DiggCount = 0
	do.Diggs = []string{}
	do.Version = 1
}

func (do *WuKongPictureDO) toWuKongPictureDO(p *domain.WuKongPicture) {
	*do = WuKongPictureDO{
		Id:        p.Id,
		Owner:     p.Owner.Account(),
		OBSPath:   p.OBSPath,
		CreatedAt: p.CreatedAt,
		Diggs:     p.Diggs,
		DiggCount: p.DiggCount,
		WuKongPictureMetaDO: WuKongPictureMetaDO{
			Style: p.Style,
			Desc:  p.Desc.WuKongPictureDesc(),
		},
	}
}

func (do *WuKongPictureDO) toWuKongPicture(r *domain.WuKongPicture) (err error) {
	user, err := domain.NewAccount(do.Owner)
	if err != nil {
		return
	}

	r.Owner = user
	r.Id = do.Id
	r.OBSPath = do.OBSPath
	r.Diggs = do.Diggs
	r.DiggCount = do.DiggCount
	r.Version = do.Version
	r.CreatedAt = do.CreatedAt

	r.WuKongPictureMeta, err = do.toWuKongPictureMeta()

	return
}
