package models

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/opensourceways/xihe-server/util"

	"github.com/Authing/authing-go-sdk/lib/authentication"
	"github.com/Authing/authing-go-sdk/lib/management"
	"github.com/Authing/authing-go-sdk/lib/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	JwtString = "xihe-sadf43@98524"
)

type CreateUserInput struct {
	Username          *string  `json:"username,omitempty"`
	Email             *string  `json:"email,omitempty"`
	EmailVerified     *bool    `json:"emailVerified,omitempty"`
	Phone             *string  `json:"phone,omitempty"`
	PhoneVerified     *bool    `json:"phoneVerified,omitempty"`
	Unionid           *string  `json:"unionid,omitempty"`
	Openid            *string  `json:"openid,omitempty"`
	Nickname          *string  `json:"nickname,omitempty"`
	Photo             *string  `json:"photo,omitempty"`
	Password          *string  `json:"password,omitempty"`
	RegisterSource    []string `json:"registerSource,omitempty"`
	Browser           *string  `json:"browser,omitempty"`
	Oauth             *string  `json:"oauth,omitempty"`
	LoginsCount       *int64   `json:"loginsCount,omitempty"`
	LastLogin         *string  `json:"lastLogin,omitempty"`
	Company           *string  `json:"company,omitempty"`
	LastIP            *string  `json:"lastIP,omitempty"`
	SignedUp          *string  `json:"signedUp,omitempty"`
	Blocked           *bool    `json:"blocked,omitempty"`
	IsDeleted         *bool    `json:"isDeleted,omitempty"`
	Device            *string  `json:"device,omitempty"`
	Name              *string  `json:"name,omitempty"`
	GivenName         *string  `json:"givenName,omitempty"`
	FamilyName        *string  `json:"familyName,omitempty"`
	MiddleName        *string  `json:"middleName,omitempty"`
	Profile           *string  `json:"profile,omitempty"`
	PreferredUsername *string  `json:"preferredUsername,omitempty"`
	Website           *string  `json:"website,omitempty"`
	Gender            *string  `json:"gender,omitempty"`
	Birthdate         *string  `json:"birthdate,omitempty"`
	Zoneinfo          *string  `json:"zoneinfo,omitempty"`
	Locale            *string  `json:"locale,omitempty"`
	Address           *string  `json:"address,omitempty"`
	Formatted         *string  `json:"formatted,omitempty"`
	StreetAddress     *string  `json:"streetAddress,omitempty"`
	Locality          *string  `json:"locality,omitempty"`
	Region            *string  `json:"region,omitempty"`
	PostalCode        *string  `json:"postalCode,omitempty"`
	Country           *string  `json:"country,omitempty"`
	ExternalId        *string  `json:"externalId,omitempty"`
}
type AuthingKey struct {
	E   string `json:"e"`
	N   string `json:"n"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	Kid string `json:"Kid"`
}

type AuthingJWKS struct {
	Keys []AuthingKey
}

var AuthingClient *management.Client
var AppClient *model.Application
var UserClient *authentication.Client
var AuthingJWKSItem AuthingJWKS

func InitAuthing(userpoolid, secret string) {
	if userpoolid == "" {
		userpoolid = util.GetConfig().AuthingConfig.UserPoolID
		secret = util.GetConfig().AuthingConfig.Secret
	}
	AuthingClient = management.NewClient(userpoolid, secret)
	AppClient, _ = AuthingClient.FindApplicationById(util.GetConfig().AuthingConfig.AppID)
	UserClient = authentication.NewClient(util.GetConfig().AuthingConfig.AppID, util.GetConfig().AuthingConfig.AppSecret)

}
func ParseAuthingUserInput(userinput *CreateUserInput) *model.CreateUserRequest {
	var testuser model.CreateUserRequest
	testuser.UserInfo.Username = userinput.Username
	testuser.UserInfo.Email = userinput.Email
	return &testuser
}

func GetUserInfoByToekn(token string) error {

	resp, err := http.Get("https://openeuler-omni-manager.authing.cn/oidc/me?access_token=" + token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respDataBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_ = respDataBytes
	return nil

}

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		userInfo, err := CheckAuthorization(token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, util.ExportData(util.CodeStatusClientError, "forbidden", err))
			return
		}
		c.Keys = userInfo
	}
}

//GetJwtString GetJwtString
func GetJwtString(expire int, id, name, provider string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	now := time.Now().In(util.CnTime)
	claims["exp"] = now.Add(time.Hour * time.Duration(expire)).Unix()
	claims["iat"] = now.Unix()
	claims["id"] = id
	claims["nm"] = name
	claims["p"] = provider
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(JwtString))
	return tokenString, err
}

//check user token status
func CheckAuthorization(tokenString string) (userInfo map[string]interface{}, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtString), nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		var ok bool
		userInfo, ok = token.Claims.(jwt.MapClaims)
		if ok == false {
			return nil, fmt.Errorf("token无效")
		}
		if userInfo["id"] == nil || userInfo["id"] == "" {
			return nil, fmt.Errorf("token无效,无id")
		}
		expireTime := userInfo["exp"].(float64)
		if int(expireTime) <= int(time.Now().Local().Unix()) {
			return nil, fmt.Errorf("登陆已经过期")
		}
	}
	return
}
