package repositoryadapter

import (
	"fmt"

	"gorm.io/gorm"
)

var (
	fileScanInstance *FileScanAdapter
)

// Init initializes the database and sets up the necessary adapters.
func Init(db *gorm.DB, tables *Tables) error {
	// must set TableName before migrating
	fileScanTableName = tables.FileScan
	largeFileScanTableName = tables.LargeFileScan

	if err := db.AutoMigrate(&fileScanDO{}); err != nil {
		fmt.Printf("==========================Init err1: %v\n", err)
		return err
	}

	if err := db.AutoMigrate(&largeFileScanDO{}); err != nil {
		fmt.Printf("==========================Init err2: %v\n", err)
		return err
	}

	dbInstance = db

	filescanDao := daoImpl{table: fileScanTableName, tableLarge: largeFileScanTableName}

	fileScanInstance = &FileScanAdapter{
		daoImpl: filescanDao,
	}

	return nil
}

// ComputilityOrgAdapter returns the instance of the computilityOrgAdapter.
func NewFileScanAdapter() *FileScanAdapter {
	return fileScanInstance
}
