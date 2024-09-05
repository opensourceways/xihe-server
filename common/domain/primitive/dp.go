/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"context"
	"errors"
	"time"

	"github.com/opensourceways/xihe-server/utils"
)

var (
	chinese Language = dpLanguage("Chinese")
	english Language = dpLanguage("English")
)

// Time is an interface for time operations.
type Time interface {
	Time() int64
	TimeDate() string
}

// NewTime creates a new Time instance with the given value.
func NewTime(v int64) (Time, error) {
	if v < 0 {
		return nil, errors.New("invalid value")
	}

	return ptime(v), nil
}

type ptime int64

// Time returns the time value as an int64.
func (r ptime) Time() int64 {
	return int64(r)
}

// TimeDate returns the time value formatted as "2006-01-02".
func (r ptime) TimeDate() string {
	return time.Unix(r.Time(), 0).Format("2006-01-02")
}

// URL is an interface for URL operations.
type URL interface {
	URL() string
}

// NewURL creates a new URL instance with the given value.
func NewURL(v string) (URL, error) {
	if v == "" {
		return nil, errors.New("empty url")
	}

	return dpURL(v), nil
}

// CreateURL creates a new URL instance with the given value without validation.
func CreateURL(v string) URL {
	return dpURL(v)
}

type dpURL string

// URL returns the URL as a string.
func (r dpURL) URL() string {
	return string(r)
}

// Website is an interface for Website operations.
type Website interface {
	Website() string
}

// NewOrgWebsite creates a new Website instance with the given value.
func NewOrgWebsite(v string) (Website, error) {
	if v == "" {
		return dpWebsite(v), nil
	}

	if len(v) > websiteConfig.MaxLength {
		return nil, errors.New("invalid website length")
	}

	if !websiteConfig.regexp.MatchString(v) {
		return nil, errors.New("invalid Website")
	}

	if !utils.IsUrl(v) {
		return nil, errors.New("invalid website")
	}

	return dpWebsite(v), nil
}

// CreateOrgWebsite creates a new Website instance with the given value without validation.
func CreateOrgWebsite(v string) Website {
	return dpWebsite(v)
}

type dpWebsite string

// Website returns the Website as a string.
func (r dpWebsite) Website() string {
	return string(r)
}

// Language
type Language interface {
	Language() string
}

type dpLanguage string

// Language is an interface for language operations.
func (r dpLanguage) Language() string {
	return string(r)
}

// NewLanguage creates a new Language instance based on the given value.
func NewLanguage(v string) Language {
	switch v {
	case chinese.Language():
		return chinese

	case english.Language():
		return english

	default:
		return nil
	}
}

// SupportedLanguages returns a list of supported languages.
func SupportedLanguages() []Language {
	return []Language{chinese, english}
}

// WithContext executes the given function with a context that has a timeout.
func WithContext(f func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second, // TODO use config
	)
	defer cancel()

	return f(ctx)
}
