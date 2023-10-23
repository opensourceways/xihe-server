package bigmodels

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	libutils "github.com/opensourceways/community-robot-lib/utils"
	"github.com/opensourceways/xihe-server/bigmodel/domain"
	"github.com/sirupsen/logrus"
)

type skyWorkRequest struct {
	Inputs            string      `json:"inputs"`
	History           [][2]string `json:"history"`
	Sampling          bool        `json:"sampling"`
	TopK              int         `json:"top_k"`
	TopP              float64     `json:"top_p"`
	Temperature       float64     `json:"temperature"`
	RepetitionPenalty float64     `json:"repetition_penalty"`
}

type skyWorkResponse struct {
	Reply        string `json:"reply"`
	Code         int    `json:"code"`
	Msg          string `json:"msg"`
	StreamStatus string `json:"stream_status"`
}

type skyWorkInfo struct {
	endpoints chan string
}

func newSkyWorkInfo(cfg *Config) (info skyWorkInfo, err error) {
	ce := &cfg.Endpoints
	es, _ := ce.parse(ce.SkyWork)

	// init endpoints
	info.endpoints = make(chan string, len(es))
	for _, e := range es {
		info.endpoints <- e
	}

	return
}

func (s *service) SkyWork(ch chan string, input *domain.SkyWorkInput) (err error) {
	// input audit
	if err = s.check.check(input.Text.SkyWorkText()); err != nil {
		return
	}

	// call bigmodel skywork 13b
	f := func(ec chan string, e string) (err error) {
		err = s.genSkyWork(ec, ch, e, input)

		return
	}

	if err = s.doWaitAndEndpointNotReturned(s.skyWorkInfo.endpoints, f); err != nil {
		return
	}

	return
}

func (s *service) genSkyWork(ec, ch chan string, endpoint string, input *domain.SkyWorkInput) (
	err error,
) {
	t, err := genToken(&s.wukongInfo.cfg.CloudConfig)
	if err != nil {
		return
	}

	opt := toSkyWorkReq(input)
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

	req.Header.Set("X-Auth-Token", t)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "*/*")

	resp, err := s.hc.Client.Do(req)
	if err != nil {
		return
	}

	reader := bufio.NewReader(resp.Body)

	var (
		r     skyWorkResponse
		count int
	)
	go func() {
		defer func() { ec <- endpoint }()
		defer resp.Body.Close()

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				ch <- "done"

				logrus.Debugf("skywork read end or error: %s", err)

				return
			}

			data := strings.Replace(string(line), "data: ", "", 1)
			data = strings.TrimRight(data, "\x00")

			if err = json.Unmarshal([]byte(data), &r); err != nil {
				continue
			}

			if r.StreamStatus == "DONE" {
				ch <- "done"

				return
			}

			// response audit, skip 6 response
			if r.Reply != "" && count > 6 {
				count = 0

				if err = s.check.check(r.Reply); err != nil {
					logrus.Debugf("content audit not pass: %s", err.Error())

					ch <- "done"

					return
				}
			}

			ch <- r.Reply
			count += 1
		}
	}()

	return
}

func toSkyWorkReq(input *domain.SkyWorkInput) skyWorkRequest {
	history := make([][2]string, len(input.History))

	for i := range input.History {
		history[i] = input.History[i].History()
	}

	return skyWorkRequest{
		Inputs:            input.Text.SkyWorkText(),
		History:           history,
		Sampling:          input.Sampling,
		TopK:              input.TopK.TopK(),
		TopP:              input.TopP.TopP(),
		Temperature:       input.Temperature.Temperature(),
		RepetitionPenalty: input.RepetitionPenalty.RepetitionPenalty(),
	}
}
