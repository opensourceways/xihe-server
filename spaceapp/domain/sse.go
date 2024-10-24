/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package domain provides domain models and functionality for managing space apps.
package domain

import (
	"context"
)

// SeverSentStream represents a server-sent stream.
type SeverSentStream struct {
	Parameter   StreamParameter
	Ctx         context.Context
	StreamWrite func(doOnce func() ([]byte, error))
}

// StreamParameter is a type alias for StreamParameter.
type StreamParameter struct {
	StreamUrl string `json:"stream_url" required:"true"`
}

// SeverSentEvent represents a server-sent event.
type SeverSentEvent interface {
	Request(*SeverSentStream) error
}
