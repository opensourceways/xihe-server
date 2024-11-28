package controller

import (
	"fmt"
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
	"github.com/opensourceways/xihe-server/utils"
)

func AddRouterForProjectController(
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
	spaceProducer spacedomain.SpaceEventProducer,
	repoPg spacerepo.ProjectPg,
) {
	ctl := ProjectController{
		user:    user,
		repo:    repo,
		repoPg:  repoPg,
		model:   model,
		dataset: dataset,
		tags:    tags,
		like:    like,
		s: spaceapp.NewProjectService(
			user, repo, repoPg, model, dataset, activity, sender, computility, spaceProducer,
		),
		newPlatformRepository: newPlatformRepository,
	}

	rg.POST("/v1/project", checkUserEmailMiddleware(&ctl.baseController), ctl.Create)
	rg.PUT("/v1/project/:owner/:id", checkUserEmailMiddleware(&ctl.baseController), ctl.Update)
	rg.DELETE("/v1/project/:owner/:name", checkUserEmailMiddleware(&ctl.baseController), ctl.Delete)
	rg.GET("/v1/project/:owner/:name", ctl.Get)
	rg.GET("/v1/project/:owner/:name/check", ctl.Check)
	rg.GET("/v1/project/:owner", ctl.List)
	rg.GET("/v1/project", ctl.ListGlobal)

	rg.POST("/v1/project/:owner/:id", checkUserEmailMiddleware(&ctl.baseController), ctl.Fork)

	rg.PUT("/v1/project/relation/:owner/:id/model",
		checkUserEmailMiddleware(&ctl.baseController), ctl.AddRelatedModel)
	rg.DELETE("/v1/project/relation/:owner/:id/model",
		checkUserEmailMiddleware(&ctl.baseController), ctl.RemoveRelatedModel)

	rg.PUT("/v1/project/relation/:owner/:id/dataset",
		checkUserEmailMiddleware(&ctl.baseController), ctl.AddRelatedDataset)
	rg.DELETE("/v1/project/relation/:owner/:id/dataset",
		checkUserEmailMiddleware(&ctl.baseController), ctl.RemoveRelatedDataset)

	rg.PUT("/v1/project/:owner/:id/tags", checkUserEmailMiddleware(&ctl.baseController), ctl.SetTags)
}

type ProjectController struct {
	baseController

	user   userrepo.User
	repo   spacerepo.Project
	repoPg spacerepo.ProjectPg

	s spaceapp.ProjectService

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
func (ctl *ProjectController) Check(ctx *gin.Context) {
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

// @Summary		Create
// @Description	create project
// @Tags			Project
// @Param			body	body	projectCreateRequest	true	"body of creating project"
// @Accept			json
// @Produce		json
// @Router			/v1/project [post]
func (ctl *ProjectController) Create(ctx *gin.Context) {
	req := projectCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	tags, err := ctl.tags.List(apiConfig.Tags.ProjectTagDomains)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	cmd, err := req.toCmd(tags)
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

	if pl.isNotMe(cmd.Owner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed,
			"can't create project for other user",
		))

		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "create project")

	pr := ctl.newPlatformRepository(
		pl.PlatformToken, pl.PlatformUserNamespaceId,
	)

	d, err := ctl.s.Create(&cmd, pr)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	utils.DoLog("", pl.DomainAccount().Account(), "create project",
		fmt.Sprintf("%s/%s", d.Owner, d.Name), "success")

	ctx.JSON(http.StatusCreated, newResponseData(d))
}

// @Summary		Delete
// @Description	delete project
// @Tags			Project
// @Param			owner	path	string	true	"owner of project"
// @Param			name	path	string	true	"name of project"
// @Accept			json
// @Success		204
// @Produce		json
// @Router			/v1/project/{owner}/{name} [delete]
func (ctl *ProjectController) Delete(ctx *gin.Context) {
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
		ctx.JSON(http.StatusNotFound, newResponseCodeMsg(
			errorResourceNotExists,
			"can't access other's project",
		))

		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "deleta project")

	proj, err := ctl.repoPg.GetByName(owner, name)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	pr := ctl.newPlatformRepository(
		pl.PlatformToken, pl.PlatformUserNamespaceId,
	)

	if err := ctl.s.Delete(&proj, pr); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		utils.DoLog("", pl.DomainAccount().Account(), "delete project",
			fmt.Sprintf("%s/%s", proj.Owner, proj.Name), "success")

		ctx.JSON(http.StatusNoContent, newResponseData("success"))
	}
}

