package repositoryimpl

import (
	"context"

	"github.com/opensourceways/xihe-server/course/domain"
	"github.com/opensourceways/xihe-server/course/domain/repository"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
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

// List
func (impl *courseRepoImpl) FindCourses(opt *repository.CourseListOption) (
	[]domain.CourseSummary, error) {
	var v []DCourse

	f := func(ctx context.Context) error {
		filter := bson.M{}
		if opt.Status != nil {
			filter[fieldStatus] = opt.Status.CourseStatus()
		}
		if opt.Type != nil {
			filter[fieldType] = opt.Type.CourseType()
		}
		if len(opt.CourseIds) > 0 {
			filter[fieldId] = bson.M{
				"$in": opt.CourseIds,
			}
		}

		return impl.cli.GetDocs(ctx, filter, nil, &v)
	}

	if err := withContext(f); err != nil || len(v) == 0 {
		return nil, err
	}

	r := make([]domain.CourseSummary, len(v))
	for i := range v {
		if err := v[i].toCourseSummary(&r[i]); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (impl *courseRepoImpl) FindAssignments(cid string) (
	[]domain.Assignment, error) {
	var v []DCourse

	f := func(ctx context.Context) error {
		filter := impl.docFilter(cid)

		return impl.cli.GetDocs(ctx, filter, nil, &v)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}

		return nil, err
	}
	a := v[0].Assignments
	r := make([]domain.Assignment, len(a))
	for i := range a {
		if err := a[i].toAssignment(&r[i]); err != nil {
			return nil, err
		}
	}
	return r, nil
}

func (impl *courseRepoImpl) FindAssignment(cid string, aid string) (
	a domain.Assignment, err error) {
	var v DCourse

	f := func(ctx context.Context) error {
		filter := bson.M{
			fieldId:   cid,
			fieldAsgs: bson.M{mongoCmdElemMatch: bson.M{fieldId: aid}},
		}

		return impl.cli.GetDoc(ctx, filter, nil, &v)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}
	}

	if v.Assignments == nil || len(v.Assignments) == 0 {
		err = repoerr.NewErrorResourceNotExists(err)
		return
	}

	if err = v.Assignments[0].toAssignment(&a); err != nil {
		return
	}

	return
}
