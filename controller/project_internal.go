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
	spacedomain "github.com/opensourceways/xihe-server/space/domain"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
)

func AddRouterForProjectInternalController(
	rg *gin.RouterGroup,
	user userrepo.User,

	model repository.Model,
	dataset repository.Dataset,
	activity repository.Activity,
	tags repository.Tags,
	like repository.Like,
	sender message.ResourceProducer,
	newPlatformRepository func(token, namespace string) platform.Repository,
	computility computilityapp.ComputilityInternalAppService,
	spaceProducer spacedomain.SpaceEventProducer,
	repoPg spacerepo.ProjectPg,
) {
	ctl := ProjectInternalController{
		user:    user,
		repoPg:  repoPg,
		model:   model,
		dataset: dataset,
		tags:    tags,
		like:    like,
		s: spaceapp.NewProjectService(
			user, repoPg, model, dataset, activity, sender, computility, spaceProducer,
		),
		newPlatformRepository: newPlatformRepository,
	}

	rg.GET("/v1/space/:id", internalApiCheckMiddleware(&ctl.baseController), ctl.GetSpaceById)
	rg.PUT("/v1/space/:id/notify_update_code", internalApiCheckMiddleware(&ctl.baseController), ctl.NotifyUpdateCode)
}

type ProjectInternalController struct {
	baseController

	user   userrepo.User
	repoPg spacerepo.ProjectPg
	s      spaceapp.ProjectService

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
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}
	proj, err := ctl.s.GetByRepoId(id)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(proj))
}

// @Summary  NotifyUpdateCode space
// @Description  NotifyUpdateCode space
// @Tags     SpaceInternal
// @Param    id    path  string            true  "id of space" MaxLength(20)
// @Param    body  body  reqToNotifyUpdateCode  true  "body"
// @Accept   json
// @Security Internal
// @Success  202   {object}   responseData{data=string,msg=string,code=string}
// @Router   /v1/space/{id}/notify_update_code [put]
func (ctl *ProjectInternalController) NotifyUpdateCode(ctx *gin.Context) {
	req := reqToNotifyUpdateCode{}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	spaceId, err := domain.NewIdentity(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if err = ctl.s.NotifyUpdateCodes(spaceId, &cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(""))
	}
}
