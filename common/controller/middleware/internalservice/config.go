/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package internalservice provides functionality for internal services.
package internalservice

var config Config

// Init initializes the internal service with the provided configuration.
func Init(cfg *Config) {
	config = *cfg
}

// Config holds the configuration for the internal service.
type Config struct {
	TokenHash string `json:"token_hash" required:"true"`
	Salt      string `json:"salt" required:"true"`
}
