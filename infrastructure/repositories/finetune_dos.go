package repositories

import "github.com/opensourceways/xihe-server/domain"

type FinetuneIndexDO struct {
	Id    string
	Owner string
}

func (impl finetuneImpl) toFinetuneIndexDO(obj *domain.FinetuneIndex) FinetuneIndexDO {
	return FinetuneIndexDO{
		Id:    obj.Id,
		Owner: obj.Owner.Account(),
	}
}

type UserFinetuneDO struct {
	FinetuneIndexDO

	FinetuneConfigDO

	CreatedAt int64
}

func (impl finetuneImpl) toUserFinetuneDO(
	user domain.Account, obj *domain.Finetune, do *UserFinetuneDO,
) {
	p := obj.Param

	*do = UserFinetuneDO{
		FinetuneIndexDO: FinetuneIndexDO{
			Id:    obj.Id,
			Owner: user.Account(),
		},
		FinetuneConfigDO: FinetuneConfigDO{
			Name:            obj.Name.FinetuneName(),
			Model:           p.Model(),
			Task:            p.Task(),
			Hyperparameters: p.Hyperparameters(),
		},
		CreatedAt: obj.CreatedAt,
	}
}

type FinetuneConfigDO struct {
	Name            string
	Task            string
	Model           string
	Hyperparameters map[string]string
}

func (do *FinetuneConfigDO) toFinetuneConfig(cfg *domain.FinetuneConfig) (err error) {
	if cfg.Name, err = domain.NewFinetuneName(do.Name); err != nil {
		return
	}

	cfg.Param, err = domain.NewFinetuneParameter(
		do.Model, do.Task, do.Hyperparameters,
	)

	return
}

type FinetuneJobDO = domain.FinetuneJob
type FinetuneJobInfoDO = domain.FinetuneJobInfo
type FinetuneJobDetailDO = domain.FinetuneJobDetail

type FinetuneDetailDO struct {
	Id        string
	CreatedAt int64

	FinetuneConfigDO

	Job       FinetuneJobInfoDO
	JobDetail FinetuneJobDetailDO
}

func (do *FinetuneDetailDO) toUserFinetune(obj *domain.Finetune) (err error) {
	err = do.FinetuneConfigDO.toFinetuneConfig(&obj.FinetuneConfig)
	if err != nil {
		return
	}

	obj.Id = do.Id
	obj.Job = do.Job
	obj.JobDetail = do.JobDetail
	obj.CreatedAt = do.CreatedAt

	return
}

type FinetuneSummaryDO struct {
	Id        string
	Name      string
	CreatedAt int64

	FinetuneJobDetailDO
}

func (do *FinetuneSummaryDO) toFinetuneSummary(obj *domain.FinetuneSummary) (err error) {
	if obj.Name, err = domain.NewFinetuneName(do.Name); err != nil {
		return
	}

	obj.Id = do.Id
	obj.CreatedAt = do.CreatedAt
	obj.FinetuneJobDetail = do.FinetuneJobDetailDO

	return
}

type UserFinetunesDO struct {
	Expiry int64

	Data []FinetuneSummaryDO
}
