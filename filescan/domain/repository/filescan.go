package repository

import (
	"context"

	"github.com/opensourceways/xihe-server/filescan/domain"
)

type FileScanAdapter interface {
	Get(string, string) ([]domain.FilescanRes, error)
	Find(int64) (domain.FileScan, error)
	Save(*domain.FileScan) error
	Remove(context.Context, []domain.FileScan) error
}
