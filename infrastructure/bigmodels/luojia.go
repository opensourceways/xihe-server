package bigmodels

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/opensourceways/xihe-server/domain"
)

type luojiaInfo struct {
	bucket    string
	endpoints chan string
}

func newLuoJiaInfo(cfg *Config) luojiaInfo {
	v := luojiaInfo{
		endpoints: make(chan string, 1),
	}

	v.endpoints <- cfg.EndpointsOfLuoJia
	v.bucket = cfg.OBS.LuoJiaBucket

	return v
}

func (s *service) LuoJiaUploadPicture(f io.Reader, user domain.Account) error {
	return s.obs.createObject(
		f,
		s.luojiaInfo.bucket,
		fmt.Sprintf("infer/%s/input.png", user.Account()),
	)
}

func (s *service) LuoJia(question string) (answer string, err error) {
	s.doIfFree(s.luojiaInfo.endpoints, func(e string) error {
		answer, err = s.sendReqToLuojia(e, question)

		return err
	})

	return
}

func (s *service) sendReqToLuojia(endpoint, userName string) (answer string, err error) {
	t, err := s.token()
	if err != nil {
		return
	}

	body := []byte(fmt.Sprintf(`{"user_name":"%s"}`, userName))

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
		Status int    `json:"status"`
	}

	if _, err = s.hc.ForwardTo(req, &r); err != nil {
		return
	}

	if r.Status != 200 {
		err = errors.New("failed")
	} else {
		answer = r.Result
	}

	return
}
