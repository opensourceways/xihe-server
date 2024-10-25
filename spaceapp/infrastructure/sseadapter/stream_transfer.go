/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package sseadapter provides an adapter implementation for working with the repository of space applications.
package sseadapter

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

const (
	lineSplitSize = 2
)

type streamTransfer struct {
	input bufio.Reader
}

// Event object is a representation of single chunk of data in event stream.
type Event struct {
	ID    string
	Event string
	Data  []byte
}

func (impl *streamTransfer) readAndWriteOnce() ([]byte, error) {
	event, err := impl.parseEvent(&impl.input)
	if err != nil {
		return nil, err
	}
	if event.Event == "finish" {
		return nil, errors.New("finish")
	}
	// ignore empty events
	if len(event.Data) == 0 {
		return nil, nil
	}
	return event.Data, nil
}

// parseEvent reads a single Event fromthe event stream.
func (impl *streamTransfer) parseEvent(r *bufio.Reader) (*Event, error) {
	line, err := r.ReadBytes('\n')
	if err != nil {
		if errors.Is(err, io.EOF) && len(line) != 0 {
			err = errors.New("incomplete event at the end of the stream")
		}
		return nil, err
	}

	line = impl.chomp(line)
	return impl.parseProtocolData(line)
}

// chomp removes \r or \n or \r\n suffix from the given byte slice.
func (impl *streamTransfer) chomp(b []byte) []byte {
	if len(b) > 0 && b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	if len(b) > 0 && b[len(b)-1] == '\r' {
		b = b[:len(b)-1]
	}
	return b
}

func (impl *streamTransfer) parseProtocolData(line []byte) (*Event, error) {
	event := &Event{
		ID:    "",
		Event: "message",
	}
	if len(line) == 0 {
		return event, nil
	}
	parts := bytes.SplitN(line, []byte(":"), lineSplitSize)

	// Make sure parts[1] always exist
	if len(parts) == 1 {
		parts = append(parts, nil)
	}

	// Chomp space after ":"
	if len(parts[1]) > 0 && parts[1][0] == ' ' {
		parts[1] = parts[1][1:]
	}
	switch string(parts[0]) {
	case "id":
		event.ID = string(parts[1])
	case "event":
		event.Event = string(parts[1])
	case "data":
		if event.Data != nil {
			event.Data = append(event.Data, '\n')
		}
		event.Data = append(event.Data, parts[1]...)
	default:
		return nil, errors.New("invalid event")
	}
	return event, nil
}
