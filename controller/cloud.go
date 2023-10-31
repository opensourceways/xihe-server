package controller

import (
	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/cloud/app"
)

func AddRouterForCloudController(
	rg *gin.RouterGroup,
	s app.CloudService,
) {
	ctl := CloudController{
		s: s,
	}

	rg.GET("/v1/cloud", ctl.List)
}

type CloudController struct {
	baseController

	s app.CloudService
}

// @Summary		List
// @Description	list cloud config
// @Tags			Cloud
// @Accept			json
// @Success		200	{object}		[]app.CloudDTO
// @Failure		500	system_error	system	error
// @Router			/v1/cloud [get]
func (ctl *CloudController) List(ctx *gin.Context) {
	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	cmd := new(app.GetCloudConfCmd)
	if visitor {
		cmd.ToCmd(nil, visitor)
	} else {
		cmd.ToCmd(pl.DomainAccount(), visitor)
	}

	data, err := ctl.s.ListCloud(cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}
