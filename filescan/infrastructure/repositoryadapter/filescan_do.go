package repositoryadapter

import (
	"time"

	d "github.com/opensourceways/xihe-server/domain"

	"github.com/opensourceways/xihe-server/filescan/domain"

	"github.com/opensourceways/xihe-server/filescan/domain/primitive"
)

var (
	fileScanTableName = ""
)

type fileScanDO struct {
	Id               int64     `gorm:"column:id;primaryKey;autoIncrement"`
	RepoId           int64     `gorm:"column:repo_id"`
	Owner            string    `gorm:"column:owner"`
	RepoName         string    `gorm:"column:repo_name"`
	Dir              string    `gorm:"column:dir"`
	File             string    `gorm:"column:file"`
	IsLFS            bool      `gorm:"column:is_lfs"`
	IsText           bool      `gorm:"column:is_text"`
	Hash             string    `gorm:"column:hash"`
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
		ModerationStatus: do.ModerationStatus,
		ModerationResult: do.ModerationResult,
	}
}

func (do *fileScanDO) toFileScan() domain.FileScan {
	fs := domain.NewFileScan()
	fs.Id = do.Id
	fs.RepoId = do.RepoId
	fs.Owner = d.CreateAccount(do.Owner)
	fs.RepoName = do.RepoName
	fs.Dir = do.Dir
	fs.File = do.File
	fs.LFS = do.IsLFS
	fs.Text = do.IsText
	fs.Hash = do.Hash
	fs.ModerationStatus = primitive.CreateModerationStatus(do.ModerationStatus)
	fs.ModerationResult = primitive.CreateModerationResult(do.ModerationResult)
	fs.SensitiveItem = do.SensitiveItem
	return fs
}

func toFileScanDO(f *domain.FileScan) fileScanDO {
	return fileScanDO{
		Id:               f.Id,
		RepoId:           f.RepoId,
		Owner:            f.Owner.Account(),
		RepoName:         f.RepoName,
		Dir:              f.Dir,
		File:             f.File,
		IsLFS:            f.LFS,
		IsText:           f.Text,
		Hash:             f.Hash,
		SensitiveItem:    f.SensitiveItem,
		ModerationStatus: f.ModerationStatus.Status(),
		ModerationResult: f.ModerationResult.Result(),
	}
}
