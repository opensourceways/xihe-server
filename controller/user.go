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

func AddRouterForUserController(
	rg *gin.RouterGroup,
	repo repository.User,
	ps platform.User,
	auth authing.User,
	sender message.Sender,
) {
	pc := UserController{
		auth: auth,
		repo: repo,
		s:    app.NewUserService(repo, ps, sender),
	}

	rg.POST("/v1/user", pc.Create) // TODO: delete
	rg.PUT("/v1/user", pc.Update)
	rg.GET("/v1/user", pc.Get)

	rg.POST("/v1/user/following", pc.AddFollowing)
	rg.DELETE("/v1/user/following/:account", pc.RemoveFollowing)
	rg.GET("/v1/user/following/:account", pc.ListFollowing)

	rg.GET("/v1/user/follower/:account", pc.ListFollower)
	rg.GET("/v1/user/:account/gitlab", pc.GitlabToken)
}

type UserController struct {
	baseController

	repo repository.User
	auth authing.User
	s    app.UserService
}

// @Summary Create
// @Description create user
// @Tags  User
// @Param	body	body 	userCreateRequest	true	"body of creating user"
// @Accept json
// @Success 201 {object} app.UserDTO
// @Failure 400 bad_request_body    can't parse request body
// @Failure 400 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Failure 500 duplicate_creating  create user repeatedly
// @Router /v1/user [post]
func (ctl *UserController) Create(ctx *gin.Context) {
	token := ctx.GetHeader(headerPrivateToken)

	if token != apiConfig.DefaultPassword {
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

	token, err = ctl.newApiToken(ctx, oldUserTokenPayload{
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

	ctl.setRespToken(ctx, token)
	ctx.JSON(http.StatusCreated, newResponseData(d))
}

// @Summary Update
// @Description update user basic info
// @Tags  User
// @Param	body	body 	userBasicInfoUpdateRequest	true	"body of updating user"
// @Accept json
// @Produce json
// @Router /v1/user [put]
func (uc *UserController) Update(ctx *gin.Context) {
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

	if err := uc.s.UpdateBasicInfo(nil, cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(m))
}

// @Summary Get
// @Description get user
// @Tags  User
// @Param	account	query	string	false	"account"
// @Accept json
// @Success 200 {object} userDetail
// @Failure 400 bad_request_param   account is invalid
// @Failure 401 resource_not_exists user does not exist
// @Failure 500 system_error        system error
// @Router /v1/user [get]
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

	resp := func(u *app.UserDTO, isFollower bool) {
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

// @Title GitLabToken
// @Description get code platform info of user
// @Tags  User
// @Param	account	path	string	true	"account"
// @Accept json
// @Success 200 {object} platformInfo
// @Failure 400 bad_request_param   account is invalid
// @Failure 401 not_allowed         can't get info of other user
// @Router /{account}/gitlab [get]
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
