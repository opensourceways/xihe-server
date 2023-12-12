package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/promotion/app"
	pc "github.com/opensourceways/xihe-server/promotion/controller"
)

func AddRouterForPromotionController(
	rg *gin.RouterGroup,
	pros app.PromotionService,
	ps app.PointsService,
) {
	ctl := PromotionController{
		pros: pros,
		ps:   ps,
	}

	rg.POST("/v1/promotion/:id/apply", ctl.Apply)
}

type PromotionController struct {
	baseController

	pros app.PromotionService
	ps   app.PointsService
}

// @Summary		Apply
// @Description	apply the Promotion
// @Tags			Promotion
// @Param			id		path	string						true	"promotion id"
// @Param			body	body	pc.PromotionApplyReq	true	"body of applying"
// @Accept			json
// @Success		201
// @Failure		500	system_error	system	error
// @Router			/v1/promotion/{id}/apply [post]
func (ctl *PromotionController) Apply(ctx *gin.Context) {
	req := pc.PromotionApplyReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "apply the promotion")

	cmd, err := req.ToCmd(ctx.Param("id"), pl.DomainAccount())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestParam(err))

		return
	}

	if code, err := ctl.pros.UserRegister(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, "success")
	}
}
