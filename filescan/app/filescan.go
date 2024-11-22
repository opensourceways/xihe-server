// Package app the app of fileScan
package app

import (
	"context"

	filescan "github.com/opensourceways/xihe-server/filescan/domain"
	repo "github.com/opensourceways/xihe-server/filescan/domain/repository"
)

// FileScanAppService is the app service of file scan.
type FileScanService interface {
	Get(string, string) ([]filescan.FilescanRes, error)
	Update(context.Context, CmdToUpdateFileScan) error
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

func (s *fileScanService) Get(owner string, repoName string) ([]filescan.FilescanRes, error) {
	return s.FileScanAdapter.Get(owner, repoName)
}

// Update updates a file scan.
func (s *fileScanService) Update(ctx context.Context, cmd CmdToUpdateFileScan) error {
	fileInfo, err := s.FileScanAdapter.Find(cmd.Id)

	if err != nil {
		return err
	}

	fileInfo.HandleScanDone(cmd.ModerationResult)

	return s.FileScanAdapter.Save(&fileInfo)
}
