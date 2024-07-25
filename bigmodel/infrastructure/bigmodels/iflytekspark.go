package bigmodels

import (
	"bufio"
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"

	libutils "github.com/opensourceways/community-robot-lib/utils"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
)

const (
	skipStepFly = 5

	doneStatusFly = "DONE"
	lenThreshold  = 500
)

type iflyteksparkRequest struct {
	Inputs            string  `json:"question"`
	Sampling          bool    `json:"do_sample"`
	TopK              int     `json:"top_k"`
	Temperature       float64 `json:"temperature"`
	RepetitionPenalty float64 `json:"repetition_penalty"`
}

type iflyteksparkResponse struct {
	Reply        string `json:"reply"`
	Code         int    `json:"code"`
	Msg          string `json:"msg"`
	StreamStatus string `json:"stream_status"`
}

type iflyteksparkInfo struct {
	auth CloudConfig

	endpoints     chan string
	endpointsLong chan string
}

func newiflyteksparkInfo(cfg *Config) (info iflyteksparkInfo, err error) {
	ce := &cfg.Endpoints
	es, err := ce.parse(ce.IFlytekspark)
	if err != nil {
		return
	}
	esLong, err := ce.parse(ce.IFlyteksparkLong)
	if err != nil {
		return
	}

	info.auth = cfg.CloudGY

	// init endpoints
	info.endpoints = make(chan string, len(es))
	for _, e := range es {
		info.endpoints <- e
	}

	info.endpointsLong = make(chan string, len(esLong))
	for _, e := range esLong {
		info.endpointsLong <- e
	}

	return
}

func (s *service) IFlytekSpark(ch chan string, input *domain.IFlytekSparkInput) (err error) {
	// input audit
	if err = s.check.check(input.Text.IFlytekSparkText()); err != nil {
		return
	}

	// call bigmodel iflytekspark
	f := func(ec chan string, e string) (err error) {
		err = s.geniflytekspark(ec, ch, e, input)

		return
	}

	if utf8.RuneCountInString(input.Text.IFlytekSparkText()) > lenThreshold {
		err = s.doWaitAndEndpointNotReturned(s.iflyteksparkInfo.endpointsLong, f)
		return
	}

	err = s.doWaitAndEndpointNotReturned(s.iflyteksparkInfo.endpoints, f)

	return
}

func (s *service) geniflytekspark(ec, ch chan string, endpoint string, input *domain.IFlytekSparkInput) (
	err error,
) {
	t, err := genToken(&s.iflyteksparkInfo.auth)
	if err != nil {
		return
	}

	opt := toiflyteksparkReq(input)
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
		r     iflyteksparkResponse
		count int
	)

	go func() {
		defer close(ch)
		defer func() { ec <- endpoint }()
		defer resp.Body.Close()

		for {
			line, err := reader.ReadString('\n')

			if count != 1 && err != nil {
				ch <- "done"

				return
			}

			data := strings.Replace(string(line), "data: ", "", 1)
			data = strings.TrimRight(data, "\n")

			if err = json.Unmarshal([]byte(data), &r); err != nil {
				continue
			}

			if r.StreamStatus == doneStatusFly {
				ch <- "done"

				return
			}

			if r.Reply != "" && count > skipStepFly {
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

func toiflyteksparkReq(input *domain.IFlytekSparkInput) iflyteksparkRequest {
	history := make([][2]string, len(input.History))

	for i := range input.History {
		history[i] = input.History[i].History()
	}

	return iflyteksparkRequest{
		Inputs:            input.Text.IFlytekSparkText(),
		Sampling:          input.Sampling,
		TopK:              input.TopK.TopK(),
		Temperature:       input.Temperature.Temperature(),
		RepetitionPenalty: input.RepetitionPenalty.RepetitionPenalty(),
	}
}
