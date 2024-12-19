package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	spacedomain "github.com/opensourceways/xihe-server/space/domain"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
)

type ReverselyRelatedResourceInfo = domain.ReverselyRelatedResourceInfo

// dataset
type DatasetMessageService interface {
	AddRelatedProject(*ReverselyRelatedResourceInfo) error
	RemoveRelatedProject(*ReverselyRelatedResourceInfo) error

	AddRelatedModel(*ReverselyRelatedResourceInfo) error
	RemoveRelatedModel(*ReverselyRelatedResourceInfo) error

	AddLike(*domain.ResourceIndex) error
	RemoveLike(*domain.ResourceIndex) error

	IncreaseDownload(*domain.ResourceIndex) error
}

func NewDatasetMessageService(repo repository.Dataset) DatasetMessageService {
	return datasetMessageService{
		repo: repo,
	}
}

type datasetMessageService struct {
	repo repository.Dataset
}

func (s datasetMessageService) AddRelatedProject(info *ReverselyRelatedResourceInfo) error {
	return s.repo.AddRelatedProject(info)
}

func (s datasetMessageService) RemoveRelatedProject(info *ReverselyRelatedResourceInfo) error {
	return s.repo.RemoveRelatedProject(info)
}

func (s datasetMessageService) AddRelatedModel(info *ReverselyRelatedResourceInfo) error {
	return s.repo.AddRelatedModel(info)
}

func (s datasetMessageService) RemoveRelatedModel(info *ReverselyRelatedResourceInfo) error {
	return s.repo.RemoveRelatedModel(info)
}

func (s datasetMessageService) AddLike(r *domain.ResourceIndex) error {
	return s.repo.AddLike(r)
}

func (s datasetMessageService) RemoveLike(r *domain.ResourceIndex) error {
	return s.repo.RemoveLike(r)
}

func (s datasetMessageService) IncreaseDownload(index *domain.ResourceIndex) error {
	return s.repo.IncreaseDownload(index)
}

// model
type ModelMessageService interface {
	AddRelatedProject(*ReverselyRelatedResourceInfo) error
	RemoveRelatedProject(*ReverselyRelatedResourceInfo) error

	AddRelatedDataset(*ReverselyRelatedResourceInfo) error
	RemoveRelatedDataset(*ReverselyRelatedResourceInfo) error

	AddLike(*domain.ResourceIndex) error
	RemoveLike(*domain.ResourceIndex) error

	IncreaseDownload(*domain.ResourceIndex) error
}

type modelMessageService struct {
	repo repository.Model
}

func NewModelMessageService(repo repository.Model) ModelMessageService {
	return modelMessageService{
		repo: repo,
	}
}

func (s modelMessageService) AddRelatedProject(info *ReverselyRelatedResourceInfo) error {
	return s.repo.AddRelatedProject(info)
}

func (s modelMessageService) RemoveRelatedProject(info *ReverselyRelatedResourceInfo) error {
	return s.repo.RemoveRelatedProject(info)
}

func (s modelMessageService) AddRelatedDataset(info *ReverselyRelatedResourceInfo) error {
	// this case will not happen
	return nil
}

func (s modelMessageService) RemoveRelatedDataset(info *ReverselyRelatedResourceInfo) error {
	m, err := s.repo.Get(info.Resource.Owner, info.Resource.Id)
	if err != nil {
		return err
	}

	index := info.Promoter

	if !m.RelatedDatasets.Has(index) {
		return nil
	}

	param := repository.RelatedResourceInfo{
		ResourceToUpdate: s.toResourceToUpdate(&m),
		RelatedResource:  *index,
	}

	return s.repo.RemoveRelatedDataset(&param)
}

func (s modelMessageService) AddLike(r *domain.ResourceIndex) error {
	return s.repo.AddLike(r)
}

func (s modelMessageService) RemoveLike(r *domain.ResourceIndex) error {
	return s.repo.RemoveLike(r)
}

func (s modelMessageService) toResourceToUpdate(m *domain.Model) repository.ResourceToUpdate {
	return repository.ResourceToUpdate{
		Owner:     m.Owner,
		Id:        m.Id,
		Version:   m.Version,
		UpdatedAt: m.UpdatedAt,
	}
}

func (s modelMessageService) IncreaseDownload(index *domain.ResourceIndex) error {
	return s.repo.IncreaseDownload(index)
}

// project
type ProjectMessageService interface {
	AddRelatedModel(*ReverselyRelatedResourceInfo) error
	RemoveRelatedModel(*ReverselyRelatedResourceInfo) error

	AddRelatedDataset(*ReverselyRelatedResourceInfo) error
	RemoveRelatedDataset(*ReverselyRelatedResourceInfo) error

	AddLike(*domain.ResourceIndex) error
	RemoveLike(*domain.ResourceIndex) error

	IncreaseFork(*domain.ResourceIndex) error
	IncreaseDownload(*domain.ResourceIndex) error
}

type projectMessageService struct {
	repoPg spacerepo.ProjectPg
}

func NewProjectMessageService(repo spacerepo.ProjectPg, repoPg spacerepo.ProjectPg) ProjectMessageService {
	return projectMessageService{
		repoPg: repoPg,
	}
}

func (s projectMessageService) AddRelatedModel(info *ReverselyRelatedResourceInfo) error {
	// this case will not happen
	return nil
}

func (s projectMessageService) RemoveRelatedModel(info *ReverselyRelatedResourceInfo) error {
	p, err := s.repoPg.Get(info.Resource.Owner, info.Resource.Id)
	if err != nil {
		return err
	}

	index := info.Promoter

	if !p.RelatedModels.Has(index) {
		return nil
	}

	param := repository.RelatedResourceInfo{
		ResourceToUpdate: s.toResourceToUpdate(&p),
		RelatedResource:  *index,
	}

	return s.repoPg.RemoveRelatedModel(&param)
}

func (s projectMessageService) AddRelatedDataset(info *ReverselyRelatedResourceInfo) error {
	// this case will not happen
	return nil
}

func (s projectMessageService) RemoveRelatedDataset(info *ReverselyRelatedResourceInfo) error {
	p, err := s.repoPg.Get(info.Resource.Owner, info.Resource.Id)
	if err != nil {
		return err
	}

	index := info.Promoter

	if !p.RelatedDatasets.Has(index) {
		return nil
	}

	param := repository.RelatedResourceInfo{
		ResourceToUpdate: s.toResourceToUpdate(&p),
		RelatedResource:  *index,
	}

	return s.repoPg.RemoveRelatedDataset(&param)
}

func (s projectMessageService) AddLike(r *domain.ResourceIndex) error {
	return s.repoPg.AddLike(r)
}

func (s projectMessageService) RemoveLike(r *domain.ResourceIndex) error {
	return s.repoPg.RemoveLike(r)
}

func (s projectMessageService) IncreaseFork(index *domain.ResourceIndex) error {
	return s.repoPg.IncreaseFork(index)
}

func (s projectMessageService) IncreaseDownload(index *domain.ResourceIndex) error {
	return s.repoPg.IncreaseDownload(index)
}

func (s projectMessageService) toResourceToUpdate(p *spacedomain.Project) repository.ResourceToUpdate {
	return repository.ResourceToUpdate{
		Owner:     p.Owner,
		Id:        p.Id,
		Version:   p.Version,
		UpdatedAt: p.UpdatedAt,
	}
}
