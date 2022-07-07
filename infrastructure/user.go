package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Authing/authing-go-sdk/lib/authentication"
	"github.com/Authing/authing-go-sdk/lib/model"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/domain/entity"
	"github.com/opensourceways/xihe-server/util"
	"go.mongodb.org/mongo-driver/mongo"

	"golang.org/x/oauth2"
)

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

type FindUserRequest struct {
	Email          *string `json:"email,omitempty"`
	Username       *string `json:"username,omitempty"`
	Phone          *string `json:"phone,omitempty"`
	ExternalId     *string `json:"externalId,omitempty"`
	WithCustomData bool    `json:"withCustomData,omitempty"`
}

// var AuthingManageClient *management.Client

// var AppClient *model.Application
var AuthingDefaultUserClient *authentication.Client

// var AuthingJWKSItem AuthingJWKS
var oidcProvider *oidc.Provider
var OIDCConfig *oauth2.Config

type UserRepo struct {
	Collection *mongo.Collection
}

func NewUserRepository(mongodDB *mongo.Database) *UserRepo {
	repo := new(UserRepo)
	repo.Collection = mongodDB.Collection("user")

	return repo
}

func (r *UserRepo) Save(item *entity.User) (*entity.User, error) {

	return nil, nil
}

func InitAuthing() (err error) {
	if util.GetConfig().AuthingConfig.UserPoolID == "" {
		util.GetConfig().AuthingConfig.UserPoolID = os.Getenv("AUTHING_USER_POOL_ID")
	}

	if util.GetConfig().AuthingConfig.AppID == "" {
		util.GetConfig().AuthingConfig.AppID = os.Getenv("AUTHING_APP_ID")
	}
	if util.GetConfig().AuthingConfig.AppSecret == "" {
		util.GetConfig().AuthingConfig.AppSecret = os.Getenv("AUTHING_APP_SECRET")
	}

	if util.GetConfig().AuthingConfig.Secret == "" {
		util.GetConfig().AuthingConfig.Secret = os.Getenv("AUTHING_APP_SECRET")
	}

	// AuthingManageClient = management.NewClient(util.GetConfig().AuthingConfig.UserPoolID, util.GetConfig().AuthingConfig.Secret)
	AuthingDefaultUserClient = authentication.NewClient(util.GetConfig().AuthingConfig.AppID, util.GetConfig().AuthingConfig.AppSecret)
	ctx := context.Background()
	oidcProvider, err = oidc.NewProvider(ctx, util.GetConfig().AuthingConfig.AuthingUrl+"/oidc")
	if err != nil {
		return err
	}

	OIDCConfig = &oauth2.Config{
		ClientID:     util.GetConfig().AuthingConfig.AppID,
		ClientSecret: util.GetConfig().AuthingConfig.AppSecret,
		Endpoint:     oidcProvider.Endpoint(),
		RedirectURL:  util.GetConfig().AuthingConfig.RedirectURI,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "external_id", "phone"},
	}
	return err
}

func GetUserInfoByToken(access_token string) (userinfo *AuthingLoginUser, err error) {
	resp, err := http.Get(util.GetConfig().AuthingConfig.AuthingUrl + "/oidc/me?access_token=" + access_token)
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

func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		currentUserClient := authentication.NewClient(util.GetConfig().AuthingConfig.AppID, util.GetConfig().AuthingConfig.AppSecret)
		currentUser, err := currentUserClient.GetCurrentUser(&token)
		currentUserClient.SetCurrentUser(currentUser)
		// userInfo, err := CheckAuthorization(token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusUnauthorized, util.ExportData(util.CodeStatusClientError, "forbidden", err.Error()))
			return
		}
		c.Keys = make(map[string]any)
		c.Keys["me"] = currentUserClient
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
	tokenString, err := token.SignedString([]byte(entity.JwtString))
	return tokenString, err
}

//check user token status
func CheckAuthorization(tokenString string) (userInfo map[string]interface{}, err error) {
	var token *jwt.Token
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(entity.JwtString), nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		var ok bool
		userInfo, ok = token.Claims.(jwt.MapClaims)
		if !ok {
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
func GetTokenByCode(code string) (jwtToken *jwt.Token, token *entity.TokenItem, err error) {
	resp, err := http.PostForm(util.GetConfig().AuthingConfig.AuthingUrl+"/oidc/token",
		url.Values{
			"code":          {code},
			"client_id":     {util.GetConfig().AuthingConfig.AppID},
			"client_secret": {util.GetConfig().AuthingConfig.AppSecret},
			"grant_type":    {"authorization_code"},
			"redirect_uri":  {util.GetConfig().AuthingConfig.RedirectURI}})

	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	token = new(entity.TokenItem)
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, nil, err
	}

	jwtToken = new(jwt.Token)
	jwtToken.Valid = false
	jwtToken, err = jwt.Parse(token.AccessToken, func(jwtToken *jwt.Token) (interface{}, error) {
		return util.GetConfig().AuthingConfig.AppSecret, nil
	})

	if err != nil {
		return nil, nil, err
	}
	return jwtToken, token, nil
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

func Change2FindUserRequest(localFindUserRequest *FindUserRequest) (result *model.FindUserRequest) {
	result = new(model.FindUserRequest)
	result.Email = localFindUserRequest.Email
	result.ExternalId = localFindUserRequest.ExternalId
	result.Phone = localFindUserRequest.Phone
	result.Username = localFindUserRequest.Username
	result.WithCustomData = localFindUserRequest.WithCustomData
	return
}
