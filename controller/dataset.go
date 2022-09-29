package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForDatasetController(
	rg *gin.RouterGroup,
	user repository.User,
	repo repository.Dataset,
	model repository.Model,
	proj repository.Project,
	activity repository.Activity,
	tags repository.Tags,
	like repository.Like,
	newPlatformRepository func(token, namespace string) platform.Repository,
) {
	ctl := DatasetController{
		repo: repo,
		tags: tags,
		like: like,
		s:    app.NewDatasetService(user, repo, proj, model, activity, nil),

		newPlatformRepository: newPlatformRepository,
	}

	rg.POST("/v1/dataset", ctl.Create)
	rg.PUT("/v1/dataset/:owner/:id", ctl.Update)
	rg.GET("/v1/dataset/:owner/:name", ctl.Get)
	rg.GET("/v1/dataset/:owner", ctl.List)

	rg.PUT("/v1/dataset/:owner/:id/tags", ctl.SetTags)
}

type DatasetController struct {
	baseController

	repo repository.Dataset
	tags repository.Tags
	like repository.Like
	s    app.DatasetService

	newPlatformRepository func(string, string) platform.Repository
}

// @Summary Create
// @Description create dataset
// @Tags  Dataset
// @Param	body	body 	datasetCreateRequest	true	"body of creating dataset"
// @Accept json
// @Success 201 {object} app.DatasetDTO
// @Failure 400 bad_request_body    can't parse request body
// @Failure 400 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Failure 500 duplicate_creating  create dataset repeatedly
// @Router /v1/dataset [post]
func (ctl *DatasetController) Create(ctx *gin.Context) {
	req := datasetCreateRequest{}
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

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

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
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(d))
}

// @Summary Update
// @Description update property of dataset
// @Tags  Dataset
// @Param	id	path	string			true	"id of dataset"
// @Param	body	body 	datasetUpdateRequest	true	"body of updating dataset"
// @Accept json
// @Produce json
// @Router /v1/dataset/{owner}/{id} [put]
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
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(d))
}

// @Summary Get
// @Description get dataset
// @Tags  Dataset
// @Param	owner	path	string	true	"owner of dataset"
// @Param	name	path	string	true	"name of dataset"
// @Accept json
// @Success 200 {object} datasetDetail
// @Produce json
// @Router /v1/dataset/{owner}/{name} [get]
func (ctl *DatasetController) Get(ctx *gin.Context) {
	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	name, err := domain.NewDatasetName(ctx.Param("name"))
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

	ctx.JSON(http.StatusOK, newResponseData(datasetDetail{
		Liked:            liked,
		DatasetDetailDTO: &d,
	}))
}

// @Summary List
// @Description list dataset
// @Tags  Dataset
// @Param	owner		path	string	true	"owner of dataset"
// @Param	name		query	string	false	"name of dataset"
// @Param	repo_type	query	string	false	"repo type of dataset, value can be public or private"
// @Param	count_per_page	query	int	false	"count per page"
// @Param	page_num	query	int	false	"page num which starts from 1"
// @Param	sort_by		query	string	false	"sort keys, value can be update_time, first_letter, download_count"
// @Accept json
// @Success 200 {object} app.DatasetsDTO
// @Produce json
// @Router /v1/dataset/{owner} [get]
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
			cmd.RepoType, _ = domain.NewRepoType(domain.RepoTypePublic)
		} else {
			if cmd.RepoType.RepoType() != domain.RepoTypePublic {
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

	ctx.JSON(http.StatusOK, newResponseData(data))
}

// @Summary SetTags
// @Description set tags for dataset
// @Tags  Dataset
// @Param	owner	path	string				true	"owner of dataset"
// @Param	id	path	string				true	"id of dataset"
// @Param	body	body 	resourceTagsUpdateRequest	true	"body of tags"
// @Accept json
// @Success 202
// @Router /v1/dataset/{owner}/{id}/tags [put]
func (ctl *DatasetController) SetTags(ctx *gin.Context) {
	req := resourceTagsUpdateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	tags, err := ctl.tags.List(domain.ResourceTypeDataset)
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
			errorNotAllowed,
			"can't update dataset for other user",
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