// @Summary		Update
// @Description	update project
// @Tags			Project
// @Param			owner	path	string					true	"owner of project"
// @Param			id		path	string					true	"id of project"
// @Param			body	body	projectUpdateRequest	true	"body of updating project"
// @Accept			json
// @Produce		json
// @Router			/v1/project/{owner}/{id} [put]
func (ctl *ProjectController) Update(ctx *gin.Context) {
	req := projectUpdateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
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

	owner, err := domain.NewAccount(ctx.Param("owner"))
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
			errorNotAllowed,
			"can't update project for other user",
		))

		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "update project")

	proj, err := ctl.repoPg.Get(owner, ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(err))

		return
	}

	pr := ctl.newPlatformRepository(
		pl.PlatformToken, pl.PlatformUserNamespaceId,
	)

	d, err := ctl.s.Update(&proj, &cmd, pr)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(d))
}

// @Summary		Get
// @Description	get project
// @Tags			Project
// @Param			owner	path	string	true	"owner of project"
// @Param			name	path	string	true	"name of project"
// @Accept			json
// @Success		200	{object}	projectDetail
// @Produce		json
// @Router			/v1/project/{owner}/{name} [get]
func (ctl *ProjectController) Get(ctx *gin.Context) {
	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	avatar, err := ctl.user.GetUserAvatarId(owner)
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

	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	proj, err := ctl.s.GetByName(owner, name, !visitor && pl.isMyself(owner))
	if err != nil {
		if isErrorOfAccessingPrivateRepo(err) {
			ctx.JSON(http.StatusNotFound, newResponseCodeMsg(
				errorResourceNotExists,
				"can't access private project",
			))
		} else {
			ctl.sendRespWithInternalError(ctx, newResponseError(err))
		}

		return
	}

	liked := true
	if !visitor && pl.isNotMe(owner) {
		obj := &domain.ResourceObject{Type: domain.ResourceTypeProject}
		obj.Owner = owner
		obj.Id = proj.Id

		liked, err = ctl.like.HasLike(pl.DomainAccount(), obj)

		if err != nil {
			ctl.sendRespWithInternalError(ctx, newResponseError(err))

			return
		}
	}

	detail := projectDetail{
		Liked:            liked,
		ProjectDetailDTO: &proj,
	}
	if avatar != nil {
		detail.AvatarId = avatar.AvatarId()
	}

	ctx.JSON(http.StatusOK, newResponseData(detail))
}

