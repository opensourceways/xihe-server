package repositoryimpl

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/opensourceways/xihe-server/common/domain/repository"
)

var dbInstance *gorm.DB

type daoImpl struct {
	table        string
	tableTag     string
	tableDataset string
	tableModel   string
}

// Each operation must generate a new gorm.DB instance.
// If using the same gorm.DB instance by different operations, they will share the same error.
func (dao *daoImpl) db() *gorm.DB {
	if dbInstance == nil {
		return nil
	}

	return dbInstance.Table(dao.table)
}

func (dao *daoImpl) dbTag() *gorm.DB {
	if dbInstance == nil {
		return nil
	}

	return dbInstance.Table(dao.tableTag)
}

func (dao *daoImpl) dbDataset() *gorm.DB {
	if dbInstance == nil {
		return nil
	}

	return dbInstance.Table(dao.tableDataset)
}

func (dao *daoImpl) dbModel() *gorm.DB {
	if dbInstance == nil {
		return nil
	}

	return dbInstance.Table(dao.tableModel)
}

// GetRecord retrieves a single record from the database based on the provided filter
// and stores it in the result parameter.
func (dao *daoImpl) GetProjectRecord(filter, result interface{}) error {
	err := dao.db().Where(filter).First(result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

func equalQuery(field string) string {
	return fmt.Sprintf(`%s = ?`, field)
}

// DeleteByPrimaryKey deletes a single record from the database based on the primary key of the row parameter.
func (dao *daoImpl) DeleteByPrimaryKey(row interface{}) error {
	err := dao.db().Delete(row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}
