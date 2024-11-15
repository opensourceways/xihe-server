/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package domain

import "errors"

const (
	fileTypeReadme  = "readme"
	fileTypeImage   = "image"
	fileTypeDocment = "document"
	fileTypeVideo   = "video"
	fileTypeAudio   = "audio"
	fileTypeUnknow  = "unknown"
)

// FileTitle creates a new commit message with the given title.
type FileType interface {
	FileType() string
	IsReadme() bool
	IsImage() bool
	IsDocment() bool
	IsVideo() bool
	IsAudio() bool
	IsUnknown() bool
}

func CreateFileType(v string) FileType {
	return fileType(v)
}

func NewFileType(v string) (FileType, error) {
	if v != fileTypeReadme && v != fileTypeImage && v != fileTypeDocment && v != fileTypeVideo && v != fileTypeAudio {
		return nil, errors.New("invalid file type")
	}

	return fileType(v), nil
}

func NewFileTypeReadme() FileType {
	return fileType(fileTypeReadme)
}

func NewFileTypeUnknown() FileType {
	return fileType(fileTypeUnknow)
}

type fileType string

func (r fileType) FileType() string {
	return string(r)
}

func (r fileType) IsReadme() bool {
	return string(r) == fileTypeReadme
}

func (r fileType) IsImage() bool {
	return string(r) == fileTypeImage
}

func (r fileType) IsDocment() bool {
	return string(r) == fileTypeDocment
}

func (r fileType) IsVideo() bool {
	return string(r) == fileTypeVideo
}

func (r fileType) IsAudio() bool {
	return string(r) == fileTypeAudio
}

func (r fileType) IsUnknown() bool {
	return string(r) == fileTypeUnknow
}
