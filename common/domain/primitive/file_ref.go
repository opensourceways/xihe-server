/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

// FileRef is an interface that represents a file reference.
type FileRef interface {
	FileRef() string
}

// InitCodeFileRef initializes a new CodeFileRef instance with the default value "main".
func InitCodeFileRef() FileRef {
	return codeFileRef("main")
}

type codeFileRef string

// FileRef returns the file reference as a string.
func (r codeFileRef) FileRef() string {
	return string(r)
}
