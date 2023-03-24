package repositoryimpl

import (
	"context"

	"github.com/opensourceways/xihe-server/course/domain"
	"github.com/opensourceways/xihe-server/course/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
)

type workRepoImpl struct {
	cli mongodbClient
}

func NewWorkRepo(m mongodbClient) repository.Work {
	return &workRepoImpl{m}
}

func (impl *workRepoImpl) GetWork(cid string, user types.Account, asgId string, status domain.WorkStatus) (
	w domain.Work, err error,
) {
	var v []DCourseWork

	f := func(ctx context.Context) error {
		filter := bson.M{
			fieldCourseId: cid,
			fieldAsgId:    asgId,
			fieldAccount:  user.Account(),
		}
		if status != nil {
			filter[fieldStatus] = status.WorkStatus()
		}

		return impl.cli.GetDocs(ctx, filter, nil, &v)
	}

	if err = withContext(f); err != nil {
		return
	}
	if len(v) == 0 {
		err = repoerr.NewErrorResourceNotExists(err)
		return
	}

	r := make([]domain.Work, len(v))
	if err = v[0].toCourseWork(&r[0]); err != nil {
		return

	}
	return r[0], nil

}
