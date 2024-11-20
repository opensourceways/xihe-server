package repositoryadapter

import (
	"time"

	"github.com/opensourceways/xihe-server/filescan/domain"
)

var (
	largeFileScanTableName = ""
)

type largeFileScanDO struct {
	Id               int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Hash             string    `gorm:"column:hash;index"`
	SensitiveItem    string    `gorm:"column:sensitive_item"`
	ModerationStatus string    `gorm:"column:moderation_status"`
	ModerationResult string    `gorm:"column:moderation_result"`
	CreatedAt        time.Time `gorm:"column:created_at;<-:create"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (do *largeFileScanDO) toFileScan() domain.FileScan {
	return domain.FileScan{
		ModerationStatus: do.ModerationResult,
		ModerationResult: do.ModerationStatus,
	}
}

func (do *largeFileScanDO) toFilescanRes() domain.FilescanRes {
	return domain.FilescanRes{
		ModerationStatus: do.ModerationResult,
		ModerationResult: do.ModerationStatus,
	}
}
