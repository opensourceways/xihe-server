package repositoryadapter

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/opensourceways/xihe-server/common/domain/repository"
)

var dbInstance *gorm.DB

type daoImpl struct {
	table      string
	tableLarge string
}

// Each operation must generate a new gorm.DB instance.
// If using the same gorm.DB instance by different operations, they will share the same error.
func (dao *daoImpl) db() *gorm.DB {
	if dbInstance == nil {
		return nil
	}

	return dbInstance.Table(dao.table)
}

func (dao *daoImpl) dbLarge() *gorm.DB {
	if dbInstance == nil {
		return nil
	}

	return dbInstance.Table(dao.tableLarge)
}

func (dao *daoImpl) GetRecord(filter, result interface{}) error {
	err := dao.db().Where(filter).Find(result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

func (dao *daoImpl) GetRecordLarge(filter, result interface{}) error {
	err := dao.dbLarge().Where(filter).First(result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

// EqualQuery generates a query string for an "equal" filter condition.
func (dao *daoImpl) EqualQuery(field string) string {
	return fmt.Sprintf(`%s = ?`, field)
}

func (impl *daoImpl) Delete(ctx context.Context, filter any) error {
	return impl.db().Delete(filter).Error
}

func (impl *daoImpl) Create(ctx context.Context, value any) error {
	return impl.db().Create(value).Error
}

// GetRecordsOnDisjunction returns results by multiple conditions on disjunction
func (impl *daoImpl) GetRecordsOnDisjunction(ctx context.Context, filter []map[string]any, results any) error {
	db := impl.db().Where("1 = 0")
	for _, cond := range filter {
		db = db.Or(cond)
	}

	err := db.Debug().Find(results).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}
