package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func loginDocFilter(account string) bson.M {
	return bson.M{
		fieldAccount: account,
	}
}

func NewLoginMapper(name string) repositories.LoginMapper {
	return login{name}
}

type login struct {
	collectionName string
}

func (col login) Insert(do repositories.LoginDO) error {
	doc, err := col.toLoginDoc(&do)
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		_, err := cli.replaceDoc(
			ctx, col.collectionName,
			loginDocFilter(do.Account), doc,
		)

		return err
	}

	return withContext(f)
}

func (col login) Get(account string) (do repositories.LoginDO, err error) {
	var v dLogin

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName, loginDocFilter(account), nil, &v,
		)
	}

	if err = withContext(f); err == nil {
		col.toLoginDO(&v, &do)

		return
	}

	if isDocNotExists(err) {
		err = repositories.NewErrorDataNotExists(err)
	}

	return
}

func (col login) toLoginDoc(do *repositories.LoginDO) (bson.M, error) {
	docObj := dLogin{
		Account:     do.Account,
		Info:        do.Info,
		AccessToken: do.AccessToken,
	}

	return genDoc(docObj)
}

func (col login) toLoginDO(u *dLogin, do *repositories.LoginDO) {
	*do = repositories.LoginDO{
		Account:     u.Account,
		Info:        u.Info,
		AccessToken: u.AccessToken,
	}
}
