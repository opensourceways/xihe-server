/*
Copyright (c) Huawei Technologies Co., Ltd. 2023. All rights reserved
*/

// Package primitive provides a primitive function in the application.
package primitive

import (
	"errors"
	"strings"
)

const (
	SortByMostLikes       = "most_likes"
	SortByAlphabetical    = "alphabetical"
	SortByMostDownloads   = "most_downloads"
	SortByMostVisits      = "most_visits"
	SortByRecentlyUpdated = "recently_updated"
	SortByRecentlyCreated = "recently_created"
	SortByGlobal          = "global"
	TrueCondition         = "1"
)

// SortType is an interface that defines a method to return the sort type as a string.
type SortType interface {
	SortType() string
}

// NewSortType creates a new SortType based on the input string.
func NewSortType(v string) (SortType, error) {
	switch strings.ToLower(v) {
	case SortByMostLikes:
		return sortType(SortByMostLikes), nil

	case SortByAlphabetical:
		return sortType(SortByAlphabetical), nil

	case SortByMostDownloads:
		return sortType(SortByMostDownloads), nil
	case SortByMostVisits:
		return sortType(SortByMostVisits), nil

	case SortByRecentlyUpdated:
		return sortType(SortByRecentlyUpdated), nil

	case SortByRecentlyCreated:
		return sortType(SortByRecentlyCreated), nil

	case SortByGlobal:
		return sortType(SortByGlobal), nil

	default:
		return nil, errors.New("unknown sort type")
	}
}

type sortType string

// SortType returns the sortType value as a string.
func (s sortType) SortType() string {
	return string(s)
}
