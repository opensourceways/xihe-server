package app

import "fmt"

type Config struct {
	CourseRevision   string `json:"course_revision"`
	UserRevision     string `json:"user_revision"`
	FinetuneRevision string `json:"finetune_revision"`
}

type AgreementType = string

const (
	Course   AgreementType = "course"
	Finetune AgreementType = "finetune"
	User     AgreementType = "user"
)

var (
	_course_ver   string
	_user_ver     string
	_finetune_ver string
)

func Init(c *Config) {
	_course_ver = c.CourseRevision
	_user_ver = c.UserRevision
	_finetune_ver = c.FinetuneRevision
}

func GetCurrentCourseAgree() string {
	return _course_ver
}

func GetCurrentFinetuneAgree() string {
	return _finetune_ver
}

func GetCurrentUserAgree() string {
	return _user_ver
}

func (cfg *Config) Validate() error {
	if cfg.CourseRevision == "" {
		return fmt.Errorf("course revision not initialized")
	}

	if cfg.FinetuneRevision == "" {
		return fmt.Errorf("finetune revision not initialized")
	}

	if cfg.UserRevision == "" {
		return fmt.Errorf("user revision not initialized")
	}

	return nil
}
