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
	Owner            string    `gorm:"column:owner"`
	RepoName         string    `gorm:"column:repo_name"`
	Hash             string    `gorm:"column:hash;index"`
	SensitiveItem    string    `gorm:"column:sensitive_item"`
	ModerationStatus string    `gorm:"column:moderation_status"`
	ModerationResult string    `gorm:"column:moderation_result"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
	File             string    `gorm:"column:file"`
}

func (largeFileScanDO) TableName() string {
	return largeFileScanTableName
}

func (do *largeFileScanDO) toFilescanRes() domain.FilescanRes {
	return domain.FilescanRes{
		Name:             do.File,
		ModerationStatus: do.ModerationResult,
		ModerationResult: do.ModerationStatus,
	}
}
