package repository

import (
	"github.com/opensourceways/xihe-server/course/domain"
)

type RecordVersion struct {
	Record  domain.Record
	Version int
}

type Record interface {
	AddPlayRecord(*domain.Record) error
	FindPlayRecord(*domain.Record) (RecordVersion, error)
	UpdatePlayRecord(*domain.Record, int) error
}
