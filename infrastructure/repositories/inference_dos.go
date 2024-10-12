package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
)

type InferenceIndexDO struct {
	Id         string
	LastCommit string
	Project    ResourceIndexDO
}

func (impl inference) toInferenceIndexDO(index *domain.InferenceIndex) InferenceIndexDO {
	return InferenceIndexDO{
		Id:         index.Id,
		LastCommit: index.LastCommit,
		Project:    ToResourceIndexDO(&index.Project),
	}
}

type InferenceDetailDO = domain.InferenceDetail

type InferenceSummaryDO struct {
	Id string

	InferenceDetailDO
}

type InferenceDO struct {
	Id           string
	ProjectId    string
	LastCommit   string
	ProjectName  string
	ProjectOwner string

	InferenceDetailDO
}

func (impl inference) toInferenceDO(obj *domain.Inference) InferenceDO {
	return InferenceDO{
		Id:           obj.Id,
		ProjectId:    obj.Project.Id,
		LastCommit:   obj.LastCommit,
		ProjectName:  obj.ProjectName.ResourceName(),
		ProjectOwner: obj.Project.Owner.Account(),
	}
}
