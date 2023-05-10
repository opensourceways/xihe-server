package authingimpl

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

const (
	resetPassword = "CHANNEL_RESET_PASSWORD"
	verifyEmail   = "CHANNEL_VERIFY_EMAIL_LINK"
	changeEmail   = "CHANNEL_UPDATE_EMAIL"
	bindEmail     = "CHANNEL_BIND_EMAIL"
	unbindEmail   = "CHANNEL_UNBIND_EMAIL"
)

type normalEmailRes struct {
	Code string `json:"statusCode"`
	Msg  string `json:"message"`
}

func (impl *user) SendBindEmail(accessToken string) (err error) {
	return impl.sendEmail(accessToken, bindEmail)
}

func (impl *user) sendEmail(accessToken, channel string) (err error) {
	req, err := http.NewRequest(http.MethodPost, impl.sendEmailURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", accessToken)
	req.Header.Add("x-authing-app-id", impl.cfg.APPId)

	var res normalEmailRes
	if err = sendHttpRequest(req, &res); err != nil {
		return
	}

	if res.Code != "200" {
		logrus.Fatalf("send email code, err:%s", res.Msg)
		return errors.New("send email error")
	}

	return
}

func (impl *user) VerifyBindEmail(accessToken, passCode string) (err error) {
	req, err := http.NewRequest(http.MethodPost, impl.BindEmailURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", accessToken)
	req.Header.Add("x-authing-app-id", impl.cfg.APPId)

	var res normalEmailRes
	if err = sendHttpRequest(req, &res); err != nil {
		return
	}

	if res.Code != "200" {
		return fmt.Errorf("%s", res.Msg)
	}

	return
}
