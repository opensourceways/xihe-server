package app

import (
	"github.com/opensourceways/xihe-server/domain"
)

type CreateProjectCmd struct {
	Name      domain.ProjName
	Desc      domain.ProjDesc
	Type      domain.RepoType
	CoverId   string
	Protocol  domain.ProtocolName
	Training  domain.TrainingSDK
	Inference domain.InferenceSDK
}

func (cmd *CreateProjectCmd) toProject() domain.Project {
	return domain.Project{
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

	p := cmd.toProject()

	if err := s.repo.Save(p); err != nil {
		return dto, err
	}

	return s.toProjectDTO(p), nil
}

func (s createProjectService) toProjectDTO(p domain.Project) ProjectDTO {
	return ProjectDTO{}
}
