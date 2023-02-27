package repositoryimpl

import (
	"context"

	"github.com/opensourceways/xihe-server/course/domain"
	"github.com/opensourceways/xihe-server/course/domain/repository"
	"go.mongodb.org/mongo-driver/bson"

	repoerr "github.com/opensourceways/xihe-server/domain/repository"
)

func NewCourseRepo(m mongodbClient) repository.Course {
	return &courseRepoImpl{m}
}

type courseRepoImpl struct {
	cli mongodbClient
}

func (impl *courseRepoImpl) FindCourse(cid string) (
	c domain.Course,
	err error,
) {
	var v DCourse

	f := func(ctx context.Context) error {
		filter := impl.docFilter(cid)

		return impl.cli.GetDoc(ctx, filter, nil, &v)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}

		return
	}

	if err = v.toCourse(&c); err != nil {
		return
	}

	return
}

func (impl *courseRepoImpl) docFilter(cid string) bson.M {
	return bson.M{
		fieldId: cid,
	}
}
