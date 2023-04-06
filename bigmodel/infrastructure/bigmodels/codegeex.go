package bigmodels

import (
	"bytes"
	"net/http"

	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/bigmodel/domain/bigmodel"
)

type codegeexInfo struct {
	endpoints chan string
}

func newCodeGeexInfo(cfg *Config) codegeexInfo {
	ce := &cfg.Endpoints

	es, _ := ce.parse(ce.CodeGeex)

	v := codegeexInfo{
		endpoints: make(chan string, len(es)),
	}

	for _, e := range es {
		v.endpoints <- e
	}

	return v
}

func (s *service) CodeGeex(question *bigmodel.CodeGeexReq) (r bigmodel.CodeGeexResp, err error) {
	if err = s.check.check(question.Content); err != nil {
		return
	}

	s.doIfFree(s.codegeexInfo.endpoints, func(e string) error {
		r, err = s.sendReqToCodeGeex(e, question)

		return nil
	})

	return
}

func (s *service) sendReqToCodeGeex(
	endpoint string, question *bigmodel.CodeGeexReq,
) (r bigmodel.CodeGeexResp, err error) {
	opt := codegeexReq{
		Samples: question.Content,
		Lang:    question.Lang,
	}

	body, err := utils.JsonMarshal(&opt)
	if err != nil {
		return
	}

	req, err := http.NewRequest(
		http.MethodPost, endpoint, bytes.NewBuffer(body),
	)
	if err != nil {
		return
	}

	t, err := s.token()
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", t)

	_, err = s.hc.ForwardTo(req, &r)

	return
}

type codegeexReq struct {
	Samples string `json:"samples"`
	Lang    string `json:"language"`
}
