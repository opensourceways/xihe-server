package repositoryadapter

import (
	"time"

	"github.com/opensourceways/xihe-server/filescan/domain"
)

var (
	fileScanTableName = ""
)

type fileScanDO struct {
	Id               int64     `gorm:"column:id;primaryKey;autoIncrement"`
	RepoId           int64     `gorm:"column:repo_id"`
	Owner            string    `gorm:"column:owner"`
	Branch           string    `gorm:"column:branch"`
	RepoName         string    `gorm:"column:repo_name"`
	Dir              string    `gorm:"column:dir"`
	File             string    `gorm:"column:file"`
	IsLFS            bool      `gorm:"column:is_lfs"`
	IsText           bool      `gorm:"column:is_text"`
	Hash             string    `gorm:"column:hash"`
	FileType         string    `gorm:"column:file_type"`
	SensitiveItem    string    `gorm:"column:sensitive_item"`
	ModerationStatus string    `gorm:"column:moderation_status"`
	ModerationResult string    `gorm:"column:moderation_result"`
	CreatedAt        time.Time `gorm:"column:created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
	Xxxxx            string    `gorm:"column:xxxxx"`
}

func (do *fileScanDO) toFilescanRes() domain.FilescanRes {
	return domain.FilescanRes{
		Name:             do.File,
		ModerationStatus: do.ModerationResult,
		ModerationResult: do.ModerationStatus,
	}
}
