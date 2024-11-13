package infrastructure

import (
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	"github.com/opensourceways/xihe-server/domain"
)

type fileScanRepoImpl struct {
	dao fileScanDAO
}

func NewFileScanRepository() (domain.FileScanRepository, error) {
	do := fileScanDO{}

	if err := pgsql.AutoMigrate(&do); err != nil {
		return nil, err
	}

	return fileScanRepoImpl{dao: pgsql.NewDBTable(do.TableName())}, nil
}
