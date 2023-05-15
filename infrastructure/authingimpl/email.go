package authingimpl

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"

	libutils "github.com/opensourceways/community-robot-lib/utils"

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
	Code    int    `json:"statusCode"`
	Msg     string `json:"message"`
	ApiCode int    `json:"apiCode"`
}

type sendEmail struct {
	Email   string `json:"email"`
	Channel string `json:"channel"`
}

type veriEmail struct {
	Email    string `json:"email"`
	PassCode string `json:"passCode"`
}

func (impl *user) SendBindEmail(accessToken, email string) (err error) {
	return impl.sendEmail(accessToken, bindEmail, email)
}

func (impl *user) sendEmail(accessToken, channel, email string) (err error) {

	send := sendEmail{
		Email:   email,
		Channel: channel,
	}

	body, err := libutils.JsonMarshal(&send)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, impl.sendEmailURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", accessToken)
	req.Header.Add("x-authing-app-id", impl.cfg.APPId)

	var res normalEmailRes
	if err = sendHttpRequest(req, &res); err != nil {
		return
	}

	if res.Code != 200 {
		logrus.Fatalf("send email code, err:%s", res.Msg)
		return errors.New("send email error")
	}

	return
}

func (impl *user) VerifyBindEmail(accessToken, email, passCode string) (err error) {
	veri := veriEmail{
		Email:    email,
		PassCode: passCode,
	}

	body, err := libutils.JsonMarshal(&veri)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, impl.BindEmailURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", accessToken)
	req.Header.Add("x-authing-app-id", impl.cfg.APPId)

	var res normalEmailRes
	if err = sendHttpRequest(req, &res); err != nil {
		return
	}

	if res.Code != 200 {
		return fmt.Errorf("%s", res.Msg)
	}

	return
}
