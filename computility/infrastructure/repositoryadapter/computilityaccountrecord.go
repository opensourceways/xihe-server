/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package repositoryadapter

import (
	"errors"
	"fmt"

	"gorm.io/gorm/clause"

	commonrepo "github.com/opensourceways/xihe-server/common/domain/repository"
	"github.com/opensourceways/xihe-server/computility/domain"
	primitive "github.com/opensourceways/xihe-server/domain"
)

type computilityAccountRecordAdapter struct {
	daoImpl
}

// Add adds a new computility account record to the database and returns an error if any occurs.
func (adapter *computilityAccountRecordAdapter) Add(
	d *domain.ComputilityAccountRecord,
) error {
	d.Id = primitive.CreateIdentity(primitive.GetId())

	do := toComputilityAccountRecordDO(d)
	return adapter.db().Clauses(clause.Returning{}).Create(&do).Error

}

// Delete deletes a computility account record in the database and returns an error if any occurs.
func (adapter *computilityAccountRecordAdapter) Delete(id primitive.Identity) error {
	return adapter.DeleteByPrimaryKey(
		&computilityAccountDO{Id: id.Integer()},
	)
}

// ListByAccountIndex lists all records based on the account index and returns an error if any occurs.
func (adapter *computilityAccountRecordAdapter) ListByAccountIndex(index domain.ComputilityAccountIndex) (
	[]domain.ComputilityAccountRecord, int, error,
) {
	var result []computilityAccountRecordDO

	sql := fmt.Sprintf(`%s = ? and %s = ?`, filedUserName, filedComputeType)
	query := adapter.daoImpl.db().Where(sql, index.UserName, index.ComputeType)

	err := query.Find(&result).Error
	if err != nil || len(result) == 0 {
		return nil, 0, err
	}

	r := make([]domain.ComputilityAccountRecord, len(result))
	for i := range result {
		r[i] = result[i].toComputilityAccountRecord()
	}

	return r, len(r), nil
}

// FindByRecordIndex finds a record based on the account record index and returns an error if any occurs.
func (adapter *computilityAccountRecordAdapter) FindByRecordIndex(index domain.ComputilityAccountRecordIndex) (
	domain.ComputilityAccountRecord, error,
) {
	do := computilityAccountRecordDO{
		UserName:    index.UserName.Account(),
		ComputeType: index.ComputeType.ComputilityType(),
		SpaceId:     index.SpaceId.Integer(),
	}

	// It must new a new DO, otherwise the sql statement will include duplicate conditions.
	result := computilityAccountRecordDO{}
	if err := adapter.daoImpl.GetRecord(&do, &result); err != nil {
		return domain.ComputilityAccountRecord{}, err
	}

	return result.toComputilityAccountRecord(), nil
}

// Save saves the record in the repository.
func (adapter *computilityAccountRecordAdapter) Save(record *domain.ComputilityAccountRecord) error {
	do := toComputilityAccountRecordDO(record)
	do.Version += 1

	v := adapter.db().Model(
		&computilityAccountRecordDO{Id: record.Id.Integer()},
	).Where(
		equalQuery(filedVersion), record.Version,
	).Select(`*`).Updates(&do)

	if v.Error != nil {
		return v.Error
	}

	if v.RowsAffected == 0 {
		return commonrepo.NewErrorConcurrentUpdating(
			errors.New("concurrent updating"),
		)
	}

	return nil
}
