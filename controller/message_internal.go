package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
)

func AddRouterForMessageInternalController(
	rg *gin.RouterGroup,
	message app.ProjectMessageService,
	repo spacerepo.Project,
	repoPg spacerepo.ProjectPg,
) {
	ctl := MessageInternalController{
		message: message,
	}
	rg.PUT("/like/", ctl.ChangeProjectLike)
	rg.PUT("/fork/", ctl.IncreaseFork)
	rg.PUT("/download/", ctl.IncreaseDownload)
}

type MessageInternalController struct {
	baseController
	message app.ProjectMessageService
}

func (ctl *MessageInternalController) ChangeProjectLike(ctx *gin.Context) {
	req := reqToChangeProjectLike{}
	if err := ctx.BindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequestBody(ctx)
		return
	}

	change := cmd.ChangeNum

	if change < 0 {
		if err = ctl.message.RemoveLike(&cmd.ResourceIndex); err != nil {
			ctl.sendRespWithInternalError(ctx, newResponseError(err))
		} else {
			ctx.JSON(http.StatusOK, newResponseData(""))
		}
	}
	if change > 0 {
		if err = ctl.message.AddLike(&cmd.ResourceIndex); err != nil {
			ctl.sendRespWithInternalError(ctx, newResponseError(err))
		} else {
			ctx.JSON(http.StatusOK, newResponseData(""))
		}
	}
}

func (ctl *MessageInternalController) IncreaseFork(ctx *gin.Context) {
	req := reqToIncrease{}
	if err := ctx.BindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequestBody(ctx)
		return
	}

	if err = ctl.message.IncreaseFork(&cmd.ResourceIndex); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(""))
	}
}

func (ctl *MessageInternalController) IncreaseDownload(ctx *gin.Context) {
	req := reqToIncrease{}
	if err := ctx.BindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)
		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequestBody(ctx)
		return
	}

	if err = ctl.message.IncreaseDownload(&cmd.ResourceIndex); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(""))
	}
}
