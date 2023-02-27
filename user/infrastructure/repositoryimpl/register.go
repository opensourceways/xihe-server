package repositoryimpl

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/user/domain/repository"
)

func NewUserRegRepo(m mongodbClient) repository.UserReg {
	return &userRepoImpl{m}
}

type userRepoImpl struct {
	cli mongodbClient
}

func (impl *userRepoImpl) AddUserRegInfo(u *domain.UserRegInfo) error {
	doc, err := impl.genUserRegInfo(u)
	if err != nil {
		return err
	}
	f := func(ctx context.Context) error {
		filter := bson.M{
			fieldAccount: u.Account.Account(),
		}

		_, err := impl.cli.NewDocIfNotExist(ctx, filter, doc)

		return err
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocExists(err) {
			err = repoerr.NewErrorDuplicateCreating(err)
		}
	}

	return err
}

func (impl *userRepoImpl) genUserRegInfo(u *domain.UserRegInfo) (bson.M, error) {
	var d DUserRegInfo
	toUserRegInfoDoc(u, &d)

	return genDoc(d)
}
