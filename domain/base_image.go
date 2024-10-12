/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package domain

import (
	"fmt"
	"strings"

	"golang.org/x/xerrors"
)

const (
	PyTorch   = "pytorch"
	MindSpore = "mindspore"
)

// Hardware is an interface that defines hardware-related operations.
type BaseImage interface {
	BaseImage() string
	IsPytorch() bool
	IsMindspore() bool
	Type() string
}

// NewHardware creates a new Hardware instance decided by sdk based on the given string.
func NewBaseImage(v string, hardware string) (BaseImage, error) {
	v = strings.ToLower(strings.TrimSpace(v))
	hardware = strings.ToLower(strings.TrimSpace(hardware))
	fmt.Printf("hardware: %+v\n", hardware)
	fmt.Printf("baseImages: %+v\n", baseImages)
	fmt.Printf("baseImages[hardware]: %+v\n", baseImages[hardware])

	if _, ok := baseImages[hardware]; hardware == "" || !ok {
		fmt.Printf("ok: %+v\n", ok)
		return nil, xerrors.Errorf("unsupported hardware: %s", hardware)
	}

	fmt.Printf("baseImages[hardware].Has(v): %v\n", baseImages[hardware].Has(v))
	if v == "" || !baseImages[hardware].Has(v) {
		return nil, xerrors.Errorf("%s unsupported base image: %s", hardware, v)
	}

	return baseImage(v), nil
}

func IsValidFramework(v string) bool {
	return strings.ToLower(v) == PyTorch || strings.ToLower(v) == MindSpore
}

// CreateHardware creates a new Hardware instance based on the given string.
func CreateBaseImage(v string) BaseImage {
	return baseImage(v)
}

type baseImage string

// BaseImage returns the base image of the base image.
func (r baseImage) BaseImage() string {
	return string(r)
}

func (r baseImage) Type() string {
	if r.IsPytorch() {
		return PyTorch
	}

	return MindSpore
}

// IsPytorch returns whether the base image is PyTorch.
func (r baseImage) IsPytorch() bool {
	return strings.Contains(strings.ToLower(string(r)), PyTorch)
}

// IsMindspore returns whether the base image is MindSpore.
func (r baseImage) IsMindspore() bool {
	return strings.Contains(strings.ToLower(string(r)), PyTorch)
}
