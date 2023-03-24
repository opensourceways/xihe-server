package repositoryimpl

import (
	"context"

	"github.com/opensourceways/xihe-server/course/domain"
	"github.com/opensourceways/xihe-server/course/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
)

func NewPlayerRepo(m mongodbClient) repository.Player {
	return &playerRepoImpl{m}
}

type playerRepoImpl struct {
	cli mongodbClient
}

func (impl *playerRepoImpl) FindPlayer(cid string, user types.Account) (p repository.PlayerVersion, err error) {
	var v DCoursePlayer

	f := func(ctx context.Context) error {
		filter := impl.docFilterFindPlayer(cid, user.Account())

		return impl.cli.GetDoc(ctx, filter, nil, &v)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}

		return
	}
	p.Version = v.Version

	if err = v.toPlayerNoStudent(&p); err != nil {
		return
	}

	return
}

func (impl *playerRepoImpl) SavePlayer(p *domain.Player) (err error) {
	doc, err := impl.genPlayerDoc(p)
	if err != nil {
		return
	}
	doc[fieldVersion] = 1

	f := func(ctx context.Context) error {
		_, err := impl.cli.NewDocIfNotExist(
			ctx, bson.M{
				fieldAccount:  p.Account.Account(),
				fieldCourseId: p.CourseId,
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

func (impl *playerRepoImpl) genPlayerDoc(p *domain.Player) (bson.M, error) {
	obj := DCoursePlayer{
		Id:        p.Id,
		CourseId:  p.CourseId,
		Name:      p.Account.Account(),
		CreatedAt: p.CreatedAt.CourseTime(),
	}

	return genDoc(obj)
}

// Player Count
func (impl *playerRepoImpl) PlayerCount(cid string) (int, error) {
	var v []struct {
		Total int `bson:"total"`
	}

	f := func(ctx context.Context) error {

		pipeline := bson.A{
			bson.M{
				mongoCmdMatch: bson.M{
					fieldCourseId: bson.M{mongoCmdEqual: cid},
				},
			},
			bson.M{mongoCmdCount: "total"},
		}

		cursor, err := impl.cli.Collection().Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, &v)
	}

	if err := withContext(f); err != nil || len(v) == 0 {
		return 0, err
	}

	return v[0].Total, nil
}

func (impl *playerRepoImpl) docFilterFindPlayer(cid, account string) bson.M {

	return bson.M{
		fieldCourseId: cid,
		fieldAccount:  account,
	}
}

func (impl *playerRepoImpl) SaveRepo(course_id string, a *domain.CourseProject, version int) error {
	f := func(ctx context.Context) error {

		return impl.cli.UpdateDoc(
			ctx,
			impl.docFilterFindPlayer(course_id, a.Owner.Account()),
			bson.M{fieldRepo: a.RepoRouting}, mongoCmdSet, version,
		)
	}

	err := withContext(f)

	if err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorConcurrentUpdating(err)
		}
	}

	return err
}

func (impl *playerRepoImpl) docFilter(account string) bson.M {
	return bson.M{
		fieldAccount: account,
	}
}

func (impl *playerRepoImpl) FindCoursesUserApplied(u types.Account) (
	cs []string, err error) {
	var v []DCoursePlayer

	f := func(ctx context.Context) error {
		filter := impl.docFilter(u.Account())
		return impl.cli.GetDocs(ctx, filter, nil, &v)
	}

	if err = withContext(f); err != nil || len(v) == 0 {
		return
	}

	cs = make([]string, len(v))
	for i := range v {
		cs[i] = v[i].CourseId
	}

	return
}
