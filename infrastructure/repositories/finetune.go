package repositories

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type FinetuneMapper interface {
	Delete(*FinetuneIndexDO) error
	Insert(*UserFinetuneDO, int) (string, error)
	Get(*FinetuneIndexDO) (FinetuneDetailDO, error)
	List(user string) (UserFinetunesDO, int, error)

	GetJob(*FinetuneIndexDO) (FinetuneJobDO, error)
	UpdateJobInfo(*FinetuneIndexDO, *FinetuneJobInfoDO) error
	UpdateJobDetail(*FinetuneIndexDO, *FinetuneJobDetailDO) error
}

func NewFinetuneRepository(mapper FinetuneMapper) repository.Finetune {
	return finetuneImpl{mapper}
}

type finetuneImpl struct {
	mapper FinetuneMapper
}

func (impl finetuneImpl) Save(user domain.Account, obj *domain.Finetune, version int) (
	string, error,
) {
	if obj.Id != "" {
		return "", errors.New("must be a new finetune")
	}

	do := new(UserFinetuneDO)
	impl.toUserFinetuneDO(user, obj, do)

	v, err := impl.mapper.Insert(do, version)
	if err != nil {
		err = convertError(err)
	}

	return v, err
}

func (impl finetuneImpl) Get(index *domain.FinetuneIndex) (
	obj domain.Finetune, err error,
) {
	do := impl.toFinetuneIndexDO(index)

	v, err := impl.mapper.Get(&do)
	if err != nil {
		err = convertError(err)
	} else {
		err = v.toUserFinetune(&obj)
	}

	return
}

func (impl finetuneImpl) Delete(index *domain.FinetuneIndex) error {
	do := impl.toFinetuneIndexDO(index)

	if err := impl.mapper.Delete(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl finetuneImpl) List(user domain.Account) (
	r repository.UserFinetunes, err error,
) {
	v, version, err := impl.mapper.List(user.Account())
	if err != nil {
		err = convertError(err)

		return
	}

	r.Version = version
	r.Expiry = v.Expiry

	if len(v.Data) == 0 {
		return
	}

	data := make([]domain.FinetuneSummary, len(v.Data))
	for i := range v.Data {
		if err = v.Data[i].toFinetuneSummary(&data[i]); err != nil {
			return
		}
	}
	r.Data = data

	return
}

func (impl finetuneImpl) GetJob(index *domain.FinetuneIndex) (domain.FinetuneJob, error) {
	do := impl.toFinetuneIndexDO(index)

	v, err := impl.mapper.GetJob(&do)
	if err != nil {
		err = convertError(err)
	}

	return v, err
}

func (impl finetuneImpl) SaveJob(
	index *domain.FinetuneIndex, job *domain.FinetuneJobInfo,
) error {
	do := impl.toFinetuneIndexDO(index)

	if err := impl.mapper.UpdateJobInfo(&do, job); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl finetuneImpl) UpdateJobDetail(
	index *domain.FinetuneIndex, job *domain.FinetuneJobDetail,
) error {
	do := impl.toFinetuneIndexDO(index)

	if err := impl.mapper.UpdateJobDetail(&do, job); err != nil {
		return convertError(err)
	}

	return nil
}
