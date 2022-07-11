package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type CreateProjectCmd struct {
	Owner     string
	Name      domain.ProjName
	Desc      domain.ProjDesc
	Type      domain.RepoType
	CoverId   domain.CoverId
	Protocol  domain.ProtocolName
	Training  domain.TrainingSDK
	Inference domain.InferenceSDK
}

func (cmd *CreateProjectCmd) validate() error {
	b := cmd.Owner != "" &&
		cmd.Name != nil &&
		cmd.Type != nil &&
		cmd.CoverId != nil &&
		cmd.Protocol != nil &&
		cmd.Training != nil &&
		cmd.Inference != nil

	if !b {
		return errors.New("invalid cmd of creating project")
	}

	return nil
}

func (cmd *CreateProjectCmd) toProject() domain.Project {
	return domain.Project{
		Owner:     cmd.Owner,
		Name:      cmd.Name,
		Desc:      cmd.Desc,
		Type:      cmd.Type,
		CoverId:   cmd.CoverId,
		Protocol:  cmd.Protocol,
		Training:  cmd.Training,
		Inference: cmd.Inference,
	}
}

type ProjectDTO struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Type      string `json:"type"`
	CoverId   string `json:"cover_id"`
	Protocol  string `json:"protocol"`
	Training  string `json:"training"`
	Inference string `json:"inference"`
}

type CreateProjectService interface {
	Create(cmd CreateProjectCmd) (ProjectDTO, error)
}

func NewCreateProjectService(repo repository.Project) CreateProjectService {
	return createProjectService{repo}
}

type createProjectService struct {
	repo repository.Project
}

func (s createProjectService) Create(cmd CreateProjectCmd) (dto ProjectDTO, err error) {
	if err = cmd.validate(); err != nil {
		return
	}

	v, err := s.repo.Save(cmd.toProject())
	if err != nil {
		return
	}

	dto.Id = v.Id

	return
}
