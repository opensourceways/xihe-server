package repositoryadapter

import (
	"gorm.io/gorm"
)

var (
	fileScanInstance *fileScanAdapter
)

// Init initializes the database and sets up the necessary adapters.
func Init(db *gorm.DB, tables *Tables) error {
	// must set TableName before migrating
	fileScanTableName = tables.FileScan
	largeFileScanTableName = tables.LargeFileScan

	if err := db.AutoMigrate(&fileScanDO{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&largeFileScanDO{}); err != nil {
		return err
	}

	dbInstance = db

	filescanDao := daoImpl{table: fileScanTableName, tableLarge: largeFileScanTableName}

	fileScanInstance = &fileScanAdapter{
		daoImpl: filescanDao,
	}

	return nil
}

// ComputilityOrgAdapter returns the instance of the computilityOrgAdapter.
func NewFileScanAdapter() *fileScanAdapter {
	return fileScanInstance
}
