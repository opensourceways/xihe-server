package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"

	computilityapp "github.com/opensourceways/xihe-server/computility/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	spaceapp "github.com/opensourceways/xihe-server/space/app"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
)

func AddRouterForProjectInternalController(
	rg *gin.RouterGroup,
	user userrepo.User,
	repo spacerepo.Project,
	model repository.Model,
	dataset repository.Dataset,
	activity repository.Activity,
	tags repository.Tags,
	like repository.Like,
	sender message.ResourceProducer,
	newPlatformRepository func(token, namespace string) platform.Repository,
	computility computilityapp.ComputilityInternalAppService,
) {
	ctl := ProjectInternalController{
		user:    user,
		repo:    repo,
		model:   model,
		dataset: dataset,
		tags:    tags,
		like:    like,
		s: spaceapp.NewProjectService(
			user, repo, model, dataset, activity, nil, sender, computility,
		),
		newPlatformRepository: newPlatformRepository,
	}

	rg.POST("/v1/space/:id", ctl.GetSpaceById)
}

type ProjectInternalController struct {
	baseController

	user userrepo.User
	repo spacerepo.Project
	s    spaceapp.ProjectService

	model   repository.Model
	dataset repository.Dataset
	tags    repository.Tags
	like    repository.Like

	newPlatformRepository func(string, string) platform.Repository
}

// @Summary		Check
// @Description	check whether the name can be applied to create a new project
// @Tags			Project
// @Param			owner	path	string	true	"owner of project"
// @Param			name	path	string	true	"name of project"
// @Accept			json
// @Success		200	{object}	canApplyResourceNameResp
// @Produce		json
// @Router			/v1/project/{owner}/{name}/check [get]
func (ctl *ProjectInternalController) GetSpaceById(ctx *gin.Context) {
	id, err := domain.NewIdentity(ctx.Param("id"))
	fmt.Printf("id: %v\n", id)
	fmt.Printf("err: %v\n", err)
}
