package controller

import (
	"net/http"

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

	rg.POST("/v1/project/create", checkUserEmailMiddleware(&ctl.baseController), ctl.GetSpaceById)
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
	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	name, err := domain.NewResourceName(ctx.Param("name"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if pl.isNotMe(owner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed, "not allowed",
		))

		return
	}

	b := ctl.s.CanApplyResourceName(owner, name)

	ctx.JSON(http.StatusOK, newResponseData(canApplyResourceNameResp{b}))
}
