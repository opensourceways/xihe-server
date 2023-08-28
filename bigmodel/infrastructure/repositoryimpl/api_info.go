package repositoryimpl

import (
	"context"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	"github.com/opensourceways/xihe-server/bigmodel/domain/repository"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
)

func NewApiInfo(m mongodbClient) repository.ApiInfo {
	return &apiInfoRepoImpl{m}
}

type apiInfoRepoImpl struct {
	cli mongodbClient
}

func (impl *apiInfoRepoImpl) GetApiInfo(m domain.ModelName) (u domain.ApiInfo, err error) {
	v := new(dApiInfo)

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(
			ctx,
			bson.M{fieldId: m.ModelName()},
			nil,
			&v,
		)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}
		return
	}

	err = v.toApiInfo(&u)
	return
}
