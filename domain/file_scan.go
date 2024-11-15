package domain

type FileScan struct {
	Id               int64
	LFS              bool
	Text             bool
	Dir              string
	File             string
	Hash             string
	Owner            Account
	Branch           FileRef
	RepoId           int64
	status           FileScanStatus
	ModerationStatus FileModerationStatus
	ModerationResult FileModerationResult
	RepoName         ResourceName
	Size             int64
	FileType         FileType
	ScanItem
}

// ScanItem represents the scan item.
type ScanItem struct {
	SensitiveItem SensitiveItemResult
}
