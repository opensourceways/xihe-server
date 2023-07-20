package domain

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/opensourceways/xihe-server/utils"
)

const rootDirectory = ""

var (
	reDirectory = regexp.MustCompile("^[a-zA-Z0-9_/-]+$")
	reFile      = regexp.MustCompile("^[a-zA-Z0-9_.-]+$")
)

// TrainingName
type TrainingName interface {
	TrainingName() string
}

func NewTrainingName(v string) (TrainingName, error) {
	max := DomainConfig.MaxTrainingNameLength
	min := DomainConfig.MinTrainingNameLength

	if n := len(v); n > max || n < min {
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

func NewComputeFlavor(v string) (ComputeFlavor, error) {
	if v == "" {
		return nil, errors.New("empty compute flavor")
	}

	return computeFlavor(v), nil
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
	return customizedValue(v), nil
}

type customizedValue string

func (r customizedValue) CustomizedValue() string {
	return string(r)
}
