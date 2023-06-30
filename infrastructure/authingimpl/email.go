package authingimpl

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"

	libutils "github.com/opensourceways/community-robot-lib/utils"
	"github.com/sirupsen/logrus"
)

const (
	resetPassword = "CHANNEL_RESET_PASSWORD"
	verifyEmail   = "CHANNEL_VERIFY_EMAIL_LINK"
	changeEmail   = "CHANNEL_UPDATE_EMAIL"
	bindEmail     = "CHANNEL_BIND_EMAIL"
	unbindEmail   = "CHANNEL_UNBIND_EMAIL"

	accountTypeEmail = "email"

	errorCodeError          = "email_code_error"
	errorEmailDuplicateBind = "email_email_duplicate_bind"
	ErrorUserDuplicateBind  = "email_user_duplicate_bind"
	errorEmailDuplicateSend = "email_email_duplicate_send"
)

type managerBody struct {
	AppId     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
	GrantType string `json:"grant_type"`
}

type managerToken struct {
	Status int    `json:"status"`
	Token  string `json:"token"`
	Msg    string `json:"msg"`
}

type normalEmailRes struct {
	Code   int `json:"code"`
	Status int `json:"status"`
}

type sendEmail struct {
	Account             string `json:"account"`
	Channel             string `json:"channel"`
	CaptchaVerification string `json:"captchaVerification"`
}

type veriEmail struct {
	Account     string `json:"account"`
	Code        string `json:"code"`
	UserId      string `json:"user_id"`
	AccountType string `json:"account_type"`
}

func (impl *user) getManagerToken() (token string, err error) {
	b := managerBody{
		AppId:     impl.cfg.APPId,
		AppSecret: impl.cfg.Secret,
		GrantType: "token",
	}

	body, err := libutils.JsonMarshal(&b)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, impl.getManagerTokenURL, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	var res managerToken
	if err = sendHttpRequest(req, &res); err != nil {
		return
	}

	if res.Status != 200 {
		err = errors.New("get token error")
		return
	}

	token = res.Token

	return
}

func (impl *user) SendBindEmail(email, capt string) (code string, err error) {
	token, err := impl.getManagerToken()
	if err != nil {
		return
	}

	return impl.sendEmail(token, bindEmail, email, capt)
}

func (impl *user) sendEmail(token, channel, email, capt string) (code string, err error) {

	send := sendEmail{
		Account:             email,
		Channel:             channel,
		CaptchaVerification: capt,
	}

	body, err := libutils.JsonMarshal(&send)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, impl.sendEmailURL, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Add("token", token)

	var res normalEmailRes
	err = sendHttpRequest(req, &res)

	if res.Status != 200 {
		code = errorReturn(err)
		err = errors.New("send email error")

		return
	}

	return
}

func (impl *user) VerifyBindEmail(email, passCode, userid string) (code string, err error) {
	token, err := impl.getManagerToken()
	if err != nil {
		return
	}

	return impl.verifyBindEmail(token, email, passCode, userid)
}

func (impl *user) verifyBindEmail(token, email, passCode, userid string) (code string, err error) {
	veri := veriEmail{
		Account:     email,
		Code:        passCode,
		UserId:      userid,
		AccountType: accountTypeEmail,
	}

	body, err := libutils.JsonMarshal(&veri)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, impl.bindEmailURL, bytes.NewBuffer(body))
	if err != nil {
		return
	}

	req.Header.Add("token", token)

	var res normalEmailRes
	err = sendHttpRequest(req, &res)

	if res.Code != 200 {
		code = errorReturn(err)
		err = fmt.Errorf("bind email error")

		return
	}

	return
}

func errorReturn(err error) (code string) {
	logrus.Debugf("email error:", err)

	if err == nil {
		return
	}

	errinfo := err.Error()

	if strings.Contains(errinfo, "E0002") {
		code = errorCodeError
		return
	}

	if strings.Contains(errinfo, "E0004") {
		code = errorEmailDuplicateBind
	}

	if strings.Contains(errinfo, "E00016") {
		code = ErrorUserDuplicateBind
		return
	}

	if strings.Contains(errinfo, "E00049") {
		code = errorEmailDuplicateSend
		return
	}

	return
}
