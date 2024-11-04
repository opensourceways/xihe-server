/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package sseadapter provides an adapter implementation for working with the repository of space applications.
package sseadapter

import (
	"bufio"
	"io"
	"net/http"

	"github.com/opensourceways/xihe-server/spaceapp/domain"
	"github.com/opensourceways/xihe-server/utils"
)

// const for http
const (
	httpMaxRetries = 3
	httpTimeout    = 3600
)

type Config struct {
	SSEToken string `json:"sse_token"`
}

// StreamSentAdapter creates and returns a new instance of the streamSentAdapter
func StreamSentAdapter(cfg *Config) *streamSentAdapter {
	return &streamSentAdapter{
		utils.NewHttpClient(httpMaxRetries, httpTimeout), cfg,
	}
}

// streamSentAdapter is an adapter for sending server-sent stream requests.
type streamSentAdapter struct {
	cli utils.HttpClient
	cfg *Config
}

// Request sends a server-sent stream request based on the provided SeverSentStream object.
func (sse *streamSentAdapter) Request(q *domain.SeverSentStream) error {
	accessToken := sse.cfg.SSEToken

	req, err := http.NewRequestWithContext(q.Ctx, http.MethodGet, q.Parameter.StreamUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Add("TOKEN", accessToken)

	return sse.cli.SendAndHandle(req, func(h http.Header, respBody io.Reader) error {
		st := streamTransfer{
			input: *bufio.NewReader(respBody),
		}

		q.StreamWrite(st.readAndWriteOnce)

		return nil
	})
}
