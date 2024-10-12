/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package repositoryadapter

import (
	"gorm.io/gorm/clause"

	primitive "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/computility/domain"
)

type computilityDetailAdapter struct {
	daoImpl
}

// Add adds a new computility detail record to the database and returns an error if any occurs.
func (adapter *computilityDetailAdapter) Add(d *domain.ComputilityDetail) error {
	d.Id = primitive.CreateIdentity(primitive.GetId())

	do := toComputilityDetailDO(d)
	return adapter.db().Clauses(clause.Returning{}).Create(&do).Error
}

// Delete deletes a computility detail record in the database and returns an error if any occurs.
func (adapter *computilityDetailAdapter) Delete(id primitive.Identity) error {
	return adapter.DeleteByPrimaryKey(
		&computilityDetailDO{Id: id.Integer()},
	)
}

// FindByIndex finds a computility detail record by index and returns an error if any occurs.
func (adapter *computilityDetailAdapter) FindByIndex(d *domain.ComputilityIndex) (domain.ComputilityDetail, error) {
	do := computilityDetailDO{UserName: d.UserName.Account(), OrgName: d.OrgName.Account()}

	// It must new a new DO, otherwise the sql statement will include duplicate conditions.
	result := computilityDetailDO{}

	if err := adapter.daoImpl.GetRecord(&do, &result); err != nil {
		return domain.ComputilityDetail{}, err
	}

	return result.toComputilityDetail(), nil
}

// GetMembers gets all computility detail records related to org by org_name and returns an error if any occurs.
func (adapter *computilityDetailAdapter) GetMembers(orgName primitive.Account) (
	[]domain.ComputilityDetail, error,
) {
	var result []computilityDetailDO

	query := adapter.daoImpl.db().Where(equalQuery(filedOrgName), orgName)

	err := query.Find(&result).Error
	if err != nil || len(result) == 0 {
		return nil, err
	}

	r := make([]domain.ComputilityDetail, len(result))
	for i := range result {
		r[i] = result[i].toComputilityDetail()
	}

	return r, nil
}
