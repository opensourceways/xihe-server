package repositoryadapter

import "github.com/opensourceways/xihe-server/filescan/domain"

// "errors"

// "github.com/opensourceways/xihe-server/common/domain/repository"

// primitive "github.com/opensourceways/xihe-server/domain"

type FileScanAdapter struct {
	daoImpl
}

// // FindByRecordIndex finds a record based on the account record index and returns an error if any occurs.
// func (adapter *FileScanAdapter) FindByRecordIndex(index domain.FileScanIndex) (
// 	domain.FileScan, error,
// ) {
// 	do := fileScanDO{
// 		Owner:    index.UserName.Account(),
// 		RepoName: index.RepoName.Account(),
// 	}

// 	// It must new a new DO, otherwise the sql statement will include duplicate conditions.
// 	result := fileScanDO{}
// 	if err := adapter.daoImpl.GetRecord(&do, &result); err != nil {
// 		return domain.FileScan{}, err
// 	}

// 	return result.toFileScan(), nil
// }

func (adapter *FileScanAdapter) Get(owner string, repoName string) ([]domain.FilescanRes, error) {

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

func (adapter *FileScanAdapter) GetLarge(owner string, repoName string) ([]domain.FilescanRes, error) {

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
