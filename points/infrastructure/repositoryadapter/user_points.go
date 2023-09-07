package repositoryadapter

import (
	"errors"

	common "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/points/domain"
	"github.com/opensourceways/xihe-server/points/domain/repository"
)

func UserPointsAdapter() *userPointsAdapter {
	return &userPointsAdapter{}
}

type userPointsAdapter struct {
}

func (impl *userPointsAdapter) SavePointsItem(*domain.UserPoints, *domain.PointsItem) error {
	return errors.New("unimplemented")
}

func (impl *userPointsAdapter) Find(account common.Account, date string) (domain.UserPoints, error) {
	return domain.UserPoints{}, errors.New("unimplemented")
}

func (impl *userPointsAdapter) FindAll(account common.Account) (repository.UserPointsDetails, error) {
	return repository.UserPointsDetails{}, errors.New("unimplemented")
}
