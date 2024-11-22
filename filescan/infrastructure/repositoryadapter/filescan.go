package repositoryadapter

import (
	"fmt"

	"github.com/opensourceways/xihe-server/filescan/domain"
)

const (
	fieldId = "id"
)

type fileScanAdapter struct {
	daoImpl
}

func (adapter *fileScanAdapter) Get(owner string, repoName string) ([]domain.FilescanRes, error) {

	do := fileScanDO{
		Owner:    owner,
		RepoName: repoName,
	}
	var results []fileScanDO

	if err := adapter.daoImpl.GetRecord(&do, &results); err != nil {
		return nil, err
	}

	var fileScans []domain.FilescanRes
	for _, result := range results {
		fileScan := result.toFilescanRes()
		fileScans = append(fileScans, fileScan)
	}

	return fileScans, nil
}

func (adapter *fileScanAdapter) GetLarge(owner string, repoName string) ([]domain.FilescanRes, error) {

	do := fileScanDO{
		Owner:    owner,
		RepoName: repoName,
	}
	var results []largeFileScanDO

	if err := adapter.daoImpl.GetRecordLarge(&do, &results); err != nil {
		return nil, err
	}

	var fileScans []domain.FilescanRes
	for _, result := range results {
		fileScan := result.toFilescanRes()
		fileScans = append(fileScans, fileScan)
	}

	return fileScans, nil
}

func (adapter *fileScanAdapter) Find(id int64) (domain.FileScan, error) {

	do := fileScanDO{
		Id: id,
	}
	var result fileScanDO

	if err := adapter.daoImpl.GetRecord(&do, &result); err != nil {
		return domain.FileScan{}, err
	}

	return result.toFileScan(), nil
}

func (adapter *fileScanAdapter) Save(file *domain.FileScan) error {
	do := toFileScanDO(file)
	fmt.Printf("===================do: %+v\n", do)
	return adapter.db().Where(adapter.daoImpl.EqualQuery(fieldId), do.Id).Save(&do).Error
}
