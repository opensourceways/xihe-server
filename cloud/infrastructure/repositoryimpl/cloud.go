package repositoryimpl

import (
	"context"

	"github.com/opensourceways/xihe-server/cloud/domain"
	"github.com/opensourceways/xihe-server/cloud/domain/repository"
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
		v[i].toCloudConf(&conf[i])
	}

	return
}
