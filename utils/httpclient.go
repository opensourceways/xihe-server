/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package utils provides utility functions for various purposes.
package utils

import (
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

// const for http
const (
	statusCodeUpLimit   = 200
	statusCodeDownLimit = 299
	defaultBackoff      = 10 * time.Millisecond
)

// HttpClient is a http client
type HttpClient struct {
	client     http.Client
	maxRetries int
}

// NewHttpClient returns a new http client
func newClient(timeout int) http.Client {
	return http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
}

// NewHttpClient returns a new http client
func NewHttpClient(n, timeout int) HttpClient {
	return HttpClient{
		maxRetries: n,
		client:     newClient(timeout),
	}
}

// SendAndHandle sends request and handle response
func (hc *HttpClient) SendAndHandle(req *http.Request, handle func(http.Header, io.Reader) error) error {
	resp, err := hc.do(req)
	if err != nil || resp == nil {
		return xerrors.Errorf("send request error: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < statusCodeUpLimit || resp.StatusCode > statusCodeDownLimit {
		rb, err := io.ReadAll(resp.Body)
		logrus.Infof("get logs failed, body:%s, err:%s", rb, err)
		return xerrors.Errorf("No logs")
	}

	if handle != nil {
		err = handle(resp.Header, resp.Body)
		if err != nil {
			err = xerrors.Errorf("handle response error: %w", err)
		}

		return err
	}

	return nil
}

// do sends request
func (hc *HttpClient) do(req *http.Request) (resp *http.Response, err error) {
	if resp, err = hc.client.Do(req); err == nil {
		return
	}

	maxRetries := hc.maxRetries
	backoff := defaultBackoff

	for retries := 1; retries < maxRetries; retries++ {
		time.Sleep(backoff)
		backoff *= 2

		if resp, err = hc.client.Do(req); err == nil {
			break
		}
	}
	return
}
