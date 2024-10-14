package repositories

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type InferenceMapper interface {
	Insert(*InferenceDO, int) (string, error)
	Get(*InferenceIndexDO) (InferenceSummaryDO, error)
	UpdateDetail(*InferenceIndexDO, *InferenceDetailDO) error
	List(*ResourceIndexDO, string) ([]InferenceSummaryDO, int, error)
}

func NewInferenceRepository(mapper InferenceMapper) repository.Inference {
	return inference{mapper}
}

type inference struct {
	mapper InferenceMapper
}

func (impl inference) Save(ut *domain.Inference, version int) (string, error) {
	if ut.Id != "" {
		return "", errors.New("must be a new project")
	}

	do := impl.toInferenceDO(ut)

	v, err := impl.mapper.Insert(&do, version)
	if err != nil {
		return "", convertError(err)
	}

	return v, nil
}

func (impl inference) FindInstance(info *domain.InferenceIndex) (
	r repository.InferenceSummary, err error,
) {
	index := impl.toInferenceIndexDO(info)
	v, err := impl.mapper.Get(&index)
	if err != nil {
		err = convertError(err)

		return
	}

	r.InferenceDetail = v.InferenceDetailDO
	r.Id = v.Id

	return
}

func (impl inference) FindInstances(info *domain.ResourceIndex, lastCommit string) (
	r []repository.InferenceSummary, version int, err error,
) {
	index := ToResourceIndexDO(info)
	v, version, err := impl.mapper.List(&index, lastCommit)
	if err != nil {
		err = convertError(err)

		return
	}

	r = make([]repository.InferenceSummary, len(v))

	for i := range v {
		r[i].InferenceDetail = v[i].InferenceDetailDO
		r[i].Id = v[i].Id
	}

	return
}

func (impl inference) UpdateDetail(
	info *domain.InferenceIndex, detail *domain.InferenceDetail,
) error {
	index := impl.toInferenceIndexDO(info)

	if err := impl.mapper.UpdateDetail(&index, detail); err != nil {
		return convertError(err)
	}

	return nil
}
