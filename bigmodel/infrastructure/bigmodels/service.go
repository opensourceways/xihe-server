package bigmodels

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/bigmodel/domain/bigmodel"
)

var fm *service

func Init(cfg *Config) error {
	obs, err := initOBS(&cfg.OBS.OBSAuthInfo)
	if err != nil {
		return err
	}

	check := initTextCheck(&cfg.Moderation)

	fm = &service{
		obs:   obs,
		check: check,
		cfg:   cfg.Cloud,
		hc:    utils.NewHttpClient(3),
	}

	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	fm.vqaInfo = newVQAInfo(cfg)
	fm.luojiaInfo = newLuoJiaInfo(cfg)
	fm.aiDetectorInfo = newAIDetectorInfo(cfg)
	fm.wukongInfo, err = newWuKongInfo(cfg)
	fm.baichuanInfo, err = newBaiChuanInfo(cfg)
	fm.glm2Info, err = newGLM2Info(cfg)
	fm.llama2Info, err = newLLAMA2Info(cfg)

	return err
}

func NewBigModelService() bigmodel.BigModel {
	return fm
}

type service struct {
	cfg   CloudConfig
	obs   obsService
	check textCheckService

	hc utils.HttpClient

	vqaInfo        vqaInfo
	wukongInfo     wukongInfo
	luojiaInfo     luojiaInfo
	aiDetectorInfo aiDetectorInfo
	baichuanInfo   baichuanInfo
	glm2Info       glm2Info
	llama2Info     llama2Info
}

func (s *service) token() (string, error) {
	return genToken(&s.cfg)
}

func genToken(cfg *CloudConfig) (string, error) {
	str := `
{
    "auth":{
       "identity":{
          "methods":[
             "password"
          ],
          "password":{
             "user":{
                "name":"%v",
                "password":"%v",
                "domain":{
                   "name":"%v"
                }
             }
          }
       },
       "scope":{
          "project":{
             "name":"%s"
          }
       }
    }
}
	`

	body := fmt.Sprintf(
		str, cfg.User, cfg.Password, cfg.Domain, cfg.Project,
	)

	resp, err := http.Post(
		cfg.AuthEndpoint, "application/json",
		strings.NewReader(body),
	)
	if err != nil {
		return "", err
	}

	t := resp.Header.Get("x-subject-token")

	if err = resp.Body.Close(); err != nil {
		return "", err
	}

	return t, nil
}

func (s *service) doIfFree(
	ec chan string,
	f func(string) error,
) error {
	select {
	case e := <-ec:
		err := f(e)
		ec <- e

		return err

	default:
		return bigmodel.NewErrorBusySource(errors.New("access overload, please try again later"))
	}
}

func (s *service) doWaitAndEndpointNotReturned(
	ec chan string,
	f func(chan string, string) error,
) error {
	select {
	case e := <-ec:
		return f(ec, e)
	case <-time.After(2 * 60 * time.Second):
		return bigmodel.NewErrorBusySource(errors.New("access overload, please try again later"))
	}
}

func (s *service) GetIdleEndpoint(bid string) (int, error) {
	switch bid {
	case "wukong":
		return len(s.wukongInfo.endpoints), nil
	case "wukong_4img":
		return len(s.wukongInfo.endpoints4), nil
	default:
		return 0, errors.New("internal error, cannot found this bigmodel")
	}
}
