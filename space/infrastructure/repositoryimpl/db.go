package repositoryimpl

import "gorm.io/gorm"

var (
	projectAdapterInstance *projectAdapter
)

// Init initializes the database and sets up the necessary adapters.
func Init(db *gorm.DB, tables *Tables) error {
	// must set TableName before migrating

	projectTableName = tables.Project
	tagsTableName = tables.Tags
	datasetTableName = tables.Dataset
	modelTableName = tables.Model

	if err := db.AutoMigrate(&projectDO{}); err != nil {
		return err
	}

	dbInstance = db

	projectDao := daoImpl{
		table: projectTableName, tableTag: tagsTableName, tableDataset: datasetTableName, tableModel: modelTableName}

	projectAdapterInstance = &projectAdapter{
		daoImpl: projectDao,
	}

	return nil
}

// ComputilityAccountRecordAdapter returns the instance of the computilityAccountRecordAdapter.
func ProjectAdapter() *projectAdapter {
	return projectAdapterInstance
}
