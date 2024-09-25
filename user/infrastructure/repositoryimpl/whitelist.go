package repositoryimpl

import (
	"context"

	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/user/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
)

func NewWhiteListRepo(m mongodbClient) repository.WhiteList {
	return &userWhiteListImpl{m}
}

type userWhiteListImpl struct {
	cli mongodbClient
}

func (impl *userWhiteListImpl) GetWhiteListInfo(
	account types.Account, wtype string,
) (u domain.WhiteListInfo, err error) {
	var v DWhiteListInfo

	f := func(ctx context.Context) error {
		filter := bson.M{
			fieldAccount:       account.Account(),
			fieldWhiteListType: wtype,
		}

		return impl.cli.GetDoc(ctx, filter, nil, &v)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = nil
		}

		return
	}

	if err = v.toWhiteListInfo(&u); err != nil {
		return
	}

	return
}

func (impl *userWhiteListImpl) GetWhiteListInfoItems(
	accout types.Account, whitelistTypeList []string,
) ([]domain.WhiteListInfo, error) {
	var v []DWhiteListInfo

	f := func(ctx context.Context) error {
		filter := bson.M{
			fieldAccount:       accout.Account(),
			fieldWhiteListType: bson.M{"$in": whitelistTypeList},
		}

		return impl.cli.GetDocs(ctx, filter, nil, &v)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			return nil, nil
		}

		return nil, err
	}

	whitelistItems := make([]domain.WhiteListInfo, len(v))
	for i := range v {
		if err := v[i].toWhiteListInfo(&whitelistItems[i]); err != nil {
			return nil, err
		}
	}

	return whitelistItems, nil
}
