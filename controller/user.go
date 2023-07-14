package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/authing"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	userapp "github.com/opensourceways/xihe-server/user/app"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
	userlogincli "github.com/opensourceways/xihe-server/user/infrastructure/logincli"
)

func AddRouterForUserController(
	rg *gin.RouterGroup,
	repo userrepo.User,
	ps platform.User,
	auth authing.User,
	login app.LoginService,
	sender message.Sender,
) {

	us := userapp.NewUserService(repo, ps, sender, encryptHelperToken)

	ctl := UserController{
		auth: auth,
		repo: repo,
		s:    us,
		email: userapp.NewEmailService(
			auth, userlogincli.NewLoginCli(login),
			us,
		),
	}

	rg.POST("/v1/user", ctl.Create) // TODO: delete
	rg.PUT("/v1/user", ctl.Update)
	rg.GET("/v1/user", ctl.Get)

	rg.POST("/v1/user/following", ctl.AddFollowing)
	rg.DELETE("/v1/user/following/:account", ctl.RemoveFollowing)
	rg.GET("/v1/user/following/:account", ctl.ListFollowing)

	rg.GET("/v1/user/follower/:account", ctl.ListFollower)

	rg.GET("/v1/user/:account/gitlab", checkUserEmailMiddleware(&ctl.baseController), ctl.GitlabToken)
	rg.GET("/v1/user/:account/gitlab/refresh", checkUserEmailMiddleware(&ctl.baseController), ctl.RefreshGitlabToken)

	// email
	rg.GET("/v1/user/check_email", checkUserEmailMiddleware(&ctl.baseController))
	rg.POST("/v1/user/email/sendbind", ctl.SendBindEmail)
	rg.POST("/v1/user/email/bind", ctl.BindEmail)
}

type UserController struct {
	baseController

	repo  userrepo.User
	auth  authing.User
	s     userapp.UserService
	email userapp.EmailService
}

// @Summary		Create
// @Description	create user
// @Tags			User
// @Param			body	body	userCreateRequest	true	"body of creating user"
// @Accept			json
// @Success		201	{object}			app.UserDTO
// @Failure		400	bad_request_body	can't	parse		request	body
// @Failure		400	bad_request_param	some	parameter	of		body	is	invalid
// @Failure		500	system_error		system	error
// @Failure		500	duplicate_creating	create	user	repeatedly
// @Router			/v1/user [post]
func (ctl *UserController) Create(ctx *gin.Context) {
	token, _, err := ctl.getToken(ctx)

	if err != nil || token != apiConfig.DefaultPassword {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed, "not allow",
		))

		return
	}

	req := userCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	d, err := ctl.s.Create(&cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	token, csrftoken, err := ctl.newApiToken(ctx, oldUserTokenPayload{
		Account:                 d.Account,
		Email:                   d.Email,
		PlatformToken:           d.Platform.Token,
		PlatformUserNamespaceId: d.Platform.NamespaceId,
	})
	if err != nil {
		ctl.sendRespWithInternalError(
			ctx, newResponseCodeError(errorSystemError, err),
		)

		return
	}

	ctl.setRespToken(ctx, token, csrftoken, d.Account)
	ctx.JSON(http.StatusCreated, newResponseData(d))
}

// @Summary		Update
// @Description	update user basic info
// @Tags			User
// @Param			body	body	userBasicInfoUpdateRequest	true	"body of updating user"
// @Accept			json
// @Produce		json
// @Router			/v1/user [put]
func (ctl *UserController) Update(ctx *gin.Context) {
	m := userBasicInfoUpdateRequest{}

	if err := ctx.ShouldBindJSON(&m); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := m.toCmd()
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

	if err := ctl.s.UpdateBasicInfo(pl.DomainAccount(), cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(m))
}

// @Summary		Get
// @Description	get user
// @Tags			User
// @Param			account	query	string	false	"account"
// @Accept			json
// @Success		200	{object}			userDetail
// @Failure		400	bad_request_param	account	is		invalid
// @Failure		401	resource_not_exists	user	does	not	exist
// @Failure		500	system_error		system	error
// @Router			/v1/user [get]
func (ctl *UserController) Get(ctx *gin.Context) {
	var target domain.Account

	if account := ctl.getQueryParameter(ctx, "account"); account != "" {
		v, err := domain.NewAccount(account)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}

		target = v
	}

	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	resp := func(u *userapp.UserDTO, isFollower bool) {
		ctx.JSON(http.StatusOK, newResponseData(
			userDetail{
				UserDTO:    u,
				IsFollower: isFollower,
			}),
		)
	}

	if visitor {
		if target == nil {
			ctx.JSON(http.StatusOK, newResponseData(nil))
			return
		}

		// get by visitor
		if u, err := ctl.s.GetByAccount(target); err != nil {
			ctl.sendRespWithInternalError(ctx, newResponseError(err))
		} else {
			u.Email = ""
			resp(&u, false)
		}

		return
	}

	if target != nil && pl.isNotMe(target) {
		// get by follower, and pl.Account is follower
		if u, isFollower, err := ctl.s.GetByFollower(target, pl.DomainAccount()); err != nil {
			ctl.sendRespWithInternalError(ctx, newResponseError(err))
		} else {
			u.Email = ""
			resp(&u, isFollower)
		}

		return
	}

	// get user own info
	if u, err := ctl.s.GetByAccount(pl.DomainAccount()); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		resp(&u, true)
	}
}

