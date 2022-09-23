package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForActivityController(
	rg *gin.RouterGroup,
	repo repository.Activity,
	user repository.User,
	proj repository.Project,
	model repository.Model,
	dataset repository.Dataset,
) {
	ctl := ActivityController{
		s: app.NewActivityService(repo, user, model, proj, dataset),
	}

	rg.GET("/v1/user/activity/:account", ctl.List)
}

type ActivityController struct {
	baseController

	s app.ActivityService
}

// @Title List
// @Description list activitys
// @Tags  Activity
// @Param	account	path	string	true	"the account the activities belong to"
// @Accept json
// @Success 200 {object} app.ActivityDTO
// @Failure 500 system_error        system error
// @Router /v1/user/activity/{account} [get]
func (ctl *ActivityController) List(ctx *gin.Context) {
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
