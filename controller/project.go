package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/interfaces"
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
	p := interfaces.Project{}

	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, newResponse(
			errorBadRequestBody,
			"can't fetch request body",
			nil,
		))

		return
	}

	s := app.NewCreateProjectService(pc.repo)

	d, err := s.Create("", p.GenCreateProjectCmd())
	if err != nil {

	}

	c.JSON(http.StatusOK, newResponseData(d))
}
