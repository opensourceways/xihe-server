/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package domain

import "errors"

const (
	statusInit    = "init"
	statusSkip    = "skip"
	statusScanned = "scanned"
)

// FileScanStatus is an interface that defines the status of a file scan.
type FileScanStatus interface {
	Status() string
	IsDone() bool
	IsSkip() bool
	IsNone() bool
}

// FileScanStatusValidate returns true if the status is valid.
func FileScanStatusValidate(v string) bool {
	return v == statusInit || v == statusSkip || v == statusScanned || v == ""
}

// NewFileScanStatus returns a new FileScanStatus.
func NewFileScanStatus(v string) (FileScanStatus, error) {
	if !FileScanStatusValidate(v) {
		return nil, errors.New("invalid file scan status")
	}

	return fileScanStatus(v), nil
}

// NewInitStatus returns a new FileScanStatus.
func NewInitStatus() FileScanStatus {
	return fileScanStatus(statusInit)
}

// NewScannedStatus returns a new FileScanStatus.
func NewScannedStatus() FileScanStatus {
	return fileScanStatus(statusScanned)
}

// CreateStatus returns a new FileScanStatus.
func CreateStatus(v string) FileScanStatus {
	return fileScanStatus(v)
}

// fileScanStatus is the status of a file scan.
type fileScanStatus string

// Status returns the status of the file scan.
func (s fileScanStatus) Status() string {
	return string(s)
}

// IsDone returns true if the file scan is done.
func (s fileScanStatus) IsDone() bool {
	return s.Status() == statusScanned
}

// IsSkip returns true if the file scan is skipped.
func (s fileScanStatus) IsSkip() bool {
	return s.Status() == statusSkip
}

// IsNone returns true if the file scan status is none.
func (s fileScanStatus) IsNone() bool {
	return s.Status() == ""
}
