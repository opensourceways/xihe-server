/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive define the interface of FileRef
package domain

type FileRef interface {
	FileRef() string
}

func NewCodeFileRef(v string) (FileRef, error) {
	return codeFileRef(v), nil
}

func InitCodeFileRef() FileRef {
	return codeFileRef("main")
}

func CreateFileRef(v string) FileRef {
	return codeFileRef(v)
}

type codeFileRef string

func (r codeFileRef) FileRef() string {
	return string(r)
}
