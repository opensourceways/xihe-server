package bigmodels

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type questionOpt struct {
	Picture  string `json:"image_path"`
	Question string `json:"question"`
}

func (q *questionOpt) serialize() ([]byte, error) {
	return utils.JsonMarshal(q)
}

type questionResp struct {
	Status interface{} `json:"status"`
	Msg    string      `json:"msg"`

	Inference struct {
		Instances string `json:"instances"`
	} `json:"inference_result"`
}

func (q *questionResp) answer() (string, error) {
	if status, ok := q.Status.(string); ok && status == "200" {
		return q.Inference.Instances, nil
	}

	return "", errors.New(q.Msg)
}

type questionHFOpt struct {
	File     *bytes.Buffer `json:"file"`
	Question string        `json:"question"`
}

type vqaInfo struct {
	endpoint   string
	endpointHF string
	bucket     string
}

func newVQAInfo(cfg *Config) vqaInfo {
	ce := &cfg.Endpoints

	es, _ := ce.parse(ce.VQA)
	eshf, _ := ce.parse(ce.VQAHF)

	return vqaInfo{
		bucket:     cfg.OBS.VQABucket,
		endpoint:   es[0],
		endpointHF: eshf[0],
	}
}

func (s *service) Ask(q domain.Question, f string) (string, error) {
	if err := s.check.check(q.Question()); err != nil {
		return "", err
	}

	opt := questionOpt{
		Picture:  f,
		Question: q.Question(),
	}

	body, err := opt.serialize()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost, s.vqaInfo.endpoint, bytes.NewBuffer(body),
	)
	if err != nil {
		return "", err
	}

	t, err := s.token()
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", t)

	v := new(questionResp)

	if _, err = s.hc.ForwardTo(req, v); err != nil {
		return "", err
	}

	return v.answer()
}

func (s *service) VQAUploadPicture(f io.Reader, user types.Account, fileName string) error {
	return s.obs.createObject(f, s.vqaInfo.bucket, filepath.Join(user.Account(), fileName))
}

func (s *service) AskHF(f io.Reader, user types.Account, ask string) (string, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	file, err := writer.CreateFormFile("file", "WechatIMG2645.jpeg")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, f)
	if err != nil {
		return "", err
	}

	writer.WriteField("question", ask)

	writer.Close()

	req, err := http.NewRequest(http.MethodPost, s.vqaInfo.endpointHF, buf)
	if err != nil {
		return "", err
	}

	t, err := s.token()
	if err != nil {
		return "", err
	}

	req.Header.Set("X-Auth-Token", t)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	answer := new(questionResp)

	if _, err = s.hc.ForwardTo(req, answer); err != nil {
		return "", err
	}

	return answer.answer()
}
