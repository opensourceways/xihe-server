package repository

import (
	"github.com/opensourceways/xihe-server/filescan/domain"
)

type FileScanService interface {
	Get(string, string) []domain.FilescanRes
}
