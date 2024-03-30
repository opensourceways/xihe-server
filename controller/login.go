package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/authing"
	"github.com/opensourceways/xihe-server/domain/platform"
	userapp "github.com/opensourceways/xihe-server/user/app"
	"github.com/opensourceways/xihe-server/utils"
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

func AddRouterForLoginController(
	rg *gin.RouterGroup,
	us userapp.UserService,
	auth authing.User,
	login app.LoginService,
) {
	pc := LoginController{
		auth: auth,
		us:   us,
		ls:   login,
	}

	pc.password, _ = domain.NewPassword(apiConfig.DefaultPassword)

	rg.GET("/v1/login", pc.Login)
	rg.GET("/v1/login/:account", pc.Logout)
	rg.PUT("/v1/signin", pc.SignIn)
}

type LoginController struct {
	baseController

	auth     authing.User
	us       userapp.UserService
	ls       app.LoginService
	password domain.Password
}

// @Title			Login
// @Description	callback of authentication by authing
// @Tags			Login
// @Param			code			query	string	true	"authing code"
// @Param			redirect_uri	query	string	true	"redirect uri"
// @Accept			json
// @Success		200	{object}			app.UserDTO
// @Failure		500	system_error		system	error
// @Failure		501	duplicate_creating	create	user	repeatedly	which	should	not	happen
// @Router			/v1/login [get]
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
	defer utils.ClearStringMemory(info.AccessToken)

	user, err := ctl.us.GetByAccount(info.Name)
	if err != nil {
		if d := newResponseError(err); d.Code != errorResourceNotExists {
			ctl.sendRespWithInternalError(ctx, d)

			return
		}

		if user, err = ctl.newUser(ctx, info); err != nil {
			utils.DoLog(user.Id, user.Account, "logup", "", "failed")

			return
		}

		utils.DoLog(user.Id, user.Account, "logup", "", "success")
	}

	prepareOperateLog(ctx, user.Account, OPERATE_TYPE_USER, "user login")

	if err := ctl.newLogin(ctx, info); err != nil {
		return
	}

	payload := oldUserTokenPayload{
		Account:                 user.Account,
		Email:                   user.Email,
		PlatformToken:           user.Platform.Token,
		PlatformUserNamespaceId: user.Platform.NamespaceId,
	}

	token, csrftoken, err := ctl.newApiToken(ctx, payload)
	if err != nil {
		ctl.sendRespWithInternalError(
			ctx, newResponseCodeError(errorSystemError, err),
		)

		return
	}

	if err = ctl.setRespToken(ctx, token, csrftoken, user.Account, apiConfig.SessionDomain); err != nil {
		ctl.sendRespWithInternalError(
			ctx, newResponseCodeError(errorSystemError, err),
		)

		return
	}

	utils.DoLog(user.Id, user.Account, "login", "", "success")

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

	err = ctl.ls.Create(&app.LoginCreateCmd{
		Account: info.Name,
		Email:   info.Email,
		Info:    idToken,
		UserId:  info.UserId,
	})
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	}

	return
}

func (ctl *LoginController) newUser(ctx *gin.Context, info authing.Login) (user userapp.UserDTO, err error) {
	cmd := userapp.UserCreateCmd{
		Email:    info.Email,
		Account:  info.Name,
		Password: ctl.password,
		Bio:      info.Bio,
		AvatarId: info.AvatarId,
	}

	if user, err = ctl.us.Create(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	if cmd.Email.Email() != "" {
		if err = ctl.newPlateformAccount(&cmd); err != nil {
			ctl.sendRespWithInternalError(ctx, newResponseError(err))

			return
		}

		// update userdto
		user, err = ctl.us.GetByAccount(cmd.Account)
		if err != nil {
			ctl.sendRespWithInternalError(ctx, newResponseError(err))

			return
		}
	}

	return
}

func (ctl *LoginController) newPlateformAccount(cmd *userapp.UserCreateCmd) (err error) {
	in := userapp.CreatePlatformAccountCmd{
		Email:    cmd.Email,
		Account:  cmd.Account,
		Password: cmd.Password,
	}

	err = ctl.us.NewPlatformAccountWithUpdate(&in)

	return
}

// @Title			Logout
// @Description	get info of login
// @Tags			Login
// @Param			account	path	string	true	"account"
// @Accept			json
// @Success		200	{object}			app.LoginDTO
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	login	info	of	other	user
// @Failure		500	system_error		system	error
// @Router			/v1/login/{account} [get]
func (ctl *LoginController) Logout(ctx *gin.Context) {
	account, err := domain.NewAccount(ctx.Param("account"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, _, ok := ctl.checkUserApiTokenNoRefresh(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "user logout")

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

		utils.DoLog(info.UserId, "", "logout", "", "failed")

		return
	}

	v, err := ctl.decryptData(info.Info)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseCodeError(
			errorSystemError, err,
		))

		return
	}

	utils.DoLog(info.UserId, "", "logout", "", "success")

	ctl.cleanCookie(ctx, apiConfig.SessionDomain)

	info.Info = string(v)
	ctx.JSON(http.StatusOK, newResponseData(info))
}

// @Title			SignIn
// @Description		user sign in
// @Tags			Login
// @Accept			json
// @Success		202
// @Failure		500	system_error		system	error
// @Router			/v1/signin [put]
func (ctl *LoginController) SignIn(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiTokenNoRefresh(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "user sign in")

	if err := ctl.ls.SignIn(pl.DomainAccount()); err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfPut(ctx, "success")
	}
}
