package domain

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/opensourceways/xihe-server/utils"
)

const (
	ResourceProject = "project"
	ResourceDataset = "dataset"
	ResourceModel   = "model"

	SortTypeUpdateTime    = "update_time"
	SortTypeFirstLetter   = "first_letter"
	SortTypeDownloadCount = "download_count"
)

var (
	ReName         = regexp.MustCompile("^[a-zA-Z0-9_-]+$")
	reResourceName = ReName

	resourceLevelMap = map[string]int{
		"official": 2,
		"good":     1,
	}

	ResourceTypeProject = resourceType(ResourceProject)
	ResourceTypeModel   = resourceType(ResourceModel)
	ResourceTypeDataset = resourceType(ResourceDataset)
)

// DomainValue
type DomainValue interface {
	DomainValue() string
}

func IsSameDomainValue(a, b DomainValue) bool {
	if a == nil && b == nil {
		return true
	}

	if a != nil && b != nil {
		return a.DomainValue() == b.DomainValue()
	}

	return false
}

// Name
type ResourceName interface {
	ResourceName() string
	FirstLetterOfName() byte
}

func NewResourceName(v string) (ResourceName, error) {
	max := DomainConfig.MaxNameLength
	min := DomainConfig.MinNameLength

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

func NewEmptyResourceName() ResourceName {
	return emptyResourceName("")
}

type emptyResourceName string

func (e emptyResourceName) ResourceName() string {
	return ""
}

func (e emptyResourceName) FirstLetterOfName() byte {
	return byte(0)
}

// ResourceTitle
type ResourceTitle interface {
	ResourceTitle() string
	DomainValue
}

func NewResourceTitle(v string) (ResourceTitle, error) {
	if v == "" {
		return resourceTitle(v), nil
	}

	max := DomainConfig.MaxTitleLength
	min := DomainConfig.MinTitleLength
	if n := utils.StrLen(v); n > max || n < min {
		return nil, fmt.Errorf("title's length should be between %d to %d", min, max)
	}

	return resourceTitle(v), nil
}

type resourceTitle string

func (r resourceTitle) ResourceTitle() string {
	return string(r)
}

func (r resourceTitle) DomainValue() string {
	return string(r)
}

// ResourceType
type ResourceType interface {
	ResourceType() string
}

func NewResourceType(v string) (ResourceType, error) {
	b := v == ResourceProject ||
		v == ResourceModel ||
		v == ResourceDataset
	if b {
		return resourceType(v), nil
	}

	return nil, errors.New("invalid resource type")
}

type resourceType string

func (r resourceType) ResourceType() string {
	return string(r)
}

// ResourceDesc
type ResourceDesc interface {
	ResourceDesc() string

	DomainValue
}

func NewResourceDesc(v string) (ResourceDesc, error) {
	if v == "" {
		return resourceDesc(v), nil
	}

	v = utils.XSSFilter(v)

	if max := DomainConfig.MaxDescLength; utils.StrLen(v) > max {
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

func (r resourceDesc) DomainValue() string {
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
		return nil, errors.New("invalid sort type")
	}

	return sortType(v), nil
}

type sortType string

func (s sortType) SortType() string {
	return string(s)
}
