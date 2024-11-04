package pgsql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/opensourceways/xihe-server/common/domain/repository"
	"gorm.io/gorm"
)

var (
	errRowExists   = errors.New("row exists")
	errRowNotFound = errors.New("row not found")

	errorCodes errorCode
)

type SortByColumn struct {
	Column string
	Ascend bool
}

func (s SortByColumn) order() string {
	v := " ASC"
	if !s.Ascend {
		v = " DESC"
	}
	return s.Column + v
}

type Pagination struct {
	PageNum      int
	CountPerPage int
}

func (p Pagination) pagination() (limit, offset int) {
	limit = p.CountPerPage

	if limit > 0 && p.PageNum > 0 {
		offset = (p.PageNum - 1) * limit
	}

	return
}

type dbTable struct {
	name string
}

func NewDBTable(name string) dbTable {
	return dbTable{name: name}
}

func (t dbTable) DB() *gorm.DB {
	return db
}

func (t dbTable) Create(result interface{}) error {
	return db.Table(t.name).
		Create(result).
		Error
}

func (t dbTable) Updates(filter, result interface{}) error {
	return db.Table(t.name).
		Where(filter).
		Updates(result).
		Error
}

func (t dbTable) GetRecords(
	filter, result interface{}, p Pagination,
	sort []SortByColumn,
) (err error) {
	query := db.Table(t.name).Where(filter)

	var orders []string
	for _, v := range sort {
		orders = append(orders, v.order())
	}

	if len(orders) >= 0 {
		query.Order(strings.Join(orders, ","))
	}

	if limit, offset := p.pagination(); limit > 0 {
		query.Limit(limit).Offset(offset)
	}

	err = query.Find(result).Error

	return
}

func (t dbTable) Count(filter interface{}) (int, error) {
	var total int64
	err := db.Table(t.name).Where(filter).Count(&total).Error

	return int(total), err
}

func (t dbTable) Filter(filter, result interface{}) error {
	return db.Table(t.name).
		Find(result, filter).Error
}

func (t dbTable) First(filter, result interface{}) error {
	return db.Table(t.name).
		First(result, filter).
		Error
}

func (t dbTable) GetOrderOneRecord(filter, order, result interface{}) error {
	err := db.Table(t.name).Where(filter).Order(order).Limit(1).First(result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errRowNotFound
	}

	return nil
}

func (t dbTable) GetRecord(filter, result interface{}) error {
	err := db.Table(t.name).Where(filter).First(result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errRowNotFound)
	}

	return err
}

func (t dbTable) UpdateRecord(filter, update interface{}) (err error) {
	query := db.Table(t.name).Where(filter).Updates(update)
	if err = query.Error; err != nil {
		return
	}

	if query.RowsAffected == 0 {
		err = errRowNotFound
	}

	return
}

func (t dbTable) IsRowNotFound(err error) bool {
	return errors.Is(err, errRowNotFound)
}

func (t dbTable) IsRowExists(err error) bool {
	return errors.Is(err, errRowExists)
}

func (t dbTable) UpdateWithOmittingSpecificFields(filter, values any, columns ...string) error {
	r := db.Table(t.name).Where(filter).Select(`*`).Omit(columns...).Updates(values)
	if r.Error != nil {
		return r.Error
	}

	if r.RowsAffected == 0 {
		return repository.NewErrorConcurrentUpdating(
			errors.New("concurrent updating"),
		)
	}

	return nil
}

func (t dbTable) GetByPrimaryKey(row any) error {
	err := db.Table(t.name).First(row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

func (t dbTable) Update(filter, values any) error {
	r := db.Table(t.name).Where(filter).Select(`*`).Updates(values)
	if r.Error != nil {
		return r.Error
	}

	if r.RowsAffected == 0 {
		return repository.NewErrorConcurrentUpdating(
			errors.New("concurrent updating"),
		)
	}

	return nil
}

// EqualQuery generates a query string for an "equal" filter condition.
func (dao dbTable) EqualQuery(field string) string {
	return fmt.Sprintf(`%s = ?`, field)
}

// IsRecordExists checks if the given error indicates that a unique constraint violation occurred.
func (dao dbTable) IsRecordExists(err error) bool {
	var pgError *pgconn.PgError
	ok := errors.As(err, &pgError)
	if !ok {
		return false
	}

	return pgError != nil && pgError.Code == errorCodes.UniqueConstraint
}
