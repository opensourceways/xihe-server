/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/xerrors"
)

// Avatar is an interface for CdnImageURL operations.
type Avatar interface {
	URL() string
	Storage() string
}

// NewAvatar creates a new URL instance with the given value.
func NewAvatar(v string) (Avatar, error) {
	if v == "" {
		return avatar(v), nil
	}

	imageURL, err := url.ParseRequestURI(v)
	if err != nil {
		return nil, xerrors.Errorf("invalid url: %w", err)
	}

	if skipAvatarids.Has(v) {
		return avatar(""), nil
	}

	for _, domain := range acceptableAvatarDomains {
		if strings.HasPrefix(v, domain) {
			return avatar(imageURL.String()), nil
		}
	}

	if strings.HasPrefix(v, cdnUrlConfig) {
		return avatar(imageURL.String()), nil
	}

	return nil, errors.New("invalid image url")
}

// CreateAvatar is usually called internally, such as repository.
func CreateAvatar(v string) Avatar {
	if skipAvatarids.Has(v) {
		return avatar("")
	}

	return avatar(v)
}

type avatar string

// Storage returns the avatar storage format as a string.
func (r avatar) Storage() string {
	if strings.HasPrefix(string(r), cdnUrlConfig) {
		s := strings.TrimPrefix(string(r), cdnUrlConfig)

		return s
	}

	return string(r)
}

// URL returns URL as string.
func (r avatar) URL() string {
	if r == "" {
		return ""
	}

	if skipAvatarids.Has(string(r)) {
		return ""
	}

	for _, domain := range acceptableAvatarDomains {
		if strings.HasPrefix(string(r), domain) {
			return string(r)
		}
	}

	return fmt.Sprintf("%s%s", cdnUrlConfig, r)
}
