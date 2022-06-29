package app

import (
	"github.com/opensourceways/xihe-server/domain"
)

type CreateProjectCmd struct {
	Name      string
	Desc      string
	Type      string
	CoverId   string
	Protocol  string
	Training  string
	Inference string
}

func (cmd *CreateProjectCmd) toDomainProject() (p domain.Project, err error) {
	p.NewOne = true

	n, err := domain.NewProjName(cmd.Name)
	if err != nil {
		return
	}
	p.Name = n

	t, err := domain.NewRepoType(cmd.Type)
	if err != nil {
		return
	}
	p.Type = t

	d, err := domain.NewProjDesc(cmd.Desc)
	if err != nil {
		return
	}
	p.Desc = d

	// TODO: check cover id
	p.CoverId = cmd.CoverId

	pv, err := domain.NewProtocolName(cmd.Protocol)
	if err != nil {
		return
	}
	p.Protocol = pv

	tv, err := domain.NewTrainingSDK(cmd.Training)
	if err != nil {
		return
	}
	p.Training = tv

	iv, err := domain.NewInferenceSDK(cmd.Inference)
	if err != nil {
		return
	}
	p.Inference = iv

	err = p.Validate()

	return
}

type ProjectDTO struct {
	Name      string
	Desc      string
	Type      string
	CoverId   string
	Protocol  string
	Training  string
	Inference string
}

type CreateProjectService interface {
	Create(userId string, cmd CreateProjectCmd) (ProjectDTO, error)
}

func NewCreateProjectService(repo ProjectRepository) CreateProjectService {
	return createProjectService{repo}
}

type createProjectService struct {
	repo ProjectRepository
}

func (s createProjectService) Create(userId string, cmd CreateProjectCmd) (ProjectDTO, error) {
	dto := ProjectDTO{}

	p, err := cmd.toDomainProject()
	if err != nil {
		return dto, err
	}

	if err := s.repo.Save(p); err != nil {
		return dto, err
	}

	return s.toProjectDTO(p), nil
}

func (s createProjectService) toProjectDTO(p domain.Project) ProjectDTO {
	return ProjectDTO{}
}
