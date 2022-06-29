package domain

import "errors"

const (
	RepoTypePublic  = "public"
	RepoTypePrivate = "priviate"
)

// RepoType
type RepoType interface {
	RepoType() string
}

func NewRepoType(v string) (RepoType, error) {
	if v != RepoTypePublic && v != RepoTypePrivate {
		return nil, errors.New("unknown repo type")
	}

	return repoType(v), nil
}

type repoType string

func (r repoType) RepoType() string {
	return string(r)
}

// ProjName
type ProjName interface {
	ProjName() string
}

func NewProjName(v string) (ProjName, error) {
	// TODO: limited length for name

	return projName(v), nil
}

type projName string

func (r projName) ProjName() string {
	return string(r)
}

// ProjDesc
type ProjDesc interface {
	ProjDesc() string
}

func NewProjDesc(v string) (ProjDesc, error) {
	// TODO: limited length for name

	return projDesc(v), nil
}

type projDesc string

func (r projDesc) ProjDesc() string {
	return string(r)
}

// ProtocolName
type ProtocolName interface {
	ProtocolName() string
}

func NewProtocolName(v string) (ProtocolName, error) {
	// TODO: limited value

	return protocolName(v), nil
}

type protocolName string

func (r protocolName) ProtocolName() string {
	return string(r)
}

// TrainingSDK
type TrainingSDK interface {
	TrainingSDK() string
}

func NewTrainingSDK(v string) (TrainingSDK, error) {
	// TODO: limited value

	return trainingSDK(v), nil
}

type trainingSDK string

func (r trainingSDK) TrainingSDK() string {
	return string(r)
}

// InferenceSDK
type InferenceSDK interface {
	InferenceSDK() string
}

func NewInferenceSDK(v string) (InferenceSDK, error) {
	// TODO: limited value

	return inferenceSDK(v), nil
}

type inferenceSDK string

func (r inferenceSDK) InferenceSDK() string {
	return string(r)
}
