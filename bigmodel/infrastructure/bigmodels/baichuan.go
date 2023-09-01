package bigmodels

import (
	"bytes"
	"errors"
	"net/http"

	libutils "github.com/opensourceways/community-robot-lib/utils"
	"github.com/opensourceways/xihe-server/bigmodel/domain"
)

type baichuanInfo struct {
	endpoints chan string
}

type baichuanRequest struct {
	Prompt            string  `json:"prompt"`
	Sampling          bool    `json:"sampling"`
	TopK              int     `json:"top_k"`
	TopP              float64 `json:"top_p"`
	Temperature       float64 `json:"temperature"`
	RepetitionPenalty float64 `json:"repetition_penalty"`
}

type baichuanResponse struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result []struct {
		TextGenerationText []string `json:"text_generation_text"`
	} `json:"result"`
}

func (req baichuanResponse) getText() string {
	return req.Result[0].TextGenerationText[0]
}

func newBaiChuanInfo(cfg *Config) (info baichuanInfo, err error) {

	ce := &cfg.Endpoints
	es, _ := ce.parse(ce.BaiChuan)

	// init endpoints
	info.endpoints = make(chan string, len(es))
	for _, e := range es {
		info.endpoints <- e
	}

	return
}

func (s *service) BaiChuan(input *domain.BaiChuanInput) (code, r string, err error) {
	// input check
	if err = s.check.check(input.Text.BaiChuanText()); err != nil {
		code = CodeInputTextAuditError

		return
	}

	// call bigmodel baichuan
	var resp baichuanResponse
	f := func(e string) (err error) {
		resp, err = s.genBaiChuan(e, input)

		return
	}

	if err = s.doIfFree(s.baichuanInfo.endpoints, f); err != nil {
		return
	}

	// output check
	if resp.Code != 200 {
		code = CodeBaiChuanGenerationError
		err = errors.New("bigmodel baichuan inference generation error")

		return
	}

	if err = s.check.check(resp.getText()); err != nil {
		code = CodeOutputTextAuditError

		return
	}

	return "", resp.getText(), nil
}

func (s *service) genBaiChuan(
	endpoint string, d *domain.BaiChuanInput,
) (resp baichuanResponse, err error) {
	t, err := genToken(&s.wukongInfo.cfg.CloudConfig)
	if err != nil {
		return
	}

	opt := toBaiChuanReq(d)
	body, err := libutils.JsonMarshal(&opt)
	if err != nil {
		return
	}

	req, err := http.NewRequest(
		http.MethodPost, endpoint, bytes.NewBuffer(body),
	)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", t)

	if _, err = s.hc.ForwardTo(req, &resp); err != nil {
		return
	}

	return
}

func toBaiChuanReq(d *domain.BaiChuanInput) baichuanRequest {
	return baichuanRequest{
		Prompt:            d.Text.BaiChuanText(),
		Sampling:          d.Sampling,
		TopK:              d.TopK.TopK(),
		TopP:              d.TopP.TopP(),
		Temperature:       d.Temperature.Temperature(),
		RepetitionPenalty: d.RepetitionPenalty.RepetitionPenalty(),
	}
}
