package repositoryimpl

import "gorm.io/gorm"

type SpaceAppDAO interface {
	DB() *gorm.DB
	EqualQuery(field string) string
	IsRecordExists(err error) bool
	UpdateWithOmittingSpecificFields(filter, values any, columns ...string) error
	GetRecord(filter, result any) error
	GetByPrimaryKey(row any) error
	Update(filter, values any) error
}
