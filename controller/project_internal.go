package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"

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

	rg.GET("/v1/space/:id", ctl.GetSpaceById)
	rg.PUT("/v1/space/:id/notify_update_code", ctl.NotifyUpdateCode)
	rg.GET("/v1/repo/:type/:user/:name/files", ctl.List)
	rg.GET("/v1/repo/:type/:user/:name/file/:path", ctl.Download)
}

type ProjectInternalController struct {
	baseController

	user userrepo.User
	repo spacerepo.Project
	s    spaceapp.ProjectService
	s1   app.RepoFileService

	model   repository.Model
	dataset repository.Dataset
	tags    repository.Tags
	like    repository.Like

	newPlatformRepository func(string, string) platform.Repository
}

func (ctl *ProjectInternalController) getRepoInfo(ctx *gin.Context, user domain.Account) (
	s resourceSummary, err error,
) {
	rt, err := domain.NewResourceType(ctx.Param("type"))
	if err != nil {
		return
	}

	name, err := domain.NewResourceName(ctx.Param("name"))
	if err != nil {
		return
	}

	s.rt = rt

	switch rt.ResourceType() {
	case domain.ResourceTypeModel.ResourceType():
		s.ResourceSummary, err = ctl.model.GetSummaryByName(user, name)

	case domain.ResourceTypeProject.ResourceType():
		s.ResourceSummary, err = ctl.repo.GetSummaryByName(user, name)

	case domain.ResourceTypeDataset.ResourceType():
		s.ResourceSummary, err = ctl.dataset.GetSummaryByName(user, name)
	}

	return
}

func (ctl *ProjectInternalController) checkForView(ctx *gin.Context) (
	pl *oldUserTokenPayload,
	u platform.UserInfo,
	repoInfo resourceSummary, ok bool,
) {
	user, err := domain.NewAccount(ctx.Param("user"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, visitor, b := ctl.checkUserApiToken(ctx, true)
	if !b {
		return
	}

	repoInfo, err = ctl.getRepoInfo(ctx, user)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	viewOther := visitor || pl.isNotMe(user)

	if viewOther && !repoInfo.IsPublic() {
		ctx.JSON(http.StatusNotFound, newResponseCodeMsg(
			errorResourceNotExists,
			"can't access other people's private/online project",
		))

		return
	}

	if viewOther {
		u.User = user
	} else {
		u = pl.PlatformUserInfo()
	}

	ok = true

	return
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

// todo: 修复此函数
func (ctl *ProjectInternalController) Download(ctx *gin.Context) {
	pl, u, repoInfo, ok := ctl.checkForView(ctx)
	if !ok {
		return
	}

	cmd := app.RepoFileDownloadCmd{
		Type:     repoInfo.rt,
		MyToken:  u.Token,
		Resource: repoInfo.ResourceSummary,
	}
	if pl.Account != "" {
		cmd.MyAccount = pl.DomainAccount()
	}

	var err error
	if cmd.Path, err = domain.NewFilePath(ctx.Param("path")); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if v, err := ctl.s1.Download(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}

// todo: 修复此函数
func (ctl *ProjectInternalController) List(ctx *gin.Context) {
	_, u, repoInfo, ok := ctl.checkForView(ctx)
	if !ok {
		return
	}

	var err error
	info := app.RepoDir{
		RepoName: repoInfo.Name,
	}

	info.Path, err = domain.NewDirectory(ctl.getQueryParameter(ctx, "path"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	v, err := ctl.s1.List(&u, &info)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(v))
}

// @Summary  NotifyUpdateCode space
// @Description  NotifyUpdateCode space
// @Tags     SpaceInternal
// @Param    id    path  string            true  "id of space" MaxLength(20)
// @Param    body  body  reqToNotifyUpdateCode  true  "body"
// @Accept   json
// @Security Internal
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
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
