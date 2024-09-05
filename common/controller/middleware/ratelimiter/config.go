/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package ratelimiter config for rate limiter
package ratelimiter

var config Config

// Init initializes the internal service with the provided configuration.
func Init(cfg *Config) {
	config = *cfg
}

// Config holds the configuration for the rate limiter
type Config struct {
	RequestNum int `json:"request_num" required:"true"`
	BurstNum   int `json:"burst_num"   required:"true"`
}
