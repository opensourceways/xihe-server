package repositoryimpl

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/competition/domain"
	"github.com/opensourceways/xihe-server/competition/domain/repository"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
)

func NewWorkRepo(m mongodbClient) repository.Work {
	return workRepoImpl{m}
}

type workRepoImpl struct {
	cli mongodbClient
}

func (impl workRepoImpl) docFilter(index *domain.WorkIndex) bson.M {
	return bson.M{
		fieldCid: index.CompetitionId,
		fieldPid: index.PlayerId,
	}
}

func (impl workRepoImpl) SaveWork(w *domain.Work) error {
	doc, err := genDoc(dWork{
		CompetitionId: w.CompetitionId,
		PlayerId:      w.PlayerId,
		PlayerName:    w.PlayerName,
		Final:         []dSubmission{},
		Preliminary:   []dSubmission{},
	})
	if err != nil {
		return err
	}
	doc[fieldVersion] = 0

	f := func(ctx context.Context) error {
		_, err := impl.cli.NewDocIfNotExist(
			ctx, impl.docFilter(&w.WorkIndex), doc,
		)

		return err
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocExists(err) {
			err = repoerr.NewErrorDuplicateCreating(err)
		}
	}

	return err
}

func (impl workRepoImpl) SaveRepo(w *domain.Work, version int) error {
	f := func(ctx context.Context) error {
		return impl.cli.UpdateDoc(
			ctx, impl.docFilter(&w.WorkIndex),
			bson.M{fieldRepo: w.Repo}, mongoCmdSet, version,
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

func (impl workRepoImpl) AddSubmission(
	w *domain.Work, cs *domain.PhaseSubmission, version int,
) error {
	doc, err := genDoc(dSubmission{
		Id:       cs.Id,
		Status:   cs.Status,
		OBSPath:  cs.OBSPath,
		SubmitAt: cs.SubmitAt,
	})
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		field := fieldPreliminary
		if cs.Phase.IsFinal() {
			field = fieldFinal
		}

		return impl.cli.UpdateDoc(
			ctx, impl.docFilter(&w.WorkIndex),
			bson.M{field: doc}, mongoCmdPush, version,
		)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorConcurrentUpdating(err)
		}
	}

	return err
}

func (impl workRepoImpl) SaveSubmission(
	w *domain.Work, submission *domain.PhaseSubmission,
) error {
	doc, err := genDoc(dSubmission{
		Status: submission.Status,
		Score:  float64(submission.Score),
	})
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		field := fieldPreliminary
		if submission.Phase.IsFinal() {
			field = fieldFinal
		}

		_, err := impl.cli.ModifyArrayElem(
			ctx, field, impl.docFilter(&w.WorkIndex),
			bson.M{fieldId: submission.Id}, doc, mongoCmdSet,
		)

		return err
	}

	return withContext(f)
}

func (impl workRepoImpl) FindWork(index domain.WorkIndex, Phase domain.CompetitionPhase) (
	w domain.Work, version int, err error,
) {
	var v dWork

	f := func(ctx context.Context) error {
		filter := impl.docFilter(&index)

		project := bson.M{}
		if Phase.IsPreliminary() {
			project[fieldFinal] = 0
		} else {
			project[fieldPreliminary] = 0
		}

		return impl.cli.GetDoc(ctx, filter, project, &v)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorResourceNotExists(err)
		}
	} else {
		version = v.Version
		v.toWork(&w)
	}

	return
}

func (impl workRepoImpl) FindWorks(cid string) (ws []domain.Work, err error) {
	var v []dWork

	f := func(ctx context.Context) error {
		return impl.cli.GetDocs(ctx, bson.M{fieldCid: cid}, nil, &v)
	}

	if err = withContext(f); err != nil || len(v) == 0 {
		return
	}

	ws = make([]domain.Work, len(v))
	for i := range v {
		v[i].toWork(&ws[i])
	}

	return
}
