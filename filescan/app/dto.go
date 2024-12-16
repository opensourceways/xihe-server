package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/filescan/domain/primitive"
)

type FilescanDTO struct {
	ModerationStatus string `json:"moderation_status"`
	ModerationResult string `json:"moderation_result"`
}

// // CmdToUpdateFileScan is the command to update a file scan.
type CmdToUpdateFileScan struct {
	Id     int64
	Status primitive.FileModerationResult
	// domain.ScanItem
	FileType         string
	ModerationResult primitive.FileModerationResult
}

type RemoveFileScanListCmd struct {
	RepoId  int64
	Removed []string
}

type CreateFileScanListCmd struct {
	RepoId   int64
	Owner    domain.Account
	Branch   string
	RepoName string
	Added    []string
}

type LaunchModerationCmd struct {
	RepoId   int64
	Owner    domain.Account
	Branch   string
	RepoName string
	Modified []string
}
