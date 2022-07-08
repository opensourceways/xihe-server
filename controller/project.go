package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
)

func AddRouterForProjectController(rg *gin.RouterGroup, repo app.ProjectRepository) {
	pc := ProjectController{
		repo: repo,
	}

	rg.POST("/v1/project", pc.Create)
}

type ProjectController struct {
	repo app.ProjectRepository
}

// @Summary create project
// @Description create project
// @Tags  Project
// @Accept json
// @Produce json
// @Router /v1/project [post]
func (pc *ProjectController) Create(c *gin.Context) {
	p := projectModel{}

	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, newResponse(
			errorBadRequestBody,
			"can't fetch request body",
			nil,
		))

		return
	}

	s := app.NewCreateProjectService(pc.repo)

	cmd, err := pc.genCreateProjectCmd(&p)
	if err != nil {

	}

	d, err := s.Create("", cmd)
	if err != nil {

	}

	c.JSON(http.StatusOK, newResponseData(d))
}

func (pc *ProjectController) genCreateProjectCmd(p *projectModel) (cmd app.CreateProjectCmd, err error) {
	n, err := domain.NewProjName(p.Name)
	if err != nil {
		return
	}
	cmd.Name = n

	t, err := domain.NewRepoType(p.Type)
	if err != nil {
		return
	}
	cmd.Type = t

	d, err := domain.NewProjDesc(p.Desc)
	if err != nil {
		return
	}
	cmd.Desc = d

	// TODO: check cover id in db
	cmd.CoverId = p.CoverId

	pv, err := domain.NewProtocolName(p.Protocol)
	if err != nil {
		return
	}
	cmd.Protocol = pv

	tv, err := domain.NewTrainingSDK(p.Training)
	if err != nil {
		return
	}
	cmd.Training = tv

	iv, err := domain.NewInferenceSDK(p.Inference)
	if err != nil {
		return
	}
	cmd.Inference = iv

	return
}
