package bigmodels

import (
	"bytes"
	"errors"
	"net/http"
	"strings"

	"github.com/opensourceways/community-robot-lib/utils"
	"github.com/opensourceways/xihe-server/bigmodel/domain"
)

type aiDetectorInfo struct {
	endpoint string
}

type aiDetectorReq struct {
	Lang string `json:"lang"`
	Text string `json:"text"`
}

type aiDetectorResp struct {
	Code           int    `json:"code"`
	Msg            string `json:"msg"`
	MachineWritten bool   `json:"machine_written"`
}

func newAIDetectorInfo(cfg *Config) aiDetectorInfo {
	e := &cfg.Endpoints

	es, _ := e.parse(e.AIDetector)

	return aiDetectorInfo{
		endpoint: es[0],
	}
}

func (s *service) AIDetector(input domain.AIDetectorInput) (ismachine bool, err error) {
	resp, err := s.aiDetector(input)
	if err != nil {
		if strings.Contains(err.Error(), "error_code") {
			err = NewErrorConcurrentRequest(err)
		}

		return
	}

	if resp.Code != 200 {
		err = errors.New(resp.Msg)

		return
	}

	ismachine = resp.MachineWritten

	return
}

func (s *service) aiDetector(input domain.AIDetectorInput) (resp aiDetectorResp, err error) {
	t, err := genToken(&s.wukongInfo.cfg.CloudConfig)
	if err != nil {
		return
	}

	in := aiDetectorReq{
		Lang: input.Lang.Lang(),
		Text: input.Text.AIDetectorText(),
	}

	body, err := utils.JsonMarshal(&in)
	if err != nil {
		return
	}

	es := s.aiDetectorInfo.endpoint
	req, err := http.NewRequest(
		http.MethodPost, es, bytes.NewBuffer(body),
	)
	if err != nil {
		return
	}

	req.Header.Set("X-Auth-Token", t)
	req.Header.Set("Content-Type", "application/json")

	_, err = s.hc.ForwardTo(req, &resp)
	if err != nil {
		return
	}

	return
}
