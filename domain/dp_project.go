package domain

import (
	"errors"
	"strings"
)

const (
	RepoTypePublic  = "public"
	RepoTypePrivate = "private"
	RepoTypeOnline  = "online"
)

// RepoType
type RepoType interface {
	RepoType() string
}

func NewRepoType(v string) (RepoType, error) {
	if v != RepoTypePublic && v != RepoTypePrivate && v != RepoTypeOnline {
		return nil, errors.New("unknown repo type")
	}

	return repoType(v), nil
}

type repoType string

func (r repoType) RepoType() string {
	return string(r)
}

// TrainingPlatform
type CoverId interface {
	CoverId() string
}

func NewCoverId(v string) (CoverId, error) {
	if !DomainConfig.hasCover(v) {
		return nil, errors.New("invalid cover id")
	}

	return coverId(v), nil
}

type coverId string

func (c coverId) CoverId() string {
	return string(c)
}

// ProtocolName
type ProtocolName interface {
	ProtocolName() string
}

func NewProtocolName(v string) (ProtocolName, error) {
	if !DomainConfig.hasProtocol(v) {
		return nil, errors.New("unsupported protocol")
	}

	return protocolName(v), nil
}

type protocolName string

func (r protocolName) ProtocolName() string {
	return string(r)
}

// ProjType
type ProjType interface {
	ProjType() string
}

func NewProjType(v string) (ProjType, error) {
	if !DomainConfig.hasProjectType(v) {
		return nil, errors.New("unsupported project type")
	}

	return projType(v), nil
}

type projType string

func (r projType) ProjType() string {
	return string(r)
}

// TrainingPlatform
type TrainingPlatform interface {
	TrainingPlatform() string
}

func NewTrainingPlatform(v string) (TrainingPlatform, error) {
	if !DomainConfig.hasPlatform(v) {
		return nil, errors.New("unsupported training platform")
	}

	return trainingPlatform(v), nil
}

type trainingPlatform string

func (r trainingPlatform) TrainingPlatform() string {
	return string(r)
}

// Hardware is an interface that defines hardware-related operations.
type Hardware interface {
	Hardware() string
	IsNpu() bool
	IsCpu() bool
}

// NewHardware creates a new Hardware instance decided by sdk based on the given string.
func NewHardware(v string, sdk string) (Hardware, error) {
	v = strings.ToLower(strings.TrimSpace(v))
	sdk = strings.ToLower(strings.TrimSpace(sdk))

	if _, ok := sdkObjects[sdk]; sdk == "" || !ok {
		return nil, errors.New("unsupported sdk")
	}

	if v == "" || !sdkObjects[sdk].Has(v) {
		return nil, errors.New("unsupported hardware")
	}

	return hardware(v), nil
}

func IsValidHardware(h string) bool {
	for _, sdk := range sdkObjects {
		if sdk.Has(h) {
			return true
		}
	}

	return false
}

// CreateHardware creates a new Hardware instance based on the given string.
func CreateHardware(v string) Hardware {
	return hardware(v)
}

type hardware string

// Hardware returns the string representation of the hardware.
func (r hardware) Hardware() string {
	return string(r)
}

func (r hardware) IsNpu() bool {
	return strings.Contains(strings.ToLower(string(r)), "npu")
}

func (r hardware) IsCpu() bool {
	return strings.Contains(strings.ToLower(string(r)), "cpu")
}
