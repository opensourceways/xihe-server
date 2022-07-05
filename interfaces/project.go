package interfaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/application"
	"github.com/opensourceways/xihe-server/domain/entity"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/util"
)

//Users struct defines the dependencies that will be used
type Project struct {
	app *application.ProjectApp
}

//Users constructor
func NewProject(repo repository.ProjectRepository) *Project {
	app := application.NewProjectAPP(repo)
	return &Project{
		app: app,
	}
}

// @Summary Save
// @Description Save
// @Tags  Project
// @Param	body		body 	entity.Project	true		"email username phone"
// @Accept json
// @Produce json
// @Router /v1/project/save [post]
func (p *Project) Save(c *gin.Context) {
	var item entity.Project
	err := c.ShouldBindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.ExportData(http.StatusBadRequest, "error input data ", err.Error()))
		return
	}
	result, err := p.app.Save(&item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.ExportData(http.StatusInternalServerError, "save error ", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.ExportData(util.CodeStatusNormal, "ok", result))
}
