// Package app the app of fileScan
package app

import (
	// "github.com/opensourceways/xihe-server/domain"
	filescan "github.com/opensourceways/xihe-server/filescan/domain"
	repo "github.com/opensourceways/xihe-server/filescan/infrastructure/repositoryadapter"
)

// FileScanAppService is the app service of file scan.
type FileScanService interface {
	Get(bool, string, string) ([]filescan.FilescanRes, error)
}

func NewFileScanService(
	f repo.FileScanAdapter,
) *fileScanService {
	return &fileScanService{
		FileScanAdapter: f,
	}
}

// fileScanAppService is the app service of file scan.
type fileScanService struct {
	FileScanAdapter repo.FileScanAdapter
}

func (s *fileScanService) Get(isLFSFile bool, owner string, repoName string) ([]filescan.FilescanRes, error) {
	if isLFSFile {
		return s.FileScanAdapter.Get(owner, repoName)
	}
	return s.FileScanAdapter.GetLarge(owner, repoName)
}
