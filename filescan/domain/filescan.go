package domain

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/filescan/domain/primitive"
)

type LargeFileScan struct {
	Id               int64
	Hash             string
	Name             string
	ModerationStatus string
	ModerationResult string
}

type FileScan struct {
	Id               int64
	LFS              bool
	Text             bool
	Dir              string
	File             string
	Hash             string
	Owner            domain.Account
	RepoId           int64
	ModerationStatus primitive.FileModerationStatus
	ModerationResult primitive.FileModerationResult
	RepoName         string
	Size             int64
	Name             string
	SensitiveItem    string
}

type FileScanIndex struct {
	Owner    domain.Account
	RepoName domain.Account
}

type FilescanRes struct {
	ModerationStatus string
	ModerationResult string
	Name             string
}

func NewFileScan() FileScan {
	return FileScan{}
}

// HandleScanDone handles the scan done.
func (l *FileScan) HandleScanDone(m primitive.FileModerationResult) {
	if !m.IsNone() {
		l.ModerationResult = m
		if m.IsUnsupported() {
			l.ModerationStatus = primitive.NewUnsupportedModerationStatus()
		} else {
			l.ModerationStatus = primitive.NewScannedModerationStatus()
		}
	}
}
