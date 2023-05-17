package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type LoginMapper interface {
	Insert(LoginDO) error
	Get(string) (LoginDO, error)
}

// TODO: mapper can be mysql
func NewLoginRepository(mapper LoginMapper) repository.Login {
	return login{mapper}
}

type login struct {
	mapper LoginMapper
}

func (impl login) Get(account domain.Account) (r domain.Login, err error) {
	do, err := impl.mapper.Get(account.Account())
	if err != nil {
		err = convertError(err)

		return
	}

	err = do.toLogin(&r)

	return
}

func (impl login) Save(u *domain.Login) (err error) {
	if err = impl.mapper.Insert(impl.toLoginDO(u)); err != nil {
		err = convertError(err)
	}

	return
}

func (impl login) toLoginDO(u *domain.Login) LoginDO {
	return LoginDO{
		Account:     u.Account.Account(),
		Info:        u.Info,
		AccessToken: u.AccessToken,
	}
}

type LoginDO struct {
	Account     string
	Info        string
	AccessToken string
}

func (do *LoginDO) toLogin(r *domain.Login) (err error) {
	if r.Account, err = domain.NewAccount(do.Account); err != nil {
		return
	}

	r.Info = do.Info
	r.AccessToken = do.AccessToken

	return
}
