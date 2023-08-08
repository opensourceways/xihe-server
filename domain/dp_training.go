package domain

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/opensourceways/xihe-server/utils"
)

const (
	rootDirectory = ""

	computeTypeAscend = "Ascend-Powered-Engine"
	computeTypeMPI    = "MPI"

	computeFlaverAscend = "modelarts.kat1.xlarge.public"
	computeFlaverGPU    = "modelarts.p3.large.public"

	computeVersionCudaMS13   = "mindspore_1.3.0-cuda_10.1-py_3.7-ubuntu_1804-x86_64"
	computeVersionCannMS17   = "mindspore_1.7.0-cann_5.1.0-py_3.7-euler_2.8.3-aarch64"
	computeVersionCannMS19_1 = "mindspore_1.9.0-cann_6.0.RC1-py_3.7-ubuntu_18.04-amd64"
	computeVersionCannMS19_2 = "mindspore_1.9.0-cann_6.0.RC1-py_3.9-ubuntu_18.04-amd64"
)

var (
	reDirectory = regexp.MustCompile("^[a-zA-Z0-9_/-]+$")
	reFile      = regexp.MustCompile("^[a-zA-Z0-9_.-]+$")
)

// TrainingName
type TrainingName interface {
	TrainingName() string
}

func NewTrainingName(v string) (TrainingName, error) {
	v = utils.XSSFilter(v)

	max := DomainConfig.MaxTrainingNameLength
	min := DomainConfig.MinTrainingNameLength

	if n := utils.StrLen(v); n > max || n < min {
		return nil, fmt.Errorf("name's length should be between %d to %d", min, max)
	}

	if !reName.MatchString(v) {
		return nil, errors.New("invalid name")
	}

	return trainingName(v), nil
}

type trainingName string

func (r trainingName) TrainingName() string {
	return string(r)
}

// TrainingDesc
type TrainingDesc interface {
	TrainingDesc() string
}

func NewTrainingDesc(v string) (TrainingDesc, error) {
	if v == "" {
		return trainingDesc(v), nil
	}

	v = utils.XSSFilter(v)

	max := DomainConfig.MaxTrainingDescLength
	if utils.StrLen(v) > max {
		return nil, fmt.Errorf("the length of desc should be less than %d", max)
	}

	return trainingDesc(v), nil
}

type trainingDesc string

func (r trainingDesc) TrainingDesc() string {
	return string(r)
}

// Directory
type Directory interface {
	Directory() string
	IsRootDir() bool
}

func NewDirectory(v string) (Directory, error) {
	if v == "" {
		return directory(rootDirectory), nil
	}

	v = utils.XSSFilter(v)

	if max := 50; utils.StrLen(v) > max {
		return nil, fmt.Errorf("the length of directory should be less than %d", max)
	}

	if !reDirectory.MatchString(v) {
		return nil, errors.New("invalid directory")
	}

	return directory(v), nil
}

type directory string

func (r directory) Directory() string {
	return string(r)
}

func (r directory) IsRootDir() bool {
	return string(r) == rootDirectory
}

// FilePath
type FilePath interface {
	FilePath() string
}

func NewFilePath(v string) (FilePath, error) {
	if v == "" {
		return nil, errors.New("empty file path")
	}

	v = utils.XSSFilter(v)

	if max := 50; utils.StrLen(v) > max {
		return nil, fmt.Errorf("the length of file path should be less than %d", max)
	}

	dir := filepath.Dir(v)
	if dir == "." {
		dir = ""
	}
	if dir != "" && !reDirectory.MatchString(dir) {
		return nil, errors.New("invalid filePath")
	}

	file := filepath.Base(v)
	if !reFile.MatchString(file) || file == "." || file == ".." {
		return nil, errors.New("invalid filePath")
	}

	return filePath(filepath.Join(dir, file)), nil
}

type filePath string

func (r filePath) FilePath() string {
	return string(r)
}

// ComputeType
type ComputeType interface {
	ComputeType() string
}

func NewComputeType(v string) (ComputeType, error) {
	if v == "" {
		return nil, errors.New("empty compute type")
	}

	b := v == computeTypeAscend ||
		v == computeTypeMPI
	if !b {
		return nil, errors.New("invalid compute type")
	}

	return computeType(v), nil
}

type computeType string

func (r computeType) ComputeType() string {
	return string(r)
}

// ComputeVersion
type ComputeVersion interface {
	ComputeVersion() string
}

func NewComputeVersion(v string) (ComputeVersion, error) {
	if v == "" {
		return nil, errors.New("empty compute version")
	}

	b := v == computeVersionCudaMS13 ||
		v == computeVersionCannMS17 ||
		v == computeVersionCannMS19_1 ||
		v == computeVersionCannMS19_2

	if !b {
		return nil, errors.New("invalid compute type")
	}

	return computeVersion(v), nil
}

type computeVersion string

func (r computeVersion) ComputeVersion() string {
	return string(r)
}

// ComputeFlavor
type ComputeFlavor interface {
	ComputeFlavor() string
}

func NewComputeFlavorVersion(flaver string, t string, version string) (ComputeFlavor, ComputeVersion, error) {
	if flaver == "" || t == "" {
		return nil, nil, errors.New("empty compute flavor or version")
	}

	b1 := t == computeTypeAscend && flaver == computeFlaverAscend
	b2 := t == computeTypeMPI && flaver == computeFlaverGPU
	if b1 {
		b1 = version == computeVersionCannMS17 ||
			version == computeVersionCannMS19_1 ||
			version == computeVersionCannMS19_2
	}

	if b2 {
		b2 = version == computeVersionCudaMS13
	}

	if b1 || b2 {
		return computeFlavor(flaver), computeVersion(version), nil
	}

	return nil, nil, errors.New("invalid compute flaver or version")
}

type computeFlavor string

func (r computeFlavor) ComputeFlavor() string {
	return string(r)
}

// CustomizedKey
type CustomizedKey interface {
	CustomizedKey() string
}

func NewCustomizedKey(v string) (CustomizedKey, error) {
	if v == "" {
		return nil, errors.New("empty key")
	}

	v = utils.XSSFilter(v)

	if max := 20; utils.StrLen(v) > max {
		return nil, fmt.Errorf("the length of key should be less than %d", max)
	}

	return customizedKey(v), nil
}

type customizedKey string

func (r customizedKey) CustomizedKey() string {
	return string(r)
}

// CustomizedValue
type CustomizedValue interface {
	CustomizedValue() string
}

func NewCustomizedValue(v string) (CustomizedValue, error) {
	v = utils.XSSFilter(v)

	if max := 40; utils.StrLen(v) > max {
		return nil, fmt.Errorf("the length of value should be less than %d", max)
	}

	return customizedValue(v), nil
}

type customizedValue string

func (r customizedValue) CustomizedValue() string {
	return string(r)
}

// InputFilePath
type InputeFilePath interface {
	InputeFilePath() string
}

func NewInputeFilePath(v string) (InputeFilePath, error) {
	v = utils.XSSFilter(v)

	if max := 50; utils.StrLen(v) > max {
		return nil, fmt.Errorf("the length of file path should be less than %d", max)
	}

	return inputeFilePath(v), nil
}

type inputeFilePath string

func (r inputeFilePath) InputeFilePath() string {
	return string(r)
}
