package repositoryadapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/opensourceways/xihe-server/filescan/domain"
)

const (
	fieldId = "id"
)

var ErrNoFileScan = errors.New("no file scan")

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

func (adapter *fileScanAdapter) RemoveList(ctx context.Context, ids []int64) error {
	filter := make([]fileScanDO, 0, len(ids))

	for _, id := range ids {
		filter = append(filter, fileScanDO{
			Id: id,
		})
	}

	return adapter.Delete(ctx, &filter)
}

func (adapter *fileScanAdapter) AddList(
	ctx context.Context, fileScanList []domain.FileScan,
) ([]domain.FileScan, error) {
	fileScanListDO := make([]*fileScanDO, 0)

	for i := range fileScanList {
		do := toFileScanDO(&fileScanList[i])
		fileScanListDO = append(fileScanListDO, &do)
	}

	addedFileScanList := make([]domain.FileScan, 0, len(fileScanListDO))
	if err := adapter.Create(ctx, fileScanListDO); err != nil {
		return addedFileScanList, err
	}

	for _, v := range fileScanListDO {
		addedFileScanList = append(addedFileScanList, v.toFileScan())
	}

	return addedFileScanList, nil
}

func (adapter *fileScanAdapter) FindByRepoIdAndFiles(
	ctx context.Context, queries []domain.FileScan,
) ([]domain.FileScan, error) {
	conditions := make([]map[string]any, 0, len(queries))

	for _, v := range queries {
		conditions = append(conditions, map[string]any{
			"repo_id": v.RepoId,
			"dir":     v.Dir,
			"file":    v.File,
		})
	}

	var results []fileScanDO
	if err := adapter.GetRecordsOnDisjunction(ctx, conditions, &results); err != nil {
		return nil, err
	}

	var fileScanList []domain.FileScan
	for _, result := range results {
		fileScanList = append(fileScanList, result.toFileScan())
	}

	return fileScanList, nil
}
