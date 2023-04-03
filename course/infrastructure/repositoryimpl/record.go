package repositoryimpl

import (
	"context"

	"github.com/opensourceways/xihe-server/course/domain"
	"github.com/opensourceways/xihe-server/course/domain/repository"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
)

func NewRecordRepo(m mongodbClient) repository.Record {
	return &recordRepoImpl{m}
}

type recordRepoImpl struct {
	cli mongodbClient
}

func (impl *recordRepoImpl) FindPlayRecord(r *domain.Record) (a repository.RecordVersion, err error) {
	var v DCourseRecord

	f := func(ctx context.Context) error {
		filter := impl.docFilterFindRecord(r)

		return impl.cli.GetDoc(ctx, filter, nil, &v)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}

		return
	}
	a.Version = v.Version

	if err = v.toRecord(&a); err != nil {
		return
	}

	return
}

func (impl *recordRepoImpl) AddPlayRecord(r *domain.Record) (err error) {

	doc, err := impl.genRecordDoc(r)
	if err != nil {
		return
	}
	doc[fieldVersion] = 1

	f := func(ctx context.Context) error {
		_, err := impl.cli.NewDocIfNotExist(
			ctx, bson.M{
				fieldAccount:   r.User.Account(),
				fieldCourseId:  r.Cid,
				fieldSectionId: r.SectionId,
				fieldLessonId:  r.LessonId,
				fieldPointId:   r.PointId,
			}, doc,
		)
		return err
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocExists(err) {
			err = repoerr.NewErrorDuplicateCreating(err)
		}

		return
	}
	return
}

func (impl *recordRepoImpl) genRecordDoc(p *domain.Record) (bson.M, error) {
	obj := DCourseRecord{
		CourseId:    p.Cid,
		Account:     p.User.Account(),
		SectionId:   p.SectionId.SectionId(),
		LessonId:    p.LessonId.LessonId(),
		PointId:     p.PointId,
		PlayCount:   0,
		FinishCount: 0,
	}

	return genDoc(obj)
}

func (impl *recordRepoImpl) docFilterFindRecord(r *domain.Record) bson.M {
	return bson.M{
		fieldCourseId:  r.Cid,
		fieldAccount:   r.User.Account(),
		fieldSectionId: r.SectionId.SectionId(),
		fieldLessonId:  r.LessonId.LessonId(),
		fieldPointId:   r.PointId,
	}
}

func (impl *recordRepoImpl) UpdatePlayRecord(r *domain.Record, version int) (err error) {
	f := func(ctx context.Context) error {

		return impl.cli.UpdateIncDoc(
			ctx,
			impl.docFilterFindRecord(r),
			bson.M{fieldPlayCount: r.PlayCount, fieldFinishCount: r.FinishCount},
			version,
		)
	}

	err = withContext(f)

	if err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorConcurrentUpdating(err)
		}
	}

	return
}
