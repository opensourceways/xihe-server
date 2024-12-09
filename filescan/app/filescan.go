// Package app the app of fileScan
package app

import (
	"context"
	"fmt"
	"path/filepath"

	filescan "github.com/opensourceways/xihe-server/filescan/domain"
	repo "github.com/opensourceways/xihe-server/filescan/domain/repository"
)

// FileScanAppService is the app service of file scan.
type FileScanService interface {
	Get(string, string) ([]filescan.FilescanRes, error)
	Update(context.Context, CmdToUpdateFileScan) error
	Remove(context.Context, RemoveFileScanCmd) error
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
	fmt.Printf("===================fileInfo1: %+v\n", fileInfo)

	if err != nil {
		return err
	}

	fileInfo.HandleScanDone(cmd.ModerationResult)
	fmt.Printf("===================fileInfo2: %+v\n", fileInfo)

	return s.FileScanAdapter.Save(&fileInfo)
}

func (s *fileScanService) Remove(ctx context.Context, cmd RemoveFileScanCmd) error {
	files := make([]filescan.FileScan, 0, len(cmd.Removed))

	for _, path := range cmd.Removed {
		files = append(files, filescan.FileScan{
			RepoId: cmd.RepoID,
			Dir:    filepath.Dir(path),
			File:   filepath.Base(path),
		})
	}

	return s.FileScanAdapter.Remove(ctx, files)
}
