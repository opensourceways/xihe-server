/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

package primitive

import (
	"encoding/base64"
	"errors"

	"github.com/h2non/filetype"
)

const (
	typeJPG  = "jpg"
	typePNG  = "png"
	typeJPEG = "jpeg"

	maxSize = 10 * 1024 * 1024 // 10M
)

// Image is an interface for image.
type Image interface {
	ImageType() string
	Content() []byte
	ContentOfBase64() string
}

// NewImage creates a new Image instance based on the given byte array.
func NewImage(v []byte) (Image, error) {
	i, err := filetype.Image(v)
	if err != nil {
		return nil, errors.New("not a image")
	}

	imageType := i.Extension
	if imageType != typeJPG && imageType != typePNG && imageType != typeJPEG {
		return nil, errors.New("unsupported type")
	}

	if len(v) > maxSize {
		return nil, errors.New("exceed 10M")
	}

	return image{
		imgType: imageType,
		content: v,
	}, nil
}

type image struct {
	imgType string
	content []byte
}

// ImageType returns the image type.
func (i image) ImageType() string {
	return i.imgType
}

// Content returns the image content.
func (i image) Content() []byte {
	return i.content
}

// ContentOfBase64 returns the image content encoded in base64.
func (i image) ContentOfBase64() string {
	return base64.StdEncoding.EncodeToString(i.content)
}
