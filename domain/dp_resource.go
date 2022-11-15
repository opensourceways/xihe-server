package domain

import (
	"errors"
	"fmt"
	"regexp"
)

const (
	resourceProject = "project"
	resourceDataset = "dataset"
	resourceModel   = "model"

	SortTypeUpdateTime    = "update_time"
	SortTypeFirstLetter   = "first_letter"
	SortTypeDownloadCount = "download_count"
)

var (
	reName         = regexp.MustCompile("^[a-zA-Z0-9_-]+$")
	reResourceName = reName

	resourceLevelMap = map[string]int{
		"official": 2,
		"good":     1,
	}

	ResourceTypeProject = resourceType(resourceProject)
	ResourceTypeModel   = resourceType(resourceModel)
	ResourceTypeDataset = resourceType(resourceDataset)
)

// Name
type ResourceName interface {
	ResourceName() string
	FirstLetterOfName() byte
}

func NewResourceName(v string) (ResourceName, error) {
	max := config.MaxNameLength
	min := config.MinNameLength

	if n := len(v); n > max || n < min {
		return nil, fmt.Errorf("name's length should be between %d to %d", min, max)
	}

	if !reResourceName.MatchString(v) {
		return nil, errors.New("invalid name")
	}

	return resourceName(v), nil
}

type resourceName string

func (r resourceName) ResourceName() string {
	return string(r)
}

func (r resourceName) FirstLetterOfName() byte {
	return string(r)[0]
}

// ResourceType
type ResourceType interface {
	ResourceType() string
	PrefixToName() string
}

func NewResourceType(v string) (ResourceType, error) {
	b := v == resourceProject ||
		v == resourceModel ||
		v == resourceDataset
	if b {
		return resourceType(v), nil
	}

	return nil, errors.New("invalid resource type")
}

type resourceType string

func (r resourceType) ResourceType() string {
	return string(r)
}

func (r resourceType) PrefixToName() string {
	return string(r) + "-"
}

// ResourceDesc
type ResourceDesc interface {
	ResourceDesc() string
}

func IsSameResourceDesc(a, b ResourceDesc) bool {
	if a == nil && b == nil {
		return true
	}

	if a != nil && b != nil {
		return a.ResourceDesc() == b.ResourceDesc()
	}

	return false
}

func NewResourceDesc(v string) (ResourceDesc, error) {
	if v == "" {
		return nil, nil
	}

	if max := config.MaxDescLength; len(v) > max {
		return nil, fmt.Errorf(
			"the length of desc should be less than %d", max,
		)
	}

	return resourceDesc(v), nil
}

type resourceDesc string

func (r resourceDesc) ResourceDesc() string {
	return string(r)
}

// ResourceLevel
type ResourceLevel interface {
	ResourceLevel() string
	Int() int
}

func NewResourceLevel(v string) ResourceLevel {
	for k, n := range resourceLevelMap {
		if k == v {
			return resourceLevel{
				level: n,
				desc:  k,
			}
		}
	}

	return nil

}

func NewResourceLevelByNum(v int) ResourceLevel {
	for k, n := range resourceLevelMap {
		if n == v {
			return resourceLevel{
				level: n,
				desc:  k,
			}
		}
	}

	return nil
}

type resourceLevel struct {
	level int
	desc  string
}

func (r resourceLevel) ResourceLevel() string {
	return r.desc
}

func (r resourceLevel) Int() int {
	return r.level
}

// SortType
type SortType interface {
	SortType() string
}

func NewSortType(v string) (SortType, error) {
	b := v != SortTypeUpdateTime &&
		v != SortTypeFirstLetter &&
		v != SortTypeDownloadCount

	if b {
		return nil, errors.New("invliad sort type")
	}

	return sortType(v), nil
}

type sortType string

func (s sortType) SortType() string {
	return string(s)
}
