package taskdocimpl

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"text/template"

	common "github.com/opensourceways/xihe-server/common/domain"
	"github.com/opensourceways/xihe-server/points/domain"
)

func Init(cfg *Config) (*taskDocImpl, error) {
	m := map[string]*template.Template{}

	for i := range cfg.Templates {
		item := &cfg.Templates[i]

		v, err := item.template()
		if err != nil {
			return nil, err
		}

		m[item.Language] = v
	}

	return &taskDocImpl{m}, nil
}

type taskDocImpl struct {
	templates map[string]*template.Template
}

func (impl *taskDocImpl) Doc(tasks []domain.Task, lang common.Language) ([]byte, error) {
	t, ok := impl.templates[lang.Language()]
	if !ok {
		return nil, errors.New("unsupported language")
	}

	buf := new(bytes.Buffer)

	data := make([]taskInfo, len(tasks))
	for i := range tasks {
		item := &tasks[i]

		data[i] = taskInfo{
			Name:          item.Name(lang),
			Desc:          item.RuleDesc(lang),
			MaxPoints:     item.MaxPointsDesc(lang),
			PointsPerOnce: item.Rule.PointsPerOnce,
		}
	}

	if err := renderTemplate(t, data, buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type taskInfo struct {
	Name          string
	Desc          string
	MaxPoints     string
	PointsPerOnce int
}

func newTemplate(name, path string) (*template.Template, error) {
	txtStr, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to new template: read template file failed: %s", err.Error())
	}

	tmpl, err := template.New(name).Parse(string(txtStr))
	if err != nil {
		return nil, fmt.Errorf("failed to new template: build template failed: %s", err.Error())
	}

	return tmpl, nil
}

func renderTemplate(tmpl *template.Template, data interface{}, buf *bytes.Buffer) error {
	err := tmpl.Execute(buf, data)
	if err == nil {
		return nil
	}

	return fmt.Errorf("failed to execute template(%s): %s", tmpl.Name(), err.Error())
}
