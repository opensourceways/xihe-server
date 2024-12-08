package controller

import (
	"net/http"

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

	rg.PATCH("/v1/repo/filescan/:id", internalApiCheckMiddleware(&ctl.baseController), ctl.Update)
	rg.PATCH("/v1/repo/filescan", internalApiCheckMiddleware(&ctl.baseController), ctl.LaunchModeration)
	rg.POST("/v1/repo/filescan", internalApiCheckMiddleware(&ctl.baseController), ctl.CreateList)
	rg.DELETE("/v1/repo/filescan", internalApiCheckMiddleware(&ctl.baseController), ctl.Delete)
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

func (ctl *FileScanInternalController) LaunchModeration(ctx *gin.Context) {
	req := ModifyFileScanListReq{}
	if err := ctx.BindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	if err = ctl.fileScanService.LaunchModeration(ctx.Request.Context(), cmd); err != nil {
		ctl.sendBadRequestBody(ctx)
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (ctl *FileScanInternalController) CreateList(ctx *gin.Context) {
	req := CreateFileScanListReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequestBody(ctx)
		return
	}

	if err = ctl.fileScanService.CreateList(ctx.Request.Context(), cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	}

	ctx.JSON(http.StatusAccepted, nil)
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
