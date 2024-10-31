package repositoryimpl

import (
	"errors"

	spacedomain "github.com/opensourceways/xihe-server/space/domain"
	"gorm.io/gorm/clause"
)

type projectAdapter struct {
	daoImpl
}

func (adapter *projectAdapter) Save(p *spacedomain.Project) (spacedomain.Project, error) {
	if p.Id != "" {
		err := errors.New("must be a new project")
		return spacedomain.Project{}, err
	}

	// 使用 GORM 创建记录
	do := toProjectDO(p)
	err := adapter.db().Clauses(clause.Returning{}).Create(&do).Error
	if err != nil {
		return spacedomain.Project{}, err
	}

	return *p, nil
}
