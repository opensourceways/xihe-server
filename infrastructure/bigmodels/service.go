package bigmodels

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/domain/bigmodel"
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
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // #nosec G402 -- this is a false positive
	}

	fm.vqaInfo = newVQAInfo(cfg)
	fm.panguInfo = newPanGuInfo(cfg)
	fm.luojiaInfo = newLuoJiaInfo(cfg)
	fm.codegeexInfo = newCodeGeexInfo(cfg)
	fm.pictureGenInfo = newPictureGenInfo(cfg)
	fm.pictureDescInfo = newPictureDescInfo(cfg)

	fm.wukongInfo, err = newWuKongInfo(cfg)

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

	vqaInfo         vqaInfo
	panguInfo       panguInfo
	wukongInfo      wukongInfo
	luojiaInfo      luojiaInfo
	codegeexInfo    codegeexInfo
	pictureGenInfo  pictureGenInfo
	pictureDescInfo pictureDescInfo
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
		return bigmodel.NewErrorBusySource(errors.New("busy computation resource"))
	}
}
