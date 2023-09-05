package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/points/app"
)

func AddRouterForUserPointsController(
	rg *gin.RouterGroup,
	s app.UserPointsAppService,
) {
	ctl := UserPointsController{
		s: s,
	}

	rg.GET("/v1/user_points", ctl.PointsDetails)
}

type UserPointsController struct {
	baseController

	s app.UserPointsAppService
}

//	@Summary		get user points details
//	@Description		get user points details
//	@Tags			UserPoints
//	@Accept			json
//	@Success		200	{object}	app.UserPointsDetailsDTO
//	@Failure		500	system_error	system	error
//	@Router			/v1/user_points [get]
func (ctl *UserPointsController) PointsDetails(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if v, err := ctl.s.GetPointsDetails(pl.DomainAccount()); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}
