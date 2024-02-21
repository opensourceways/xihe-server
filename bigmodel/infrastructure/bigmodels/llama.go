package bigmodels

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	libutils "github.com/opensourceways/community-robot-lib/utils"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
)

const (
	skipStepLlama = 10

	doneStatusLlama      = "DONE"
	replaceResponseLlama = "data: "
)

type llama2Request struct {
	Inputs            string      `json:"inputs"`
	History           [][2]string `json:"history"`
	Sampling          bool        `json:"sampling"`
	TopK              int         `json:"top_k"`
	TopP              float64     `json:"top_p"`
	Temperature       float64     `json:"temperature"`
	RepetitionPenalty float64     `json:"repetition_penalty"`
}

type llama2Response struct {
	Reply        string `json:"reply"`
	Code         int    `json:"code"`
	Msg          string `json:"msg"`
	StreamStatus string `json:"stream_status"`
}

type llama2Info struct {
	endpoints chan string
}

func newLLAMA2Info(cfg *Config) (info llama2Info, err error) {
	ce := &cfg.Endpoints
	es, _ := ce.parse(ce.LLAMA2)

	// init endpoints
	info.endpoints = make(chan string, len(es))
	for _, e := range es {
		info.endpoints <- e
	}

	return
}

func (s *service) LLAMA2(ch chan string, input *domain.LLAMA2Input) (err error) {
	// input audit
	if err = s.check.check(input.Text.LLAMA2Text()); err != nil {
		logrus.Debugf("content audit not pass: %s", err.Error())
		return
	}

	// call bigmodel llama2
	f := func(ec chan string, e string) (err error) {
		err = s.genllama2(ec, ch, e, input)

		return
	}

	if err = s.doWaitAndEndpointNotReturned(s.llama2Info.endpoints, f); err != nil {
		return
	}

	return
}

func (s *service) genllama2(ec, ch chan string, endpoint string, input *domain.LLAMA2Input) (
	err error,
) {
	t, err := genToken(&s.wukongInfo.cfg.CloudConfig)
	if err != nil {
		return
	}

	opt := toLLAMA2Req(input)
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
		r     llama2Response
		count int
	)
	go func() {
		defer func() { ec <- endpoint }()
		defer resp.Body.Close()
		defer close(ch)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				ch <- "done"

				logrus.Debugf("llama read end or error: %s", err)

				return
			}

			data := strings.Replace(string(line), replaceResponseLlama, "", 1)
			data = strings.TrimRight(data, "\x00")

			if err = json.Unmarshal([]byte(data), &r); err != nil {
				continue
			}

			if r.StreamStatus == doneStatusLlama {
				ch <- "done"

				return
			}

			if r.Reply != "" && count > skipStepLlama {
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

func toLLAMA2Req(input *domain.LLAMA2Input) llama2Request {
	history := make([][2]string, len(input.History))

	for i := range input.History {
		history[i] = input.History[i].History()
	}

	return llama2Request{
		Inputs:            input.Text.LLAMA2Text(),
		History:           history,
		Sampling:          input.Sampling,
		TopK:              input.TopK.TopK(),
		TopP:              input.TopP.TopP(),
		Temperature:       input.Temperature.Temperature(),
		RepetitionPenalty: input.RepetitionPenalty.RepetitionPenalty(),
	}
}
