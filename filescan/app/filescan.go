// Package app the app of fileScan
package app

import (
	"context"
	"fmt"
	"path/filepath"

	filescan "github.com/opensourceways/xihe-server/filescan/domain"
	"github.com/opensourceways/xihe-server/filescan/domain/primitive"
	repo "github.com/opensourceways/xihe-server/filescan/domain/repository"
	"github.com/sirupsen/logrus"
)

// FileScanAppService is the app service of file scan.
type FileScanService interface {
	Get(string, string) ([]filescan.FilescanRes, error)
	Update(context.Context, CmdToUpdateFileScan) error
	RemoveList(context.Context, RemoveFileScanListCmd) error
	CreateList(context.Context, CreateFileScanListCmd) error
	LaunchModeration(context.Context, LauchModerationCmd) error
}

func NewFileScanService(
	f repo.FileScanAdapter,
	p filescan.ModerationEventPublisher,
) *fileScanService {
	return &fileScanService{
		FileScanAdapter:     f,
		moderationPublisher: p,
	}
}

// fileScanAppService is the app service of file scan.
type fileScanService struct {
	FileScanAdapter     repo.FileScanAdapter
	moderationPublisher filescan.ModerationEventPublisher
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

func (s *fileScanService) RemoveList(ctx context.Context, cmd RemoveFileScanListCmd) error {
	queries := make([]filescan.FileScan, 0, len(cmd.Removed))

	for _, path := range cmd.Removed {
		queries = append(queries, filescan.FileScan{
			RepoId: cmd.RepoId,
			Dir:    filepath.Dir(path),
			File:   filepath.Base(path),
		})
	}

	removedFileScanList, err := s.FileScanAdapter.FindByRepoIdAndFiles(ctx, queries)
	if err != nil {
		return err
	}

	removedIds := make([]int64, 0, len(removedFileScanList))
	for _, v := range removedFileScanList {
		removedIds = append(removedIds, v.Id)
	}

	return s.FileScanAdapter.RemoveList(ctx, removedIds)
}

func (s *fileScanService) CreateList(ctx context.Context, cmd CreateFileScanListCmd) error {
	fileScanList := make([]filescan.FileScan, 0, len(cmd.Added))

	for _, path := range cmd.Added {
		fileScanList = append(fileScanList, filescan.FileScan{
			RepoId:           cmd.RepoId,
			Owner:            cmd.Owner,
			RepoName:         cmd.RepoName,
			Dir:              filepath.Dir(path),
			File:             filepath.Base(path),
			ModerationStatus: primitive.NewInitModerationStatus(),
			ModerationResult: primitive.NewInitModerationResult(),
		})
	}

	addedFileScanList, err := s.FileScanAdapter.AddList(ctx, fileScanList)
	if err != nil {
		return err
	}

	for _, v := range addedFileScanList {
		if err := s.moderationPublisher.Publish(filescan.ModerationEvent{
			ID:       v.Id,
			Owner:    v.Owner.Account(),
			RepoName: v.RepoName,
			Dir:      v.Dir,
			File:     v.File,
		}); err != nil {
			logrus.WithFields(logrus.Fields{
				"operation":   "create",
				"filescan_id": v.Id,
			}).Warnf("fail to publish moderation event, err: %s", err.Error())
		}
	}

	return nil
}

func (s *fileScanService) LaunchModeration(ctx context.Context, cmd LauchModerationCmd) error {
	queries := make([]filescan.FileScan, 0, len(cmd.Modified))

	for _, path := range cmd.Modified {
		queries = append(queries, filescan.FileScan{
			RepoId: cmd.RepoId,
			Dir:    filepath.Dir(path),
			File:   filepath.Base(path),
		})
	}

	fileScanList, err := s.FileScanAdapter.FindByRepoIdAndFiles(ctx, queries)
	if err != nil {
		return err
	}

	for _, v := range fileScanList {
		if err := s.moderationPublisher.Publish(filescan.ModerationEvent{
			ID:       v.Id,
			Owner:    v.Owner.Account(),
			RepoName: v.RepoName,
			Dir:      v.Dir,
			File:     v.File,
		}); err != nil {
			logrus.WithFields(logrus.Fields{
				"operation":   "modify",
				"filescan_id": v.Id,
			}).Warnf("fail to publish moderation event, err: %s", err.Error())
		}
	}

	return nil
}