// @Title			RefreshGitlabToken
// @Description	refresh platform token of user
// @Tags			User
// @Param			account	path	string	true	"account"
// @Accept			json
// @Success		200	{object}			success
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Router			/{account}/gitlab/refresh [get]
func (ctl *UserController) RefreshGitlabToken(ctx *gin.Context) {
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
			"can't refresh token of other user",
		))

		return
	}

	user, _ := ctl.s.GetByAccount(pl.DomainAccount())

	cmd := userapp.RefreshTokenCmd{
		Account:     pl.DomainAccount(),
		Id:          user.Platform.UserId,
		NamespaceId: user.Platform.NamespaceId,
	}

	if err := ctl.s.RefreshGitlabToken(&cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(err))

		return
	}

	usernew, err := ctl.s.GetByAccount(pl.DomainAccount())

	// create new token
	f := func() (token, csrftoken string) {

		if err != nil {
			return
		}

		payload := oldUserTokenPayload{
			Account:                 usernew.Account,
			Email:                   usernew.Email,
			PlatformToken:           usernew.Platform.Token,
			PlatformUserNamespaceId: usernew.Platform.NamespaceId,
		}

		token, csrftoken, err = ctl.newApiToken(ctx, payload)
		if err != nil {
			return
		}

		return
	}

	token, csrftoken := f()

	if token != "" {
		ctl.setRespToken(ctx, token, csrftoken, usernew.Account)
	}

	ctl.sendRespOfPost(ctx, "success")
}

// @Title			GitLabToken
// @Description	get code platform info of user
// @Tags			User
// @Param			account	path	string	true	"account"
// @Accept			json
// @Success		200	{object}			platformInfo
// @Failure		400	bad_request_param	account	is	invalid
// @Failure		401	not_allowed			can't	get	info	of	other	user
// @Router			/{account}/gitlab [get]
func (ctl *UserController) GitlabToken(ctx *gin.Context) {
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
			"can't get token of other user",
		))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(platformInfo{pl.PlatformToken}))
}

type platformInfo struct {
	Token string `json:"token"`
}

// @Title			CheckEmail
// @Description	check user email
// @Tags			User
// @Accept			json
// @Success		200
// @Failure		400	no	email	this	api	need	email	of	user"
// @Router			/v1/user/check_email [get]
func (ctl *UserController) CheckEmail(ctx *gin.Context) {
	ctl.sendRespOfGet(ctx, "")
}

// @Summary		SendBindEmail
// @Description	send code to user
// @Tags			User
// @Accept			json
// @Success		201	{object}			app.UserDTO
// @Failure		500	system_error		system	error
// @Failure		500	duplicate_creating	create	user	repeatedly
// @Router			/v1/user/email/sendbind [post]
func (ctl *UserController) SendBindEmail(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	req := EmailSend{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := req.toCmd(pl.DomainAccount())
	if err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	if code, err := ctl.email.SendBindEmail(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, "success")
	}
}

// @Summary		BindEmail
// @Description	bind email according the code
// @Tags			User
// @Param			body	body	userCreateRequest	true	"body of creating user"
// @Accept			json
// @Success		201	{object}			app.UserDTO
// @Failure		400	bad_request_body	can't	parse		request	body
// @Failure		400	bad_request_param	some	parameter	of		body	is	invalid
// @Failure		500	system_error		system	error
// @Failure		500	duplicate_creating	create	user	repeatedly
// @Router			/v1/user/email/bind [post]
func (ctl *UserController) BindEmail(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	req := EmailCode{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := req.toCmd(pl.DomainAccount())
	if err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	// create new token
	f := func() (token, csrftoken string) {
		user, err := ctl.s.GetByAccount(pl.DomainAccount())
		if err != nil {
			return
		}

		payload := oldUserTokenPayload{
			Account:                 user.Account,
			Email:                   user.Email,
			PlatformToken:           user.Platform.Token,
			PlatformUserNamespaceId: user.Platform.NamespaceId,
		}

		token, csrftoken, err = ctl.newApiToken(ctx, payload)
		if err != nil {
			return
		}

		return
	}

	if code, err := ctl.email.VerifyBindEmail(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		token, csrftoken := f()
		if token != "" {
			ctl.setRespToken(ctx, token, csrftoken, pl.Account)
		}

		ctl.sendRespOfPost(ctx, "success")
	}
}
