package domain

import "errors"

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

func NewConverId(v string) (CoverId, error) {
	if !config.hasCover(v) {
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
	if !config.hasProtocol(v) {
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
	if !config.hasProjectType(v) {
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
	if !config.hasPlatform(v) {
		return nil, errors.New("unsupport training platform")
	}

	return trainingPlatform(v), nil
}

type trainingPlatform string

func (r trainingPlatform) TrainingPlatform() string {
	return string(r)
}
