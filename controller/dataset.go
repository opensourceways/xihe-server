package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	auditcommon "github.com/opensourceways/xihe-server/common/domain/audit"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

func AddRouterForDatasetController(
	rg *gin.RouterGroup,
	user userrepo.User,
	repo repository.Dataset,
	model repository.Model,
	proj spacerepo.Project,
	activity repository.Activity,
	tags repository.Tags,
	like repository.Like,
	sender message.ResourceProducer,
	newPlatformRepository func(token, namespace string) platform.Repository,
	audit auditcommon.AuditService,
) {
	ctl := DatasetController{
		user: user,
		repo: repo,
		tags: tags,
		like: like,
		s:    app.NewDatasetService(user, repo, proj, model, activity, nil, sender, audit),

		newPlatformRepository: newPlatformRepository,
	}

	rg.POST("/v1/dataset", checkUserEmailMiddleware(&ctl.baseController), ctl.Create)
	rg.PUT("/v1/dataset/:owner/:id", checkUserEmailMiddleware(&ctl.baseController), ctl.Update)
	rg.DELETE("/v1/dataset/:owner/:name", checkUserEmailMiddleware(&ctl.baseController), ctl.Delete)
	rg.GET("/v1/dataset/:owner/:name/check", ctl.Check)
	rg.GET("/v1/dataset/:owner/:name", ctl.Get)
	rg.GET("/v1/dataset/:owner", ctl.List)
	rg.GET("/v1/dataset", ctl.ListGlobal)

	rg.PUT("/v1/dataset/:owner/:id/tags", checkUserEmailMiddleware(&ctl.baseController), ctl.SetTags)
}

type DatasetController struct {
	baseController

	user userrepo.User
	repo repository.Dataset
	tags repository.Tags
	like repository.Like
	s    app.DatasetService

	newPlatformRepository func(string, string) platform.Repository
}

// @Summary		Check
// @Description	check whether the name can be applied to create a new dataset
// @Tags			Dataset
// @Param			owner	path	string	true	"owner of dataset"
// @Param			name	path	string	true	"name of dataset"
// @Accept			json
// @Success		200	{object}	canApplyResourceNameResp
// @Produce		json
// @Router			/v1/dataset/{owner}/{name}/check [get]
func (ctl *DatasetController) Check(ctx *gin.Context) {
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
// @Description	create dataset
// @Tags			Dataset
// @Param			body	body	datasetCreateRequest	true	"body of creating dataset"
// @Accept			json
// @Success		201	{object}			app.DatasetDTO
// @Failure		400	bad_request_body	can't	parse		request	body
// @Failure		400	bad_request_param	some	parameter	of		body	is	invalid
// @Failure		500	system_error		system	error
// @Failure		500	duplicate_creating	create	dataset	repeatedly
// @Router			/v1/dataset [post]
func (ctl *DatasetController) Create(ctx *gin.Context) {
	req := datasetCreateRequest{}
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

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "creat dataset")

	if pl.isNotMe(cmd.Owner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed,
			"can't create dataset for other user",
		))

		return
	}

	pr := ctl.newPlatformRepository(
		pl.PlatformToken, pl.PlatformUserNamespaceId,
	)

	d, err := ctl.s.Create(&cmd, pr)
	if err != nil {
		SendError(ctx, err)

		return
	}

	utils.DoLog("", pl.DomainAccount().Account(), "create dataset",
		fmt.Sprintf("%s/%s", d.Owner, d.Name), "success")

	ctx.JSON(http.StatusOK, newResponseData(d))
}

// @Summary		Delete
// @Description	delete dataset
// @Tags			Dataset
// @Param			owner	path	string	true	"owner of dataset"
// @Param			name	path	string	true	"name of dataset"
// @Accept			json
// @Success		204
// @Produce		json
// @Router			/v1/dataset/{owner}/{name} [delete]
func (ctl *DatasetController) Delete(ctx *gin.Context) {
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

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "delete dataset")

	if pl.isNotMe(owner) {
		ctx.JSON(http.StatusNotFound, newResponseCodeMsg(
			errorResourceNotExists,
			"can't access other's dataset",
		))

		return
	}

	d, err := ctl.repo.GetByName(owner, name)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	pr := ctl.newPlatformRepository(
		pl.PlatformToken, pl.PlatformUserNamespaceId,
	)

	if err := ctl.s.Delete(&d, pr); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {

		utils.DoLog("", pl.DomainAccount().Account(), "delete dataset",
			fmt.Sprintf("%s/%s", d.Owner, d.Name), "success")

		ctx.JSON(http.StatusNoContent, newResponseData("success"))
	}
}

