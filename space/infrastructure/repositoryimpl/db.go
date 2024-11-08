package repositoryimpl

import "gorm.io/gorm"

var (
	projectAdapterInstance *projectAdapter
	datasetAdapterInstance *datasetAdapter
	modelAdapterInstance   *modelAdapter
)

// Init initializes the database and sets up the necessary adapters.
func Init(db *gorm.DB, tables *Tables) error {
	// must set TableName before migrating

	projectTableName = tables.Project
	tagsTableName = tables.Tags
	datasetTableName = tables.Dataset
	modelTableName = tables.Model

	// if err := db.AutoMigrate(&projectDO{}); err != nil {
	// 	return err
	// }

	dbInstance = db

	projectDao := daoImpl{table: projectTableName, tableTag: tagsTableName}
	datasetDao := relatedDaoImpl{table: datasetTableName}
	modelDao := relatedDaoImpl{table: modelTableName}

	projectAdapterInstance = &projectAdapter{
		daoImpl: projectDao,
	}
	datasetAdapterInstance = &datasetAdapter{
		relatedDaoImpl: datasetDao,
	}
	modelAdapterInstance = &modelAdapter{
		relatedDaoImpl: modelDao,
	}

	return nil
}

// ComputilityAccountRecordAdapter returns the instance of the computilityAccountRecordAdapter.
func ProjectAdapter() *projectAdapter {
	return projectAdapterInstance
}

func DatasetAdapter() *datasetAdapter {
	return datasetAdapterInstance
}

func ModelAdapter() *modelAdapter {
	return modelAdapterInstance
}
