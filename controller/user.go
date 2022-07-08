package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForUserController(
	rg *gin.RouterGroup,
	repoUser repository.User,
) {
	pc := UserController{
		repoUser: repoUser,
	}

	rg.POST("/v1/user", pc.Update)
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