// @Summary		Update
// @Description	update property of dataset
// @Tags			Dataset
// @Param			id		path	string					true	"id of dataset"
// @Param			body	body	datasetUpdateRequest	true	"body of updating dataset"
// @Accept			json
// @Produce		json
// @Router			/v1/dataset/{owner}/{id} [put]
func (ctl *DatasetController) Update(ctx *gin.Context) {
	req := datasetUpdateRequest{}

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

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "update property of dataset")

	if pl.isNotMe(owner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed,
			"can't update dataset for other user",
		))

		return
	}

	m, err := ctl.repo.Get(owner, ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(err))

		return
	}

	pr := ctl.newPlatformRepository(
		pl.PlatformToken, pl.PlatformUserNamespaceId,
	)

	d, err := ctl.s.Update(&m, &cmd, pr)
	if err != nil {
		SendError(ctx, err)

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(d))
}

// @Summary		Get
// @Description	get dataset
// @Tags			Dataset
// @Param			owner	path	string	true	"owner of dataset"
// @Param			name	path	string	true	"name of dataset"
// @Accept			json
// @Success		200	{object}	datasetDetail
// @Produce		json
// @Router			/v1/dataset/{owner}/{name} [get]
func (ctl *DatasetController) Get(ctx *gin.Context) {
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

	d, err := ctl.s.GetByName(owner, name, !visitor && pl.isMyself(owner))
	if err != nil {
		if isErrorOfAccessingPrivateRepo(err) {
			ctx.JSON(http.StatusNotFound, newResponseCodeMsg(
				errorResourceNotExists,
				"can't access private dataset",
			))
		} else {
			ctl.sendRespWithInternalError(ctx, newResponseError(err))
		}

		return
	}

	liked := true
	if !visitor && pl.isNotMe(owner) {
		obj := &domain.ResourceObject{Type: domain.ResourceTypeDataset}
		obj.Owner = owner
		obj.Id = d.Id

		liked, err = ctl.like.HasLike(pl.DomainAccount(), obj)

		if err != nil {
			ctl.sendRespWithInternalError(ctx, newResponseError(err))

			return
		}
	}

	detail := datasetDetail{
		Liked:            liked,
		DatasetDetailDTO: &d,
	}
	if avatar != nil {
		detail.AvatarId = avatar.AvatarId()
	}

	ctx.JSON(http.StatusOK, newResponseData(detail))
}

// @Summary		List
// @Description	list dataset
// @Tags			Dataset
// @Param			owner			path	string	true	"owner of dataset"
// @Param			name			query	string	false	"name of dataset"
// @Param			repo_type		query	string	false	"repo type of dataset, value can be public or private"
// @Param			count_per_page	query	int		false	"count per page"
// @Param			page_num		query	int		false	"page num which starts from 1"
// @Param			sort_by			query	string	false	"sort keys, value can be update_time, first_letter, download_count"
// @Accept			json
// @Success		200	{object}	datasetsInfo
// @Produce		json
// @Router			/v1/dataset/{owner} [get]
func (ctl *DatasetController) List(ctx *gin.Context) {
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

	data, err := ctl.s.List(owner, &cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	avatar, err := ctl.user.GetUserAvatarId(owner)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	result := datasetsInfo{
		Owner:       owner.Account(),
		DatasetsDTO: &data,
	}
	if avatar != nil {
		result.AvatarId = avatar.AvatarId()
	}

	ctx.JSON(http.StatusOK, newResponseData(&result))
}

// @Summary		ListGlobal
// @Description	list global public dataset
// @Tags			Dataset
// @Param			name			query	string	false	"name of dataset"
// @Param			tags			query	string	false	"tags, separate multiple tags with commas"
// @Param			tag_kinds		query	string	false	"tag kinds, separate multiple kinds with commas"
// @Param			level			query	string	false	"dataset level, such as official, good"
// @Param			count_per_page	query	int		false	"count per page"
// @Param			page_num		query	int		false	"page num which starts from 1"
// @Param			sort_by			query	string	false	"sort keys, value can be update_time, first_letter, download_count"
// @Accept			json
// @Success		200	{object}	app.GlobalDatasetsDTO
// @Produce		json
// @Router			/v1/dataset [get]
func (ctl *DatasetController) ListGlobal(ctx *gin.Context) {
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

// @Summary		SetTags
// @Description	set tags for dataset
// @Tags			Dataset
// @Param			owner	path	string						true	"owner of dataset"
// @Param			id		path	string						true	"id of dataset"
// @Param			body	body	resourceTagsUpdateRequest	true	"body of tags"
// @Accept			json
// @Success		202
// @Router			/v1/dataset/{owner}/{id}/tags [put]
func (ctl *DatasetController) SetTags(ctx *gin.Context) {
	req := resourceTagsUpdateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	tags, err := ctl.tags.List(apiConfig.Tags.DatasetTagDomains)
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

	d, ok := ctl.checkPermission(ctx)
	if !ok {
		return
	}

	pl, _, _ := ctl.checkUserApiToken(ctx, false)
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "set tags for dataset")

	if err = ctl.s.SetTags(&d, &cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData("success"))
}

func (ctl *DatasetController) checkPermission(ctx *gin.Context) (d domain.Dataset, ok bool) {
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

	d, err = ctl.repo.Get(owner, ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseError(err))

		return
	}

	ok = true

	return
}
