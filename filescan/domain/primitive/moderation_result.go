/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package primitive

import "errors"

const (
	moderationResultBlock       = "block"
	moderationResultPass        = "pass"
	moderationResultReview      = "review"
	moderationResultUnsupported = "unsupported_format"
)

// FileModerationResult is an interface that defines the result of a moderation.
type FileModerationResult interface {
	Result() string
	IsNone() bool
	IsBlock() bool
	IsPass() bool
	IsUnsupported() bool
}

// FileModerationResultValidate returns true if the result is valid.
func FileModerationResultValidate(v string) bool {
	return v == "" || v == moderationResultBlock || v == moderationResultPass ||
		v == moderationResultReview || v == moderationResultUnsupported
}

func NewUnsupportedResult() FileModerationResult {
	return fileModerationResult(moderationResultUnsupported)
}

// NewFileModerationResult returns a new FileModerationResult.
func NewFileModerationResult(v string) (FileModerationResult, error) {
	if !FileModerationResultValidate(v) {
		return nil, errors.New("invalid file moderation result")
	}

	return fileModerationResult(v), nil
}

// NewInitModerationResult returns a new InitModerationResult.
func NewInitModerationResult() FileModerationResult {
	return fileModerationResult("")
}

// CreateModerationStatus returns a new FileModerationResult.
func CreateModerationResult(v string) FileModerationResult {
	return fileModerationResult(v)
}

// fileModerationStatus is the status of a file moderation.
type fileModerationResult string

// Result returns the result of the file moderation.
func (s fileModerationResult) Result() string {
	return string(s)
}

// IsNone returns true if the file moderation result is none.
func (s fileModerationResult) IsNone() bool {
	return s.Result() == ""
}

// IsNone returns true if the file moderation result is block.
func (s fileModerationResult) IsBlock() bool {
	return s.Result() == moderationResultBlock
}

func (s fileModerationResult) IsPass() bool {
	return s.Result() == moderationResultPass
}

func (s fileModerationResult) IsUnsupported() bool {
	return s.Result() == moderationResultUnsupported
}
