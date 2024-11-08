package repositoryimpl

import (
	"errors"
	"fmt"

	// "github.com/opensourceways/xihe-server/domain"
	spacedomain "github.com/opensourceways/xihe-server/space/domain"
	"gorm.io/gorm/clause"
)

type projectAdapter struct {
	daoImpl
}

type datasetAdapter struct {
	relatedDaoImpl
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

// func (adapter *projectAdapter) GetByName(owner domain.Account, name domain.ResourceName) (
// 	r spacedomain.Project, err error,
// ) {
// 	//filter
// 	do := projectDO{
// 		Owner: owner.Account(),
// 		Name:  name.ResourceName(),
// 	}

// 	// It must new a new DO, otherwise the sql statement will include duplicate conditions.
// 	//find project
// 	result := projectDO{}
// 	if err := adapter.daoImpl.GetProjectRecord(&do, &result); err != nil {
// 		return spacedomain.Project{}, err
// 	}

// 	id := result.RepoId
// 	query := projectAdapter.daoImpl.dbTag().Where("project_id", id)

// 	//find tags
// 	var tagResults []projectTagsDO
// 	errTag := query.Find(&tagResults).Error

// 	if errTag != nil || len(tagResults) == 0 {
// 		return spacedomain.Project{}, err
// 	}

// 	err = result.toProject(&r)
// 	adapter.getProjectTags(&r, tagResults)

// 	errDataset := b.relatedDaoImpl.

// 		// .db().Where("project_id", id).Find(&datasetResults).Error
// 		// if errDataset != nil || len(tagResults) == 0 {
// 		// 	return spacedomain.Project{}, err
// 		// }
// 		datasetAdapter.getdataset(&r, tagResults)

// 	return r, nil
// }

// func (adapter *projectAdapter) getProjectTags(p *spacedomain.Project, tagResults []projectTagsDO) {
// 	p.Tags = make([]string, 0, len(tagResults))

// 	for _, tagDO := range tagResults {
// 		p.Tags = append(p.Tags, tagDO.TagName)
// 	}
// }
