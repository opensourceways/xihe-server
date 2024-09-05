/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package domain provides domain models and configuration for a specific functionality.
package domain

import (
	"github.com/opensourceways/xihe-server/organization/domain/primitive"
	"github.com/opensourceways/xihe-server/organization/infrastructure/messageadapter"
)

// Config is a structure that holds the configuration settings for the application.
type Config struct {
	MaxCountPerOwner int64                 `json:"max_count_per_owner"`
	MaxInviteCount   int64                 `json:"max_invite_count"`
	InviteExpiry     int64                 `json:"invite_expiry"`
	DefaultRole      string                `json:"default_role"`
	Tables           tables                `json:"tables"`
	Topics           messageadapter.Topics `json:"topics"`
	Primitive        primitive.Config      `json:"primitive"`
	CertificateEmail []string              `json:"certificate_email"`
	MaxStrSize       int64                 `json:"max_str_size"`
}

type tables struct {
	Member      string `json:"member"      required:"true"`
	Invite      string `json:"invite"      required:"true"`
	Certificate string `json:"certificate" required:"true"`
}

// SetDefault sets the default values for the Config struct if they are not already set.
func (cfg *Config) SetDefault() {
	if cfg.MaxCountPerOwner <= 0 {
		cfg.MaxCountPerOwner = 10
	}

	if cfg.MaxInviteCount <= 0 {
		cfg.MaxInviteCount = 100
	}
	if cfg.MaxStrSize <= 0 {
		cfg.MaxStrSize = 20
	}
}

// ConfigItems returns a slice of interface{} containing pointers to
// the configuration items in the Config struct.
func (cfg *Config) ConfigItems() []interface{} {
	return []interface{}{
		&cfg.Primitive,
	}
}

// Init initializes the application using the configuration settings provided in the Config struct.
func (cfg *Config) Init() error {
	return primitive.Init(cfg.Primitive)
}
