package infrastructure

import (
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	"github.com/opensourceways/xihe-server/domain"
)

type largeFileScanRepoImpl struct {
	dao largeFileScanDAO
}

func NewLargeFileScanRepository() (domain.LargeFileScanRepository, error) {
	do := largeFileScanDO{}

	if err := pgsql.AutoMigrate(&do); err != nil {
		return nil, err
	}

	return largeFileScanRepoImpl{dao: pgsql.NewDBTable(do.TableName())}, nil
}
