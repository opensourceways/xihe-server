package domain

import (
	"errors"
	"net/url"

	"github.com/opensourceways/xihe-server/utils"
)

const (
	cloudPodStatusStarting   = "starting"
	cloudPodStatusCreating   = "creating"
	cloudPodStatusFailed     = "failed"
	cloudPodStatusRunning    = "running"
	cloudPodStatusTerminated = "terminated"
)

// CloudName
type CloudName interface {
	CloudName() string
}

func NewCloudName(v string) (CloudName, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return cloudName(v), nil
}

type cloudName string

func (r cloudName) CloudName() string {
	return string(r)
}

// CloudSpec
type CloudSpec interface {
	CloudSpec() string
}

func NewCloudSpec(v string) (CloudSpec, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return cloudSpec(v), nil
}

type cloudSpec string

func (r cloudSpec) CloudSpec() string {
	return string(r)
}

// CloudImage
type CloudImage interface {
	CloudImage() string
}

func NewCloudImage(v string) (CloudImage, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return cloudImage(v), nil
}

type cloudImage string

func (r cloudImage) CloudImage() string {
	return string(r)
}

// CloudFeature
type CloudFeature interface {
	CloudFeature() string
}

func NewCloudFeature(v string) (CloudFeature, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return cloudFeature(v), nil
}

type cloudFeature string

func (r cloudFeature) CloudFeature() string {
	return string(r)
}

// CloudProcessor
type CloudProcessor interface {
	CloudProcessor() string
}

func NewCloudProcessor(v string) (CloudProcessor, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return cloudProcessor(v), nil
}

type cloudProcessor string

func (r cloudProcessor) CloudProcessor() string {
	return string(r)
}

// Credit
type Credit interface {
	Credit() int64
}

func NewCredit(v int64) (Credit, error) {
	if v < 0 {
		return nil, errors.New("invalid value")
	}

	return credit(v), nil
}

type credit int64

func (r credit) Credit() int64 {
	return int64(r)
}

// CloudLimited
type CloudLimited interface {
	CloudLimited() int
}

func NewCloudLimited(v int) (CloudLimited, error) {
	if v < 0 {
		return nil, errors.New("invalid value")
	}

	return cloudLimited(v), nil
}

type cloudLimited int

func (r cloudLimited) CloudLimited() int {
	return int(r)
}

// CloudRemain
type CloudRemain interface {
	CloudRemain() int
}

func NewCloudRemain(v int) (CloudRemain, error) {
	if v < 0 {
		return nil, errors.New("invalid value")
	}

	return cloudRemain(v), nil
}

type cloudRemain int

func (r cloudRemain) CloudRemain() int {
	return int(r)
}

// PodStatus
type PodStatus interface {
	PodStatus() string
	IsStarting() bool
	IsCreating() bool
	IsError() bool
	IsRunning() bool
	IsTerminated() bool
}

func NewPodStatus(v string) (PodStatus, error) {
	if v == "" {
		return nil, errors.New("empty value")
	}

	return podStatus(v), nil
}

type podStatus string

func (r podStatus) PodStatus() string {
	return string(r)
}

func (r podStatus) IsStarting() bool {
	return r.PodStatus() == cloudPodStatusStarting
}

func (r podStatus) IsCreating() bool {
	return r.PodStatus() == cloudPodStatusCreating
}

func (r podStatus) IsError() bool {
	return r.PodStatus() == cloudPodStatusFailed
}

func (r podStatus) IsRunning() bool {
	return r.PodStatus() == cloudPodStatusRunning
}

func (r podStatus) IsTerminated() bool {
	return r.PodStatus() == cloudPodStatusTerminated
}

// PodExpiry
type PodExpiry interface {
	PodExpiry() int64
	PodExpiryDate() string
}

func NewPodExpiry(v int64) (PodExpiry, error) {
	return podExpiry(v), nil
}

type podExpiry int64

func (r podExpiry) PodExpiry() int64 {
	return int64(r)
}

func (r podExpiry) PodExpiryDate() string {
	return utils.ToDate(r.PodExpiry())
}

// PodError
type PodError interface {
	PodError() string
	IsGood() bool
}

func NewPodError(v string) (PodError, error) {
	return podError(v), nil
}

type podError string

func (r podError) PodError() string {
	return string(r)
}

func (p podError) IsGood() bool {
	return p.PodError() == ""
}

// AccessURL
type AccessURL interface {
	AccessURL() string
}

func NewAccessURL(v string) (AccessURL, error) {
	if _, err := url.Parse(v); err != nil {
		return nil, errors.New("invalid url")
	}

	return accessURL(v), nil
}

type accessURL string

func (r accessURL) AccessURL() string {
	return string(r)
}
