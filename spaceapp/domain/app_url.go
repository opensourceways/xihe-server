/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

package domain

// AppURL is an interface for app url operations.
type AppURL interface {
	AppURL() string
}

// NewAppURL creates a new app URL instance with the given value.
func NewAppURL(v string) (AppURL, error) {
	return appURL(v), nil
}

// CreateAppURL creates a new app URL instance with the given value without validation.
func CreateAppURL(v string) AppURL {
	return appURL(v)
}

type appURL string

// URL returns the URL as a string.
func (r appURL) AppURL() string {
	return string(r)
}
