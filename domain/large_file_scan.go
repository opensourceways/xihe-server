package domain

type LargeFileScan struct {
	Id               int64
	Hash             string
	status           FileScanStatus
	ModerationStatus FileModerationStatus
	ModerationResult FileModerationResult
	ScanItem
}
