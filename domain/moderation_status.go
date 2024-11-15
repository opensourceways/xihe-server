/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package domain

import "errors"

const (
	moderationInit        = "init"
	moderationSkip        = "skip"
	moderationScanned     = "scanned"
	moderationToScan      = "to_scan"
	moderationUnsupported = "unsupported_format"
)

// FileModerationStatus is an interface that defines the status of a file scan.
type FileModerationStatus interface {
	Status() string
	IsDone() bool
	IsSkip() bool
	IsNone() bool
	IsBlock() bool
	IsUnsupport() bool
}

// FileModerationStatusValidate returns true if the status is valid.
func FileModerationStatusValidate(v string) bool {
	return v == "" || v == moderationInit || v == moderationSkip ||
		v == moderationScanned || v == moderationToScan || v == moderationUnsupported
}

// NewFileModerationStatus returns a new FileModerationStatus.
func NewFileModerationStatus(v string) (FileModerationStatus, error) {
	if !FileModerationStatusValidate(v) {
		return nil, errors.New("invalid file scan status")
	}

	return fileModerationStatus(v), nil
}

// NewInitModerationStatus returns a new FileModerationStatus.
func NewInitModerationStatus() FileModerationStatus {
	return fileModerationStatus(moderationInit)
}

// NewScannedModerationStatus returns a new FileModerationStatus.
func NewScannedModerationStatus() FileModerationStatus {
	return fileModerationStatus(moderationScanned)
}

// NewUnsupportedModerationStatus returns a new FileModerationStatus.
func NewUnsupportedModerationStatus() FileModerationStatus {
	return fileModerationStatus(moderationUnsupported)
}

// CreateModerationStatus returns a new FileModerationStatus.
func CreateModerationStatus(v string) FileModerationStatus {
	return fileModerationStatus(v)
}

// fileModerationStatus is the status of a file moderation.
type fileModerationStatus string

// Status returns the status of the file scan.
func (s fileModerationStatus) Status() string {
	return string(s)
}

// IsDone returns true if the file scan is done.
func (s fileModerationStatus) IsDone() bool {
	return s.Status() == moderationScanned
}

// IsSkip returns true if the file scan is skipped.
func (s fileModerationStatus) IsSkip() bool {
	return s.Status() == moderationSkip
}

// IsNone returns true if the file moderation status is none.
func (s fileModerationStatus) IsNone() bool {
	return s.Status() == ""
}

func (s fileModerationStatus) IsBlock() bool {
	return s.Status() == "block"
}

func (s fileModerationStatus) IsUnsupport() bool {
	return s.Status() == moderationUnsupported
}
