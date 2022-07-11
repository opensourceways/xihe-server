package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForProjectController(rg *gin.RouterGroup, repo repository.Project) {
	pc := ProjectController{
		repo: repo,
	}

	rg.POST("/v1/project", pc.Create)
}

type ProjectController struct {
	repo repository.Project
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

	s := app.NewCreateProjectService(pc.repo)

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
