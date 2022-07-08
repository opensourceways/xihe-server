package authing

import (
	"github.com/Authing/authing-go-sdk/lib/management"
	"github.com/Authing/authing-go-sdk/lib/model"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

var cli *management.Client

func Init(poolId, secret string) {
	cli = management.NewClient(poolId, secret)
}

func NewUserMapper() repositories.UserMapper {
	return userMapper{}
}

type userMapper struct {
}

func (u userMapper) Get(userId string) (do repositories.UserDO, err error) {
	v, err := cli.Detail(userId)
	if err != nil {
		return
	}

	do.Id = userId

	// TODO
	do.Bio = ""

	if v.Email != nil {
		do.Email = *v.Email
	}

	if v.Username != nil {
		do.Account = *v.Username
	}

	if v.Nickname != nil {
		do.Nickname = *v.Nickname
	}

	if v.Photo != nil {
		do.AvatarId = *v.Photo
	}

	if v.Phone != nil {
		do.PhoneNumber = *v.Phone
	}

	return
}

func (u userMapper) Update(do repositories.UserDO) error {
	m := model.UpdateUserInput{}

	//TODO bio
	m.Email = &do.Email
	m.Photo = &do.AvatarId
	m.Phone = &do.PhoneNumber
	m.Nickname = &do.Nickname

	_, err := cli.UpdateUser(do.Id, m)

	return err
}
