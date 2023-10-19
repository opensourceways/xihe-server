package repositories

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type LuoJiaMapper interface {
	Insert(*UserLuoJiaRecordDO) (string, error)
	List(string) ([]LuoJiaRecordDO, error)
}

func NewLuoJiaRepository(mapper LuoJiaMapper) repository.LuoJia {
	return luojia{mapper}
}

type luojia struct {
	mapper LuoJiaMapper
}

func (impl luojia) Save(ur *domain.UserLuoJiaRecord) (r domain.LuoJiaRecord, err error) {
	if ur.Id != "" {
		err = errors.New("must be a new luojia")

		return
	}

	do := new(UserLuoJiaRecordDO)
	impl.toUserLuoJiaRecordDO(ur, do)

	v, err := impl.mapper.Insert(do)
	if err != nil {
		err = convertError(err)
	} else {
		r = ur.LuoJiaRecord
		r.Id = v
	}

	return
}

func (impl luojia) List(user domain.Account) (r []domain.LuoJiaRecord, err error) {
	v, err := impl.mapper.List(user.Account())
	if err != nil {
		err = convertError(err)

		return
	}

	if len(v) == 0 {
		return
	}

	r = make([]domain.LuoJiaRecord, len(v))
	for i := range v {
		v[i].toRecord(&r[i])
	}

	return
}

func (impl luojia) toUserLuoJiaRecordDO(
	r *domain.UserLuoJiaRecord, do *UserLuoJiaRecordDO,
) {
	*do = UserLuoJiaRecordDO{
		User: r.User.Account(),
	}

	do.CreatedAt = r.CreatedAt
}

type UserLuoJiaRecordDO struct {
	User string

	LuoJiaRecordDO
}

type LuoJiaRecordDO struct {
	Id        string
	CreatedAt int64
}

func (do *LuoJiaRecordDO) toRecord(r *domain.LuoJiaRecord) {
	*r = domain.LuoJiaRecord{
		Id:        do.Id,
		CreatedAt: do.CreatedAt,
	}
}

type LuoJiaRecordIndexDO struct {
	User string
	Id   string
}
