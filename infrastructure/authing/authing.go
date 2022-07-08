package authing

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Authing/authing-go-sdk/lib/authentication"
	"github.com/Authing/authing-go-sdk/lib/management"
	"github.com/Authing/authing-go-sdk/lib/model"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

var cli *management.Client
var authingCfg *config.AuthingService
var OIDCConfig *oauth2.Config

func Init(cfg *config.AuthingService) error {
	authingCfg = cfg
	cli = management.NewClient(cfg.UserPoolId, cfg.Secret)
	return initAuthing()
}

func NewUserMapper() repositories.UserMapper {
	return userMapper{}
}

type userMapper struct {
}

type AuthingLoginUser struct {
	Birthdate           string `json:"birthdate,omitempty"`
	Gender              string `json:"gender,omitempty"`
	Name                string `json:"name,omitempty"`
	Nickname            string `json:"nickname,omitempty"`
	Picture             string `json:"picture,omitempty"`
	UpdatedAT           string `json:"updated_at,omitempty"`
	Website             string `json:"website,omitempty"`
	ExternalID          string `json:"external_id,omitempty"`
	Sub                 string `json:"sub,omitempty"`
	Email               string `json:"email,omitempty"`
	EmailVerified       bool   `json:"email_verified,omitempty"`
	PhoneNumber         string `json:"phone_number,omitempty"`
	PhoneNumberVerified bool   `json:"phone_number_verified,omitempty"`
}

func (u userMapper) Get(userId string) (do repositories.UserDO, err error) {
	v, err := cli.Detail(userId)
	if err != nil {
		return
	}

	do.Id = userId

	// TODO
	do.Bio = ""

	if v.Email != nil {
		do.Email = *v.Email
	}

	if v.Username != nil {
		do.Account = *v.Username
	}

	if v.Nickname != nil {
		do.Nickname = *v.Nickname
	}

	if v.Photo != nil {
		do.AvatarId = *v.Photo
	}

	if v.Phone != nil {
		do.PhoneNumber = *v.Phone
	}

	return
}

func (u userMapper) Update(do repositories.UserDO) error {
	m := model.UpdateUserInput{}

	//TODO bio
	m.Email = &do.Email
	m.Photo = &do.AvatarId
	m.Phone = &do.PhoneNumber
	m.Nickname = &do.Nickname

	_, err := cli.UpdateUser(do.Id, m)

	return err
}

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		currentUserClient := authentication.NewClient(authingCfg.UserPoolId, authingCfg.Secret)
		if currentUserClient == nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, &struct {
				Code string      `json:"code"`
				Msg  string      `json:"msg"`
				Data interface{} `json:"data"`
			}{
				Code: "500",
				Msg:  "authentication init error",
				Data: "",
			})
			return
		}
		currentUser, err := currentUserClient.GetCurrentUser(&token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, &struct {
				Code string      `json:"code"`
				Msg  string      `json:"msg"`
				Data interface{} `json:"data"`
			}{
				Code: "500",
				Msg:  "authentication init error",
				Data: err.Error(),
			})
			return
		}
		currentUserClient.SetCurrentUser(currentUser)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, &struct {
				Code string      `json:"code"`
				Msg  string      `json:"msg"`
				Data interface{} `json:"data"`
			}{
				Code: "401",
				Msg:  "forbidden access",
				Data: "",
			})
			return
		}
		c.Keys = make(map[string]interface{})
		c.Keys["me"] = currentUserClient
	}
}

func initAuthing() error {
	ctx := context.Background()
	oidcProvider, err := oidc.NewProvider(ctx, authingCfg.AuthingURL+"/oidc")
	if err != nil {
		return err
	}
	OIDCConfig = &oauth2.Config{
		ClientID:     authingCfg.AppID,
		ClientSecret: authingCfg.AppSecret,
		Endpoint:     oidcProvider.Endpoint(),
		RedirectURL:  authingCfg.RedirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "external_id", "phone"},
	}
	return nil
}

func GetUserInfoByToken(access_token string) (userinfo *AuthingLoginUser, err error) {
	resp, err := http.Get(authingCfg.AuthingURL + "/oidc/me?access_token=" + access_token)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respDataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	userinfo = new(AuthingLoginUser)
	err = json.Unmarshal(respDataBytes, userinfo)
	if err != nil {
		return nil, err
	}
	return userinfo, nil

}

func SetCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}

//GetJwtString GetJwtString
func GetJwtString(expire int, id, name, provider string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	now := time.Now()
	claims["exp"] = now.Add(time.Hour * time.Duration(expire)).Unix()
	claims["iat"] = now.Unix()
	claims["id"] = id
	claims["nm"] = name
	claims["p"] = provider
	token.Claims = claims
	tokenString, err := token.SignedString([]byte("xihesdf@#2334sdF"))
	return tokenString, err
}

func GetTokenFromAuthing(code string) (interface{}, error) {
	oauth2Token, err := OIDCConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	userInfo, err := GetUserInfoByToken(oauth2Token.AccessToken)
	if err != nil {

		return nil, err
	}

	token, err := GetJwtString(72, userInfo.Sub, userInfo.Name, userInfo.ExternalID)
	if err != nil {

		return nil, err
	}
	result := &struct {
		AccessToken  string            `json:"accessToken"`
		RefreshToken string            `json:"refreshToken"`
		Token        string            `json:"token"`
		User         *AuthingLoginUser `json:"user"`
	}{}
	result.User = userInfo
	result.AccessToken = oauth2Token.AccessToken
	result.RefreshToken = oauth2Token.RefreshToken
	result.Token = token

	return result, nil
}
