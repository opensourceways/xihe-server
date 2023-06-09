package authingimpl

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/authing"
)

var userInstance *user

type Config struct {
	APPId    string `json:"app_id"        required:"true"`
	Secret   string `json:"secret"        required:"true"`
	Endpoint string `json:"endpoint"      required:"true"`
}

func Init(v *Config) {
	userInstance = &user{
		cfg:                *v,
		tokenURL:           v.Endpoint + "/oidc/token",
		userInfoURL:        v.Endpoint + "/oidc/user",
		getManagerTokenURL: v.Endpoint + "/manager/token",
		sendEmailURL:       v.Endpoint + "/manager/sendcode",
		bindEmailURL:       v.Endpoint + "/manager/bind/account",
	}
}

func NewAuthingUser() *user {
	return userInstance
}

type user struct {
	cfg                Config
	tokenURL           string
	userInfoURL        string
	getManagerTokenURL string
	sendEmailURL       string
	bindEmailURL       string
}

func (impl *user) GetByAccessToken(accessToken string) (userInfo authing.UserInfo, err error) {
	if accessToken == "" {
		err = errors.New("no access token")

		return
	}

	var v struct {
		Name    string `json:"username,omitempty"`
		Picture string `json:"picture,omitempty"`
		Sub     string `json:"sub,omitempty"`
	}

	if err = impl.getUserInfoByAccessToken(accessToken, &v); err != nil {
		return
	}

	if userInfo.Name, err = domain.NewAccount(v.Name); err != nil {
		return
	}

	if userInfo.AvatarId, err = domain.NewAvatarId(v.Picture); err != nil {
		return
	}

	userInfo.UserId = v.Sub

	return
}

func (impl *user) GetByCode(code, redirectURI string) (login authing.Login, err error) {
	var v struct {
		AccessToken string `json:"access_token"`
		IdToken     string `json:"id_token"`
	}

	if err = impl.getAccessTokenByCode(code, redirectURI, &v); err != nil {
		return
	}

	if v.IdToken == "" {
		err = errors.New("no id token")

		return
	}

	info, err := impl.GetByAccessToken(v.AccessToken)
	if err == nil {
		login.IDToken = v.IdToken
		login.UserInfo = info
		login.AccessToken = v.AccessToken
	}

	return
}

func (impl *user) getAccessTokenByCode(code, redirectURI string, result interface{}) error {
	body := map[string]string{
		"client_id":     impl.cfg.APPId,
		"client_secret": impl.cfg.Secret,
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  redirectURI,
	}

	value := make(url.Values)
	for k, v := range body {
		value.Add(k, v)
	}

	req, err := http.NewRequest(
		http.MethodPost, impl.tokenURL,
		strings.NewReader(strings.TrimSpace(value.Encode())),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return sendHttpRequest(req, result)
}

func (impl *user) getUserInfoByAccessToken(accessToken string, result interface{}) error {
	req, err := http.NewRequest(http.MethodGet, impl.userInfoURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", accessToken)

	return sendHttpRequest(req, result)
}

func sendHttpRequest(req *http.Request, result interface{}) error {
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "xihe-server-authing")
	req.Header.Add("content-type", "application/json")

	hc := utils.NewHttpClient(3)

	_, err := hc.ForwardTo(req, result)

	return err
}
