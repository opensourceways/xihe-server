package repository

import (
	"context"

	"github.com/opensourceways/xihe-server/filescan/domain"
)

type FileScanAdapter interface {
	Get(string, string) ([]domain.FilescanRes, error)
	Find(int64) (domain.FileScan, error)
	Save(*domain.FileScan) error
	RemoveList(context.Context, []int64) error
	AddList(context.Context, []domain.FileScan) ([]domain.FileScan, error)
	FindByRepoIdAndFiles(ctx context.Context, queries []domain.FileScan) ([]domain.FileScan, error)
}
