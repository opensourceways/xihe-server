/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package primitive

import (
	"errors"
	"path"
	"path/filepath"
	"strings"
)

// FileName
type FileName interface {
	FileName() string
	IsPictureName() bool
	GetFormat() string
}

// NewFileName creates a new Identity instance from a string value.
func NewFileName(v string) (FileName, error) {
	if v == "" {
		return nil, errors.New("enpty filename")
	}

	return fileName(v), nil
}

type fileName string

// FileName returns the string representation of the fileName.
func (f fileName) FileName() string {
	return string(f)
}

func (f fileName) IsPictureName() bool {
	ext := filepath.Ext(f.FileName())

	return allowImageExtension.Has(strings.ToLower(ext))
}

func (f fileName) GetFormat() string {
	return path.Ext(string(f))
}
