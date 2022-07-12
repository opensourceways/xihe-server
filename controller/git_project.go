package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/infrastructure/authing"
)

func AddRouterForGitProjectController(
	rg *gin.RouterGroup,
) {
	ctl := GitProjectController{}
	rg.Use(authing.Authorize())
	rg.POST("/v1/git/project", ctl.CreateProject)
}

type GitProjectController struct {
}

// @Summary CreateProject
// @Description Create Project for gitlab
// @Tags  GitProject
// @Param	body		body 	domain.GitlabProject	true		"body for CreateProject content"
// @Accept json
// @Produce json
// @Router /v1/git/project [post]
func (gc *GitProjectController) CreateProject(ctx *gin.Context) {
	var m domain.GitlabProject
	if err := ctx.ShouldBindJSON(&m); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))
		return
	}

	ctx.JSON(http.StatusOK, newResponseData(m))
}
