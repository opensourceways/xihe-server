package bigmodels

import (
	"bytes"
	"fmt"
	"net/http"
)

type panguInfo struct {
	endpoints chan string
}

func newPanGuInfo(cfg *Config) panguInfo {
	ce := &cfg.Endpoints

	es, _ := ce.parse(ce.Pangu)

	v := panguInfo{
		endpoints: make(chan string, len(es)),
	}

	for _, e := range es {
		v.endpoints <- e
	}

	return v
}

func (s *service) PanGu(question string) (answer string, err error) {
	if err = s.check.check(question); err != nil {
		return
	}

	s.doIfFree(s.panguInfo.endpoints, func(e string) error {
		answer, err = s.sendReqToPangu(e, question)

		return err
	})

	return
}

func (s *service) sendReqToPangu(endpoint, question string) (answer string, err error) {
	t, err := s.token()
	if err != nil {
		return
	}

	body := []byte(fmt.Sprintf(`{"question":"%s"}`, question))

	req, err := http.NewRequest(
		http.MethodPost, endpoint, bytes.NewBuffer(body),
	)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", t)

	var r struct {
		Result string `json:"result"`
	}

	if _, err = s.hc.ForwardTo(req, &r); err != nil {
		return
	}

	answer = r.Result

	return
}
