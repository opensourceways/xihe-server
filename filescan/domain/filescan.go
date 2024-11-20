package domain

import (
	primitive "github.com/opensourceways/xihe-server/domain"
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
	Owner            primitive.Account
	RepoId           int64
	ModerationStatus string
	ModerationResult string
	RepoName         primitive.ResourceName
	Size             int64
	Name             string
}

// type FileScanIndex struct {
// 	UserName primitive.Account
// 	RepoName primitive.Account
// }

type FilescanRes struct {
	ModerationStatus string
	ModerationResult string
	Name             string
}
