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
	ctl := FileScanInternalController{
		fileScanService: f,
	}

	rg.PUT("repo/filescan/:id", ctl.Update)
}

// FileScanController is the controller of filescan
type FileScanInternalController struct {
	baseController
	fileScanService app.FileScanService
}

func (ctl *FileScanInternalController) Update(ctx *gin.Context) {
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

func (ctl *FileScanInternalController) Create(ctx *gin.Context) {
	req := CreateFileScansReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)
		return
	}
}

func (ctl *FileScanInternalController) Delete(ctx *gin.Context) {
	req := RemoveFileScansReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequestBody(ctx)
		return
	}

	if err = ctl.fileScanService.Remove(ctx.Request.Context(), cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
		return
	}

	ctl.sendRespOfDelete(ctx)
}
