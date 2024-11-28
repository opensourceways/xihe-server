package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
)

func AddRouterForLikeController(
	rg *gin.RouterGroup,
	repo repository.Like,
	user userrepo.User,
	proj spacerepo.Project,
	projPg spacerepo.ProjectPg,
	model repository.Model,
	dataset repository.Dataset,
	activity repository.Activity,
	sender message.LikeMessageProducer,
) {
	ctl := LikeController{
		s: app.NewLikeService(
			repo, user, model, proj,
			dataset, activity, sender,
		),
		proj:    proj,
		projPg:  projPg,
		model:   model,
		dataset: dataset,
	}

	rg.POST("/v1/user/like", ctl.Create)
	rg.DELETE("/v1/user/like", ctl.Delete)
	rg.GET("/v1/user/like/:account", ctl.List)
}

type LikeController struct {
	baseController

	s app.LikeService

	proj    spacerepo.Project
	projPg  spacerepo.ProjectPg
	model   repository.Model
	dataset repository.Dataset
}

// @Title			Create
// @Description	create a like
// @Tags			Like
// @Param			body	body	likeCreateRequest	true	"body of creating like"
// @Accept			json
// @Success		201
// @Failure		400	bad_request_body	can't	parse		request		body
// @Failure		401	bad_request_param	some	parameter	of			body	is		invalid
// @Failure		402	not_allowed			can't	add			yourself	as		your	like
// @Failure		403	resource_not_exists	the		target		of			like	does	not	exist
// @Failure		404	duplicate_creating	add		like		again
// @Failure		500	system_error		system	error
// @Router			/v1/user/like [post]
func (ctl *LikeController) Create(ctx *gin.Context) {
	req := likeCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}
	fmt.Printf("=========================req: %+v\n", req)
	cmd, ok := req.toCmd(ctx, ctl.getResourceId)
	if !ok {
		return
	}
	fmt.Printf("==========================cmd: %+v\n", cmd)
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}
	if !pl.isNotMe(cmd.ResourceOwner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed, "can't like yourself",
		))

		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "create a like")

	if err := ctl.s.Create(pl.DomainAccount(), cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData("success"))
	}
}

// @Title			Delete
// @Description	delete a like
// @Tags			Like
// @Param			body	body	likeDeleteRequest	true	"body of deleting like"
// @Accept			json
// @Success		204
// @Failure		400	bad_request_body	can't	parse		request		body
// @Failure		401	bad_request_param	some	parameter	of			body	is		invalid
// @Failure		402	not_allowed			can't	add			yourself	as		your	like
// @Failure		403	resource_not_exists	the		target		of			like	does	not	exist
// @Failure		500	system_error		system	error
// @Router			/v1/user/like [delete]
func (ctl *LikeController) Delete(ctx *gin.Context) {
	req := likeDeleteRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	cmd, ok := req.toCmd(ctx, ctl.getResourceId)
	if !ok {
		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}
	if !pl.isNotMe(cmd.ResourceOwner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed, "can't delete like of yourself",
		))

		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "delete a like")

	if err := ctl.s.Delete(pl.DomainAccount(), cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusNoContent, newResponseData("success"))
	}
}

// @Title			List
// @Description	list likes
// @Tags			Like
// @Param			account	path	string	true	"the account the likes belong to"
// @Accept			json
// @Success		200	{object}		app.LikeDTO
// @Failure		500	system_error	system	error
// @Router			/v1/user/like/{account} [get]
func (ctl *LikeController) List(ctx *gin.Context) {
	// TODO: list by page
	account, err := domain.NewAccount(ctx.Param("account"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if data, err := ctl.s.List(account); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}

func (ctl *LikeController) getResourceId(
	owner domain.Account, rt domain.ResourceType, name domain.ResourceName,
) (string, error) {
	switch rt.ResourceType() {
	case domain.ResourceTypeProject.ResourceType():
		v, err := ctl.projPg.GetByName(owner, name)
		if err != nil {
			return "", err
		}

		return v.Id, nil

	case domain.ResourceTypeDataset.ResourceType():
		v, err := ctl.dataset.GetByName(owner, name)
		if err != nil {
			return "", err
		}

		return v.Id, nil

	case domain.ResourceTypeModel.ResourceType():
		v, err := ctl.model.GetByName(owner, name)
		if err != nil {
			return "", err
		}

		return v.Id, nil

	default:
		return "", errors.New("unknown resource type")

	}

}
