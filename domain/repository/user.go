package repository

import "github.com/opensourceways/xihe-server/domain/entity"

type UserRepository interface {
	Save(item *entity.User) (*entity.User, error)
}
