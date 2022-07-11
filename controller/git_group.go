package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/infrastructure/authing"
)

func AddRouterForGitGroupController(
	rg *gin.RouterGroup,
) {
	ctl := GitGroupController{}
	rg.Use(authing.Authorize())
	rg.POST("/v1/git/group", ctl.CreateGroup)
}

type GitGroupController struct {
}

// @Summary CreateGroup
// @Description Create Group of git
// @Tags  GitGroup
// @Accept json
// @Produce json
// @Router /v1/git/group [post]
func (uc *GitGroupController) CreateGroup(ctx *gin.Context) {
	m := userBasicInfoModel{}

	if err := ctx.ShouldBindJSON(&m); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(m))
}
