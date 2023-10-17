package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/points/app"
)

func AddRouterForUserPointsController(
	rg *gin.RouterGroup,
	s app.UserPointsAppService,
	ts app.TaskAppService,
) {
	ctl := UserPointsController{
		s:  s,
		ts: ts,
	}

	rg.GET("/v1/user_points", ctl.PointsDetails)
	rg.GET("/v1/user_points/tasks", ctl.TasksOfDay)
	rg.GET("/v1/user_points/taskdoc", ctl.TasksDoc)
}

type UserPointsController struct {
	baseController

	s  app.UserPointsAppService
	ts app.TaskAppService
}

// @Summary		get user points details
// @Description		get user points details
// @Tags			UserPoints
// @Accept			json
// @Success		200	{object}	app.UserPointsDetailsDTO
// @Failure		500	system_error	system	error
// @Router			/v1/user_points [get]
func (ctl *UserPointsController) PointsDetails(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	lang, err := ctl.languageRuquested(ctx)
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if v, err := ctl.s.PointsDetails(pl.DomainAccount(), lang); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}

// @Summary		tasks of day
// @Description		tasks of day
// @Tags			UserPoints
// @Accept			json
// @Success		200	{object}	app.TasksCompletionInfoDTO
// @Failure		500	system_error	system	error
// @Router			/v1/user_points/tasks [get]
func (ctl *UserPointsController) TasksOfDay(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	lang, err := ctl.languageRuquested(ctx)
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if v, err := ctl.s.TasksOfDay(pl.DomainAccount(), lang); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}

// @Summary		task doc
// @Description		task doc
// @Tags			UserPoints
// @Accept			json
// @Success		200	{object}	app.TaskDocDTO
// @Failure		500	system_error	system	error
// @Router			/v1/user_points/taskdoc [get]
func (ctl *UserPointsController) TasksDoc(ctx *gin.Context) {
	lang, err := ctl.languageRuquested(ctx)
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if v, err := ctl.ts.Doc(lang); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}
