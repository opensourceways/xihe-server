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
		fmt.Printf("=========================doTag: %+v\n", doTag)
		if err := adapter.dbTag().Clauses(clause.Returning{}).Create(&doTag).Error; err != nil {
			return spacedomain.Project{}, err
		}
	}
	return *v, nil
}

// func (adapter *projectAdapter) GetByName(owner domain.Account, name domain.ResourceName) (
// 	r spacedomain.Project, err error,
// ) {
// 	do := projectDO{
// 		Owner: owner.Account(),
// 		Name:  name.ResourceName(),
// 	}

// 	// It must new a new DO, otherwise the sql statement will include duplicate conditions.
// 	result := projectDO{}
// 	if err := adapter.daoImpl.GetProjectRecord(&do, &result); err != nil {
// 		return spacedomain.Project{}, err
// 	}

// 	id := result.RepoId

// 	tagdo := projectTagsDO{
// 		projectId: id,
// 	}

// 	tagResults := []projectTagsDO{}

// 	err = result.toProject(&r)
// 	return
// }
