package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

func AddRouterForProjectController(rg *gin.RouterGroup, app app.ProjectApp) {
	pc := ProjectController{
		app: app,
	}

	rg.POST("/v1/project", pc.Create)
	rg.PUT("/v1/project/likeCount", pc.LikeCountIncrease)
}

type ProjectController struct {
	app app.ProjectApp
}

// @Summary create project
// @Description create project
// @Tags  Project
// @Accept json
// @Produce json
// @Router /v1/project [post]
func (pc *ProjectController) Create(ctx *gin.Context) {
	p := projectModel{}

	if err := ctx.ShouldBindJSON(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := pc.genCreateProjectCmd(&p)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(
			errorBadRequestParam, err,
		))
		return
	}

	s := app.NewCreateProjectService(nil)

	d, err := s.Create(cmd)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(
			errorSystemError, err,
		))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(d))
}

func (pc *ProjectController) genCreateProjectCmd(p *projectModel) (cmd app.CreateProjectCmd, err error) {
	cmd.Name, err = domain.NewProjName(p.Name)
	if err != nil {
		return
	}

	cmd.Type, err = domain.NewRepoType(p.Type)
	if err != nil {
		return
	}

	cmd.Desc, err = domain.NewProjDesc(p.Desc)
	if err != nil {
		return
	}

	cmd.CoverId, err = domain.NewConverId(p.CoverId)
	if err != nil {
		return
	}

	cmd.Protocol, err = domain.NewProtocolName(p.Protocol)
	if err != nil {
		return
	}

	cmd.Training, err = domain.NewTrainingSDK(p.Training)
	if err != nil {
		return
	}

	cmd.Inference, err = domain.NewInferenceSDK(p.Inference)

	return
}

// @Summary LikeCountIncrease
// @Description like count increase
// @Tags  Project
// @Param	project_id		query 	string	true		"id for project"
// @Param	user_id		query 	string	true		"id for user"
// @Accept json
// @Produce json
// @Router /v1/project/likeCount [put]
func (pc *ProjectController) LikeCountIncrease(c *gin.Context) {
	project_id := c.Query("project_id")
	user_id := c.Query("user_id")
	data, err := pc.app.LikeCountIncrease(nil, project_id, user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, newResponseError(
			errorSystemError, err,
		))
		return
	}
	c.JSON(http.StatusOK, newResponse("200", "ok", data))
}
