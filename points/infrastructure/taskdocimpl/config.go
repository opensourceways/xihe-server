package taskdocimpl

import (
	"errors"
	"text/template"

	common "github.com/opensourceways/xihe-server/common/domain"
)

type Config struct {
	Templates []templateConfig `json:"templates" required:"true"`
}

func (cfg *Config) Validate() error {
	for i := range cfg.Templates {
		if lang := common.NewLanguage(cfg.Templates[i].Language); lang == nil {
			return errors.New("unsupported language")
		}
	}

	return nil
}

type templateConfig struct {
	File     string `json:"file"  required:"true"`
	Language string `json:"language"  required:"true"`
}

func (cfg *templateConfig) template() (*template.Template, error) {
	return newTemplate(cfg.Language, cfg.File)
}
