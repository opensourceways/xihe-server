package repositoryimpl

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/cloud/domain"
	"github.com/opensourceways/xihe-server/cloud/domain/repository"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
)

func NewCloudRepo(m mongodbClient) repository.Cloud {
	return &cloudRepoImpl{m}
}

type cloudRepoImpl struct {
	cli mongodbClient
}

func (impl *cloudRepoImpl) ListCloudConf() (conf []domain.CloudConf, err error) {
	var v []DCloudConf

	f := func(ctx context.Context) error {
		return impl.cli.GetDocs(ctx, nil, nil, &v)
	}

	if err = withContext(f); err != nil {
		return
	}

	conf = make([]domain.CloudConf, len(v))
	for i := range v {
		if err = v[i].toCloudConf(&conf[i]); err != nil {
			return
		}
	}

	return
}

func (impl *cloudRepoImpl) GetCloudConf(cid string) (conf domain.CloudConf, err error) {
	var v DCloudConf

	f := func(ctx context.Context) error {
		filter := impl.docIdFilter(cid)

		return impl.cli.GetDoc(ctx, filter, nil, &v)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}

		return
	}

	err = v.toCloudConf(&conf)

	return
}

func (impl *cloudRepoImpl) docIdFilter(id string) bson.M {
	return bson.M{
		fieldId: id,
	}
}
