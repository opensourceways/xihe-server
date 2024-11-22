package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/filescan/app"
)

// AddRouterForFileScanController add router for filescan
func AddRouterForFileScanInternalController(
	rg *gin.RouterGroup,
	f app.FileScanService,
) {
	ctl := FileScanController{
		fileScanService: f,
	}

	rg.PUT("/internal/repo/filescan/:id", ctl.Update)
}

// FileScanController is the controller of filescan
type FileScanController struct {
	baseController
	fileScanService app.FileScanService
}

func (ctl *FileScanController) Update(ctx *gin.Context) {
	req := ReqToUpdateFileScan{}
	if err := ctx.BindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	cmd, err := req.ToCmdToUpdateFileScan(ctx.Param("id"))
	if err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	if err = ctl.fileScanService.Update(ctx.Request.Context(), cmd); err != nil {
		ctl.sendBadRequestBody(ctx)
	} else {
		ctl.sendRespOfPut(ctx, nil)
	}
}
