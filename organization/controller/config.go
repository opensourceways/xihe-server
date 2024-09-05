/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package controller provides functionality for managing the application's controllers.
package controller

// nolint:golint,unused
var config Config

// Init initializes the configuration.
func Init(cfg *Config) {
	config = *cfg
}

// Config represents the application configuration.
type Config struct {
	MaxAvatarFileSize int64 `json:"max_avatar_file_size"`
}

// SetDefault sets the default values for the configuration.
func (cfg *Config) SetDefault() {
	if cfg.MaxAvatarFileSize <= 0 {
		cfg.MaxAvatarFileSize = 1048576
	}
}
