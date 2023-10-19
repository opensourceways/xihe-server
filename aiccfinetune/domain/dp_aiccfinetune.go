package domain

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"

	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	rootDirectory = ""

	finetuneTaskFinetune  = "finetune"
	finetuneTaskInference = "inference"

	modelNameWukong = "wukong"
)

var (
	reDirectory = regexp.MustCompile("^[a-zA-Z0-9_/-]+$")
	reFile      = regexp.MustCompile("^[a-zA-Z0-9_.-]+$")
)

// FinetuneName
type FinetuneName interface {
	FinetuneName() string
}

func NewFinetuneName(v string) (FinetuneName, error) {
	v = utils.XSSFilter(v)

	max := 30
	min := 3

	if n := utils.StrLen(v); n > max || n < min {
		return nil, fmt.Errorf("name's length should be between %d to %d", min, max)
	}

	if !types.ReName.MatchString(v) {
		return nil, errors.New("invalid name")
	}

	return finetuneName(v), nil
}

type finetuneName string

func (r finetuneName) FinetuneName() string {
	return string(r)
}

// FinetuneDesc
type FinetuneDesc interface {
	FinetuneDesc() string
}

func NewFinetuneDesc(v string) (FinetuneDesc, error) {
	if v == "" {
		return finetuneDesc(v), nil
	}

	v = utils.XSSFilter(v)

	max := 100
	if utils.StrLen(v) > max {
		return nil, fmt.Errorf("the length of desc should be less than %d", max)
	}

	return finetuneDesc(v), nil
}

type finetuneDesc string

func (r finetuneDesc) FinetuneDesc() string {
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

// FinetuneTask
type FinetuneTask interface {
	FinetuneTask() string
}

func NewFinetuneTask(v string) (FinetuneTask, error) {
	b := v == finetuneTaskFinetune ||
		v == finetuneTaskInference
	if !b {
		return nil, fmt.Errorf("invalid task %s", v)
	}

	return finetuneTask(v), nil
}

type finetuneTask string

func (r finetuneTask) FinetuneTask() string {
	return string(r)
}

// ModelName
type ModelName interface {
	ModelName() string
}

func NewModelName(v string) (ModelName, error) {
	b := v == modelNameWukong
	if !b {
		return nil, fmt.Errorf("invalid model name %s", v)
	}

	return modelName(v), nil
}

type modelName string

func (r modelName) ModelName() string {
	return string(r)
}
