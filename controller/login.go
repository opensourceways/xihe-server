package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/authing"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type oldUserTokenPayload struct {
	Account                 string `json:"account"`
	Email                   string `json:"email"`
	PlatformToken           string `json:"token"`
	PlatformUserNamespaceId string `json:"nid"`
}

func (pl *oldUserTokenPayload) DomainAccount() domain.Account {
	a, _ := domain.NewAccount(pl.Account)

	return a
}

func (pl *oldUserTokenPayload) PlatformUserInfo() platform.UserInfo {
	v, _ := domain.NewEmail(pl.Email)

	return platform.UserInfo{
		User:  pl.DomainAccount(),
		Token: pl.PlatformToken,
		Email: v,
	}
}

func (pl *oldUserTokenPayload) isNotMe(a domain.Account) bool {
	return pl.Account != a.Account()
}

func (pl *oldUserTokenPayload) isMyself(a domain.Account) bool {
	return pl.Account == a.Account()
}

func (pl *oldUserTokenPayload) hasEmail() bool {
	return pl.Email != "" && pl.PlatformToken != ""
}

type newUserTokenPayload struct {
	AccessToken string `json:"access_token"`
}

func AddRouterForLoginController(
	rg *gin.RouterGroup,
	repo repository.User,
	ps platform.User,
	auth authing.User,
	login repository.Login,
	sender message.Sender,
) {
	pc := LoginController{
		auth: auth,
		us:   app.NewUserService(repo, ps, sender),
		ls:   app.NewLoginService(login),
	}

	pc.password, _ = domain.NewPassword(apiConfig.DefaultPassword)

	rg.GET("/v1/login", pc.Login)
	rg.GET("/v1/login/:account", pc.Logout)
}

type LoginController struct {
	baseController

	auth     authing.User
	us       app.UserService
	ls       app.LoginService
	password domain.Password
}

// @Title Login
// @Description callback of authentication by authing
// @Tags  Login
// @Param	code		query	string	true	"authing code"
// @Param	redirect_uri	query	string	true	"redirect uri"
// @Accept json
// @Success 200 {object} app.UserDTO
// @Failure 500 system_error         system error
// @Failure 501 duplicate_creating   create user repeatedly which should not happen
// @Router / [get]
func (ctl *LoginController) Login(ctx *gin.Context) {
	info, err := ctl.auth.GetByCode(
		ctl.getQueryParameter(ctx, "code"),
		ctl.getQueryParameter(ctx, "redirect_uri"),
	)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseCodeError(
			errorSystemError, err,
		))

		return
	}

	user, err := ctl.us.GetByAccount(info.Name)
	if err != nil {
		if d := newResponseError(err); d.Code != errorResourceNotExists {
			ctl.sendRespWithInternalError(ctx, d)

			return
		}

		if user, err = ctl.newUser(ctx, info); err != nil {
			return
		}
	} else {
		if user.Email != info.Email.Email() {
			ctl.us.UpdateBasicInfo(
				info.Name, app.UpdateUserBasicInfoCmd{Email: info.Email},
			)
		}
	}

	if err := ctl.newLogin(ctx, info); err != nil {
		return
	}

	payload := oldUserTokenPayload{
		Account:                 user.Account,
		Email:                   user.Email,
		PlatformToken:           user.Platform.Token,
		PlatformUserNamespaceId: user.Platform.NamespaceId,
	}

	token, err := ctl.newApiToken(ctx, payload)
	if err != nil {
		ctl.sendRespWithInternalError(
			ctx, newResponseCodeError(errorSystemError, err),
		)

		return
	}

	ctl.setRespToken(ctx, token)
	ctx.JSON(http.StatusOK, newResponseData(user))
}

func (ctl *LoginController) newLogin(ctx *gin.Context, info authing.Login) (err error) {
	idToken, err := ctl.encryptData(info.IDToken)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseCodeError(
			errorSystemError, err,
		))

		return
	}

	accessToken, err := ctl.encryptData(info.AccessToken)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseCodeError(
			errorSystemError, err,
		))

		return
	}

	err = ctl.ls.Create(&app.LoginCreateCmd{
		Account:     info.Name,
		Info:        idToken,
		AccessToken: accessToken,
	})
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	}

	return
}

func (ctl *LoginController) newUser(ctx *gin.Context, info authing.Login) (user app.UserDTO, err error) {
	cmd := app.UserCreateCmd{
		Email:    info.Email,
		Account:  info.Name,
		Password: ctl.password,
		Bio:      info.Bio,
		AvatarId: info.AvatarId,
	}

	if user, err = ctl.us.Create(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	}

	return
}

// @Title Logout
// @Description get info of login
// @Tags  Login
// @Param	account	path	string	true	"account"
// @Accept json
// @Success 200 {object} app.LoginDTO
// @Failure 400 bad_request_param   account is invalid
// @Failure 401 not_allowed         can't get login info of other user
// @Failure 500 system_error        system error
// @Router /{account} [get]
func (ctl *LoginController) Logout(ctx *gin.Context) {
	account, err := domain.NewAccount(ctx.Param("account"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if pl.isNotMe(account) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed,
			"can't get login info of other user",
		))

		return
	}

	info, err := ctl.ls.Get(account)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	v, err := ctl.decryptData(info.Info)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseCodeError(
			errorSystemError, err,
		))

		return
	}

	info.Info = string(v)
	ctx.JSON(http.StatusOK, newResponseData(info))
}
