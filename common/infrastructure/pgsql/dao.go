package pgsql

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

var (
	errRowExists   = errors.New("row exists")
	errRowNotFound = errors.New("row not found")
)

type Impl interface {
	GetRecord(ctx context.Context, filter, result interface{}) error
	GetByPrimaryKey(ctx context.Context, row interface{}) error
	DeleteByPrimaryKey(ctx context.Context, row interface{}) error
	LikeFilter(field, value string) (query, arg string)
	IntersectionFilter(field string, value []string) (query string, arg pq.StringArray)
	EqualQuery(field string) string
	NotEqualQuery(field string) string
	OrderByDesc(field string) string
	InFilter(field string) string
	NotIN(field string) string
	DB() *gorm.DB
	WithContext(context.Context) *gorm.DB
	TableName() string
}

type CommonModel struct {
	ID        int64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

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
		return errRowNotFound
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

// // DB Each operation must generate a new gorm.DB instance.
// // If using the same gorm.DB instance by different operations, they will share the same error.
// func (dao *daoImpl) DB() *gorm.DB {
// 	return db.Table(dao.table)
// }

// func (dao *daoImpl) WithContext(ctx context.Context) *gorm.DB {
// 	return db.WithContext(ctx).Table(dao.table)
// }

// // GetRecord retrieves a single record that matches the given filter criteria.
// func (dao *daoImpl) GetRecord(ctx context.Context, filter, result interface{}) error {
// 	err := dao.WithContext(ctx).Where(filter).First(result).Error

// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		return repository.NewErrorResourceNotExists(errors.New("not found"))
// 	}

// 	return err
// }

// // GetByPrimaryKey retrieves a single record by its primary key.
// func (dao *daoImpl) GetByPrimaryKey(ctx context.Context, row interface{}) error {
// 	err := dao.WithContext(ctx).First(row).Error

// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		return repository.NewErrorResourceNotExists(errors.New("not found"))
// 	}

// 	return err
// }

// // DeleteByPrimaryKey deletes a record by its primary key.
// func (dao *daoImpl) DeleteByPrimaryKey(ctx context.Context, row interface{}) error {
// 	err := dao.WithContext(ctx).Delete(row).Error

// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		return repository.NewErrorResourceNotExists(errors.New("not found"))
// 	}

// 	return err
// }

// // LikeFilter generates a query string and argument for a "like" filter condition.
// func (dao *daoImpl) LikeFilter(field, value string) (query, arg string) {
// 	query = fmt.Sprintf(`%s ilike ?`, field)

// 	arg = `%` + utils.EscapePgsqlValue(value) + `%`

// 	return
// }

// // IntersectionFilter generates a query string and argument for an "intersection" filter condition.
// func (dao *daoImpl) IntersectionFilter(field string, value []string) (query string, arg pq.StringArray) {
// 	query = fmt.Sprintf(`%s @> ?`, field)

// 	arg = pq.StringArray(value)

// 	return
// }

// // EqualQuery generates a query string for an "equal" filter condition.
// func (dao *daoImpl) EqualQuery(field string) string {
// 	return fmt.Sprintf(`%s = ?`, field)
// }

// // MultiEqualQuery generates a query string for multiple "equal" filter conditions.
// func (dao *daoImpl) MultiEqualQuery(fields ...string) string {
// 	v := make([]string, len(fields))

// 	for i, field := range fields {
// 		v[i] = dao.EqualQuery(field)
// 	}

// 	return strings.Join(v, " AND ")
// }

// // NotEqualQuery generates a query string for a "not equal" filter condition.
// func (dao *daoImpl) NotEqualQuery(field string) string {
// 	return fmt.Sprintf(`%s <> ?`, field)
// }

// // OrderByDesc generates a query string for ordering results in descending order by the specified field.
// func (dao *daoImpl) OrderByDesc(field string) string {
// 	return field + " desc"
// }

// // InFilter generates a query string and argument for an "in" filter condition.
// func (dao *daoImpl) InFilter(field string) string {
// 	return fmt.Sprintf(`%s IN ?`, field)
// }

// // NotIn generates a query string and argument for an "not in" filter condition.
// func (dao *daoImpl) NotIN(field string) string {
// 	return fmt.Sprintf(`%s NOT IN (?)`, field)
// }

// // TableName returns the name of the table associated with this daoImpl instance.
// func (dao *daoImpl) TableName() string {
// 	return dao.table
// }

// // IsRecordExists checks if the given error indicates that a unique constraint violation occurred.
// func (dao *daoImpl) IsRecordExists(err error) bool {
// 	var pgError *pgconn.PgError
// 	ok := errors.As(err, &pgError)
// 	if !ok {
// 		return false
// 	}

// 	return pgError != nil && pgError.Code == errorCodes.UniqueConstraint
// }
