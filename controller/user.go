package controller

import (
	"net/http"

	"github.com/Authing/authing-go-sdk/lib/authentication"
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/infrastructure/authing"
	"github.com/opensourceways/xihe-server/utils"
)

func AddRouterForUserController(
	rg *gin.RouterGroup,
	repoUser repository.User,
) {
	pc := UserController{
		repoUser: repoUser,
	}
	rg.Use(authing.Authorize())
	rg.PUT("/v1/user", pc.Update)
	rg.GET("/v1/user/checkLogin", pc.CheckLogin)
	rg.GET("/v1/user/callback", pc.AuthingCallback)
	rg.GET("/v1/user/getCurrentUser", pc.GetCurrentUser)

}

type UserController struct {
	repoUser repository.User
}

// @Summary Update
// @Description update user basic info
// @Tags  User
// @Accept json
// @Produce json
// @Router /v1/user [put]
func (uc *UserController) Update(ctx *gin.Context) {
	m := userBasicInfoModel{}

	if err := ctx.ShouldBindJSON(&m); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := uc.genUpdateUserBasicInfoCmd(&m)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(
			errorBadRequestParam, err,
		))

		return
	}

	s := app.NewUserService(uc.repoUser)

	if err := s.UpdateBasicInfo("", cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(
			errorSystemError, err,
		))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(m))
}

func (uc *UserController) genUpdateUserBasicInfoCmd(m *userBasicInfoModel) (
	cmd app.UpdateUserBasicInfoCmd,
	err error,
) {
	cmd.Bio, err = domain.NewBio(m.Bio)
	if err != nil {
		return
	}

	cmd.NickName, err = domain.NewNickname(m.Nickname)
	if err != nil {
		return
	}

	cmd.AvatarId, err = domain.NewAvatarId(m.AvatarId)

	return
}

// @Summary CheckLogin
// @Description CheckLogin
// @Tags  User
// @Accept json
// @Produce json
// @Router /v1/user/checkLogin [get]
func (uc *UserController) CheckLogin(c *gin.Context) {
	state, _ := utils.RandString(16)
	nonce, _ := utils.RandString(16)
	authing.SetCallbackCookie(c.Writer, c.Request, "state", state)
	authing.SetCallbackCookie(c.Writer, c.Request, "nonce", nonce)
	c.Redirect(http.StatusFound, authing.OIDCConfig.AuthCodeURL(state, oidc.Nonce(nonce)))
}

// @Summary AuthingCallback
// @Description login success callback
// @Tags  User
// @Accept json
// @Produce json
// @Router /v1/user/callback [get]
func (uc *UserController) AuthingCallback(c *gin.Context) {
	state, err := c.Request.Cookie("state")
	if err != nil {
		c.JSON(http.StatusInternalServerError, newResponse(
			"400", "Cookie Error", err.Error(),
		))
		return
	}
	if c.Request.URL.Query().Get("state") != state.Value {
		c.JSON(http.StatusInternalServerError,
			newResponse(
				"400", "state Error", err.Error(),
			))
		return
	}
	result, err := authing.GetTokenFromAuthing(c.Request.URL.Query().Get("code"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, newResponse(
			"400", "code Error", err.Error(),
		))
		return
	}
	c.JSON(http.StatusOK, newResponse(
		"200", "ok", result,
	))
}

// @Summary GetCurrentUser
// @Description 获取用户资料, 在请求的request的header中必须带有accessToken
// @Tags  User
// @Accept json
// @Produce json
// @Router /v1/user/getCurrentUser [get]
func (uc *UserController) GetCurrentUser(c *gin.Context) {
	currentUserClient := c.Keys["me"].(*authentication.Client)
	c.JSON(http.StatusOK,
		newResponse(
			"200", "ok", currentUserClient.ClientUser,
		))
}
