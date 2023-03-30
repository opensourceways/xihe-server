package repositoryimpl

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	types "github.com/opensourceways/xihe-server/domain"
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
	return impl.newUserRegInfo(u)
}

func (impl *userRepoImpl) newUserRegInfo(u *domain.UserRegInfo) error {
	doc, err := impl.genUserRegInfo(u)
	if err != nil {
		return err
	}
	doc[fieldVersion] = u.Version

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

func (impl *userRepoImpl) UpdateUserRegInfo(u *domain.UserRegInfo, version int) error {
	doc, err := impl.genUserRegInfo(u)
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		filter := bson.M{
			fieldAccount: u.Account.Account(),
		}

		return impl.cli.UpdateDoc(ctx, filter, doc, mongoCmdSet, version)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}
	}

	return err
}

func (impl *userRepoImpl) GetUserRegInfo(user types.Account) (u domain.UserRegInfo, err error) {
	var v DUserRegInfo

	f := func(ctx context.Context) error {
		filter := impl.docFilter(fieldAccount, user.Account())

		return impl.cli.GetDoc(ctx, filter, nil, &v)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = nil
		}

		return
	}

	if err = v.toUserRegInfo(&u); err != nil {
		return
	}

	return
}

func (impl *userRepoImpl) genUserRegInfo(u *domain.UserRegInfo) (bson.M, error) {
	var d DUserRegInfo
	toUserRegInfoDoc(u, &d)

	return genDoc(d)
}

func (impl *userRepoImpl) docFilter(fieldName, value string) bson.M {
	return bson.M{
		fieldName: value,
	}
}
