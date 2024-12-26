package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	filescan "github.com/opensourceways/xihe-server/filescan/app"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	uapp "github.com/opensourceways/xihe-server/user/app"
)

func AddRouterForRepoFileInternalController(
	rg *gin.RouterGroup,
	p platform.RepoFile,
	model repository.Model,
	project spacerepo.Project,
	dataset repository.Dataset,
	sender message.RepoMessageProducer,
	us uapp.UserService,
	f filescan.FileScanService,
) {
	ctl := RepoFileInternalController{
		s:       app.NewRepoFileService(p, sender, f),
		us:      us,
		model:   model,
		project: project,
		dataset: dataset,
	}

	rg.GET("/v1/repo/:type/:user/:name/files", internalApiCheckMiddleware(&ctl.baseController), ctl.List)
	rg.GET("/v1/repo/:type/:user/:name/file/:path", internalApiCheckMiddleware(&ctl.baseController), ctl.Download)
	rg.GET("/v1/resource/repo/:id/:path", internalApiCheckMiddleware(&ctl.baseController), ctl.DownloadById)
}

type RepoFileInternalController struct {
	baseController

	s       app.RepoFileService
	us      uapp.UserService
	model   repository.Model
	project spacerepo.Project
	dataset repository.Dataset
}

// @Summary		Download
// @Description	Download repo file
// @Tags			RepoFileInternal
// @Param			user			path	string	true	"user"
// @Param			name			path	string	true	"repo name"
// @Param			path			path	string	true	"repo file path"
// @Accept			json
// @Success		200	{object}			app.RepoFileDownloadDTO
// @Failure		400	bad_request_param	some	parameter	of	body	is	invalid
// @Failure		500	system_error		system	error
// @Router			/v1/repo/{type}/{user}/{name}/file/{path} [get]
func (ctl *RepoFileInternalController) Download(ctx *gin.Context) {
	u, repoInfo, ok := ctl.checkForView(ctx)
	if !ok {
		return
	}

	cmd := app.RepoFileDownloadCmd{
		Type:     repoInfo.rt,
		MyToken:  u.Token,
		Resource: repoInfo.ResourceSummary,
	}

	var err error
	if cmd.Path, err = domain.NewFilePath(ctx.Param("path")); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if v, err := ctl.s.Download(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}

// @Summary		List
// @Description	list repo file in a path
// @Tags			RepoFileInternal
// @Param			user	path	string	true	"user"
// @Param			name	path	string	true	"repo name"
// @Param			path	query	string	true	"repo file path"
// @Accept			json
// @Success		200	{object}			app.RepoPathItem
// @Failure		400	bad_request_param	some	parameter	of	body	is	invalid
// @Failure		500	system_error		system	error
// @Router			/v1/repo/{type}/{user}/{name}/files [get]
func (ctl *RepoFileInternalController) List(ctx *gin.Context) {
	u, repoInfo, ok := ctl.checkForView(ctx)
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

	v, err := ctl.s.List(&u, &info)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(v))
}

func (ctl *RepoFileInternalController) checkForView(ctx *gin.Context) (
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

	repoInfo, err = ctl.getRepoInfo(ctx, user)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	u.User = user
	ok = true

	return
}

func (ctl *RepoFileInternalController) getRepoInfo(ctx *gin.Context, user domain.Account) (
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
		s.ResourceSummary, err = ctl.project.GetSummaryByName(user, name)

	case domain.ResourceTypeDataset.ResourceType():
		s.ResourceSummary, err = ctl.dataset.GetSummaryByName(user, name)
	}

	return
}

// @Summary		DownloadById
// @Description	Download repo file by Id
// @Tags			RepoFileInternal
// @Param			id			    path	string	true	"repository id"
// @Param			path			path	string	true	"repo file path"
// @Accept			json
// @Success		200	{object}			app.RepoFileDownloadDTO
// @Failure		400	bad_request_param	some	parameter	of	body	is	invalid
// @Failure		500	system_error		system	error
// @Router			/v1/repo/{id}/{path} [get]
func (ctl *RepoFileInternalController) DownloadById(ctx *gin.Context) {
	repoId := ctx.Param("id")

	filePath, err := domain.NewFilePath(ctx.Param("path"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if v, err := ctl.s.DownloadByRepoId(repoId, filePath); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}