// @Summary		List
// @Description	list project
// @Tags			Project
// @Param			owner			path	string	true	"owner of project"
// @Param			name			query	string	false	"name of project"
// @Param			repo_type		query	string	false	"repo type of project, value can be public or private"
// @Param			count_per_page	query	int		false	"count per page"
// @Param			page_num		query	int		false	"page num which starts from 1"
// @Param			sort_by			query	string	false	"sort keys, value can be update_time, first_letter, download_count"
// @Accept			json
// @Success		200	{object}	projectsInfo
// @Produce		json
// @Router			/v1/project/{owner} [get]
func (ctl *ProjectController) List(ctx *gin.Context) {
	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	cmd, err := ctl.getListResourceParameter(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if visitor || pl.isNotMe(owner) {
		if cmd.RepoType == nil {
			type1, _ := domain.NewRepoType(domain.RepoTypePublic)
			type2, _ := domain.NewRepoType(domain.RepoTypeOnline)
			cmd.RepoType = append(cmd.RepoType, type1)
			cmd.RepoType = append(cmd.RepoType, type2)
		} else {
			if cmd.RepoType[0].RepoType() != domain.RepoTypePublic {
				ctx.JSON(http.StatusOK, newResponseData(nil))
				return
			}
		}
	}

	projs, err := ctl.s.List(owner, &cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	avatar, err := ctl.user.GetUserAvatarId(owner)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	result := projectsInfo{
		Owner:       owner.Account(),
		ProjectsDTO: &projs,
	}
	if avatar != nil {
		result.AvatarId = avatar.AvatarId()
	}

	ctx.JSON(http.StatusOK, newResponseData(&result))
}

// @Summary		ListGlobal
// @Description	list global public project
// @Tags			Project
// @Param			name			query	string	false	"name of project"
// @Param			tags			query	string	false	"tags, separate multiple tags with commas"
// @Param			tag_kinds		query	string	false	"tag kinds, separate multiple kinds with commas"
// @Param			level			query	string	false	"project level, such as official, good"
// @Param			count_per_page	query	int		false	"count per page"
// @Param			page_num		query	int		false	"page num which starts from 1"
// @Param			sort_by			query	string	false	"sort keys, value can be update_time, first_letter, download_count"
// @Accept			json
// @Success		200	{object}	app.GlobalProjectsDTO
// @Produce		json
// @Router			/v1/project [get]
func (ctl *ProjectController) ListGlobal(ctx *gin.Context) {
	cmd, err := ctl.getListGlobalResourceParameter(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	result, err := ctl.s.ListGlobal(&cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(result))
}

// @Summary		Fork
// @Description	fork project
// @Tags			Project
// @Param			owner	path	string				true	"owner of forked project"
// @Param			id		path	string				true	"id of forked project"
// @Param			body	body	projectForkRequest	true	"body of forking project"
// @Accept			json
// @Produce		json
// @Router			/v1/project/{owner}/{id} [post]
func (ctl *ProjectController) Fork(ctx *gin.Context) {
	req := projectForkRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
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

	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	tags, err := ctl.tags.List(apiConfig.Tags.ProjectTagDomains)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if !pl.isNotMe(owner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed, "no need to fork project of yourself",
		))

		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "fork project")

	proj, err := ctl.repoPg.Get(owner, ctx.Param("id"))
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	// TODO maybe the private project can be forked by special user.

	if proj.IsPrivate() {
		ctx.JSON(http.StatusNotFound, newResponseCodeMsg(
			errorResourceNotExists,
			"can't access private project",
		))

		return
	}

	if proj.IsOnline() {
		ctx.JSON(http.StatusNotFound, newResponseCodeMsg(
			errorResourceNotExists,
			"can't fork online project",
		))

		return
	}

	pr := ctl.newPlatformRepository(
		pl.PlatformToken, pl.PlatformUserNamespaceId,
	)

	cmd.From = proj
	cmd.Owner = pl.DomainAccount()
	cmd.ValidTags = tags

	data, err := ctl.s.Fork(&cmd, pr)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}
	ctx.JSON(http.StatusCreated, newResponseData(data))
}

// @Summary		AddRelatedModel
// @Description	add related model to project
// @Tags			Project
// @Param			owner	path	string						true	"owner of project"
// @Param			id		path	string						true	"id of project"
// @Param			body	body	relatedResourceAddRequest	true	"body of related model"
// @Accept			json
// @Success		202	{object}	app.ResourceDTO
// @Router			/v1/project/relation/{owner}/{id}/model [put]
func (ctl *ProjectController) AddRelatedModel(ctx *gin.Context) {
	req := relatedResourceAddRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	owner, name, err := req.toModelCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	model, err := ctl.model.GetByName(owner, name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, proj, ok := ctl.checkPermissionPg(ctx)
	if !ok {
		return
	}

	if pl.isNotMe(owner) && model.IsPrivate() {
		ctx.JSON(http.StatusNotFound, newResponseCodeMsg(
			errorResourceNotExists,
			"can't access private project",
		))

		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "add related model to project")

	index := domain.ResourceIndex{
		Owner: owner,
		Id:    model.Id,
	}

	if err = ctl.s.AddRelatedModel(&proj, &index); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(convertToRelatedResource(model)))
}

