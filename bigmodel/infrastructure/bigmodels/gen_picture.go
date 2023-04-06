package bigmodels

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/opensourceways/community-robot-lib/utils"

	types "github.com/opensourceways/xihe-server/domain"
)

type pictureGenInfo struct {
	singlePictures   chan string
	multiplePictures chan string
}

func newPictureGenInfo(cfg *Config) pictureGenInfo {
	e := &cfg.Endpoints

	se, _ := e.parse(e.SinglePicture)
	me, _ := e.parse(e.MultiplePictures)

	v := pictureGenInfo{
		singlePictures:   make(chan string, len(se)),
		multiplePictures: make(chan string, len(me)),
	}

	for _, e := range se {
		v.singlePictures <- e
	}

	for _, e := range me {
		v.multiplePictures <- e
	}

	return v
}

func (s *service) GenPicture(user types.Account, desc string) (string, error) {
	if err := s.check.check(desc); err != nil {
		return "", err
	}

	r := new(singlePicture)

	err := s.genPicture(user, desc, s.pictureGenInfo.singlePictures, r)
	if err != nil {
		return "", err
	}

	return r.picture()
}

func (s *service) GenPictures(user types.Account, desc string) ([]string, error) {
	if err := s.check.check(desc); err != nil {
		return nil, err
	}

	r := new(multiplePictures)

	err := s.genPicture(user, desc, s.pictureGenInfo.multiplePictures, r)
	if err != nil {
		return nil, err
	}

	return r.picture()
}

func (s *service) genPicture(
	user types.Account, desc string,
	ec chan string, result interface{},
) error {
	return s.doIfFree(ec, func(e string) error {
		return s.sendReqToGenPicture(user, e, desc, result)
	})
}

type pictureGenerateOpt struct {
	Desc string `json:"input_text"`
	User string `json:"user_name"`
}

type singlePicture struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Output string `json:"output_image_url"`
}

func (p *singlePicture) picture() (string, error) {
	if p.Status == -1 {
		return "", errors.New(p.Msg)
	}

	return p.Output, nil
}

type multiplePictures struct {
	Status int      `json:"status"`
	Msg    string   `json:"msg"`
	Output []string `json:"output_image_url"`
}

func (p *multiplePictures) picture() ([]string, error) {
	if p.Status == -1 {
		return nil, errors.New(p.Msg)
	}

	return p.Output, nil
}

func (s *service) sendReqToGenPicture(
	user types.Account, endpoint, desc string, r interface{},
) (err error) {
	t, err := s.token()
	if err != nil {
		return
	}

	opt := pictureGenerateOpt{
		Desc: desc,
		User: user.Account(),
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", t)

	_, err = s.hc.ForwardTo(req, r)

	return err
}
