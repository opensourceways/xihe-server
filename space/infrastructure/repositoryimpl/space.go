package repositoryimpl

import (
	"errors"
	"fmt"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	spacedomain "github.com/opensourceways/xihe-server/space/domain"
	"gorm.io/gorm/clause"
)

type projectAdapter struct {
	daoImpl
}

func (adapter *projectAdapter) Save(v *spacedomain.Project) (spacedomain.Project, error) {
	if v.Id != "" {
		err := errors.New("must be a new project")
		return spacedomain.Project{}, err
	}

	do := toProjectDO(v)
	err := adapter.db().Clauses(clause.Returning{}).Create(&do).Error
	if err != nil {
		return spacedomain.Project{}, err
	}

	doTags := toProjectTagsDO(v)
	for _, doTag := range doTags {
		if err := adapter.dbTag().Clauses(clause.Returning{}).Create(&doTag).Error; err != nil {
			return spacedomain.Project{}, err
		}
	}
	fmt.Printf("==================================v: %+v\n", *v)
	return *v, nil
}

func (adapter *projectAdapter) GetByName(owner domain.Account, name domain.ResourceName) (
	r spacedomain.Project, err error,
) {
	//filter
	do := projectDO{
		Owner: owner.Account(),
		Name:  name.ResourceName(),
	}

	// find project
	result := projectDO{}
	if err := adapter.daoImpl.GetProjectRecord(&do, &result); err != nil {
		return spacedomain.Project{}, err
	}

	id := result.RepoId
	if err = result.toProject(&r); err != nil {
		return spacedomain.Project{}, err
	}

	// find tags
	var tagResults []projectTagsDO
	if errTag := adapter.daoImpl.dbTag().Where("project_id", id).Find(&tagResults).Error; errTag != nil {
		return spacedomain.Project{}, errTag
	}
	adapter.getProjectTags(&r, tagResults)

	// get datasets
	var datasetResults []datasetDO
	if errDataset := adapter.daoImpl.dbDataset().Where("project_id", id).Find(&datasetResults).Error; errDataset != nil {
		return spacedomain.Project{}, errDataset
	}
	adapter.getDataset(&r, datasetResults)

	// get models
	var modelResults []modelDO
	if errModel := adapter.daoImpl.dbModel().Where("project_id", id).Find(&datasetResults).Error; errModel != nil {
		return spacedomain.Project{}, errModel
	}
	adapter.getModel(&r, modelResults)

	return r, nil
}

func (adapter *projectAdapter) getProjectTags(p *spacedomain.Project, tagResults []projectTagsDO) {
	p.Tags = make([]string, 0, len(tagResults))

	for _, tagDO := range tagResults {
		p.Tags = append(p.Tags, tagDO.TagName)
	}
}

func (adapter *projectAdapter) getDataset(p *spacedomain.Project, datasetResult []datasetDO) {
	if len(datasetResult) == 0 {
		return
	}

	relatedDatasets := make(domain.RelatedResources, len(datasetResult))

	for i, dataset := range datasetResult {
		relatedDatasets[i] = domain.ResourceIndex{
			Owner: domain.CreateAccount(dataset.Owner),
			Id:    dataset.DatasetId,
		}
	}

	p.RelatedDatasets = relatedDatasets

}

func (adapter *projectAdapter) getModel(p *spacedomain.Project, modelResult []modelDO) {
	if len(modelResult) == 0 {
		return
	}

	relatedModels := make(domain.RelatedResources, len(modelResult))

	for i, model := range modelResult {
		relatedModels[i] = domain.ResourceIndex{
			Owner: domain.CreateAccount(model.Owner),
			Id:    model.ModelId,
		}
	}

	p.RelatedDatasets = relatedModels

}

func (adapter *projectAdapter) AddRelatedDataset(info *repository.RelatedResourceInfo) error {
	do := toDatasetDO(info)
	fmt.Printf("================================do: %+v\n", do)
	return adapter.dbDataset().Clauses(clause.Returning{}).Create(&do).Error
}
