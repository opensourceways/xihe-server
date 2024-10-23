package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/spaceapp/app"
)

// AddRouteForSpaceappInternalController adds routes for SpaceAppInternalController to the given router group.
func AddRouteForSpaceappInternalController(
	rg *gin.RouterGroup,
	s app.SpaceappInternalAppService,
) {

	ctl := SpaceAppInternalController{
		appService: s,
	}

	rg.PUT(`/v1/space-app/serving`, ctl.NotifySpaceAppServing)
}

// SpaceAppInternalController is a struct that holds the app service
// and provides methods for handling requests related to space apps.
type SpaceAppInternalController struct {
	baseController

	appService app.SpaceappInternalAppService
}

// @Summary  NotifySpaceAppServing
// @Description  notify space app service is started
// @Tags     SpaceApp
// @Param    body  body  reqToUpdateServiceInfo  true  "body"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Security Internal
// @Router   /v1/space-app/serving [put]
func (ctl *SpaceAppInternalController) NotifySpaceAppServing(ctx *gin.Context) {
	req := reqToUpdateServiceInfo{}

	if err := ctx.BindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.appService.NotifyIsServing(ctx.Request.Context(), &cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfPut(ctx, nil)
	}
}
