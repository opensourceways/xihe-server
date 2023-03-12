package domain

import (
	"errors"
	"net/url"
	"time"
)

// Time
type Time interface {
	Time() int64
	TimeDate() string
}

func NewTime(v int64) (Time, error) {
	if v < 0 {
		return nil, errors.New("invalid value")
	}

	return ptime(v), nil
}

type ptime int64

func (r ptime) Time() int64 {
	return int64(r)
}

func (r ptime) TimeDate() string {
	return time.Unix(r.Time(), 0).Format("2006-01-02")
}

// URL
type URL interface {
	URL() string
}

func NewURL(v string) (URL, error) {
	if v == "" {
		return nil, errors.New("empty url")
	}

	if _, err := url.Parse(v); err != nil {
		return nil, errors.New("invalid url")
	}

	return dpURL(v), nil
}

type dpURL string

func (r dpURL) URL() string {
	return string(r)
}