// @Summary		RemoveRelatedModel
// @Description	remove related model to project
// @Tags			Project
// @Param			owner	path	string							true	"owner of project"
// @Param			id		path	string							true	"id of project"
// @Param			body	body	relatedResourceRemoveRequest	true	"body of related model"
// @Accept			json
// @Success		204
// @Router			/v1/project/relation/{owner}/{id}/model [delete]
func (ctl *ProjectController) RemoveRelatedModel(ctx *gin.Context) {
	req := relatedResourceRemoveRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
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

	pl, proj, ok := ctl.checkPermission(ctx)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "remove related model to project")

	index := domain.ResourceIndex{
		Owner: cmd.Owner,
		Id:    cmd.Id,
	}
	if err = ctl.s.RemoveRelatedModel(&proj, &index); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusNoContent, newResponseData("success"))
}

// @Summary		AddRelatedDataset
// @Description	add related dataset to project
// @Tags			Project
// @Param			owner	path	string						true	"owner of project"
// @Param			id		path	string						true	"id of project"
// @Param			body	body	relatedResourceAddRequest	true	"body of related dataset"
// @Accept			json
// @Success		202	{object}	app.ResourceDTO
// @Router			/v1/project/relation/{owner}/{id}/dataset [put]
func (ctl *ProjectController) AddRelatedDataset(ctx *gin.Context) {
	req := relatedResourceAddRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	owner, name, err := req.toDatasetCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	data, err := ctl.dataset.GetByName(owner, name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, proj, ok := ctl.checkPermissionPg(ctx)
	if !ok {
		return
	}

	if pl.isNotMe(owner) && data.IsPrivate() {
		ctx.JSON(http.StatusNotFound, newResponseCodeMsg(
			errorResourceNotExists,
			"can't access private project",
		))

		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "add related dataset to project")

	index := domain.ResourceIndex{
		Owner: owner,
		Id:    data.Id,
	}
	if err = ctl.s.AddRelatedDataset(&proj, &index); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(convertToRelatedResource(data)))
}

// @Summary		RemoveRelatedDataset
// @Description	remove related dataset to project
// @Tags			Project
// @Param			owner	path	string							true	"owner of project"
// @Param			id		path	string							true	"id of project"
// @Param			body	body	relatedResourceRemoveRequest	true	"body of related dataset"
// @Accept			json
// @Success		204
// @Router			/v1/project/relation/{owner}/{id}/dataset [delete]
func (ctl *ProjectController) RemoveRelatedDataset(ctx *gin.Context) {
	req := relatedResourceRemoveRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
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

	pl, proj, ok := ctl.checkPermission(ctx)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "remove related dataset to project")

	index := domain.ResourceIndex{
		Owner: cmd.Owner,
		Id:    cmd.Id,
	}
	if err = ctl.s.RemoveRelatedDataset(&proj, &index); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusNoContent, newResponseData("success"))
}

// @Summary		SetTags
// @Description	set tags for project
// @Tags			Project
// @Param			owner	path	string						true	"owner of project"
// @Param			id		path	string						true	"id of project"
// @Param			body	body	resourceTagsUpdateRequest	true	"body of tags"
// @Accept			json
// @Success		202
// @Router			/v1/project/{owner}/{id}/tags [put]
func (ctl *ProjectController) SetTags(ctx *gin.Context) {
	req := resourceTagsUpdateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	tags, err := ctl.tags.List(apiConfig.Tags.ProjectTagDomains)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	cmd, err := req.toCmd(tags)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, proj, ok := ctl.checkPermission(ctx)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "set tags for project")

	if err = ctl.s.SetTags(&proj, &cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData("success"))
}

func (ctl *ProjectController) checkPermission(ctx *gin.Context) (
	info *oldUserTokenPayload, proj spacedomain.Project, ok bool,
) {
	owner, err := domain.NewAccount(ctx.Param("owner"))
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

	proj, err = ctl.repoPg.Get(owner, ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(err))

		return
	}

	info = pl
	ok = true

	return
}

func (ctl *ProjectController) checkPermissionPg(ctx *gin.Context) (
	info *oldUserTokenPayload, proj spacedomain.Project, ok bool,
) {
	owner, err := domain.NewAccount(ctx.Param("owner"))
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

	proj, err = ctl.repoPg.Get(owner, ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(err))

		return
	}

	info = pl
	ok = true

	return
}
