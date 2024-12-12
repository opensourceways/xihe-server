/*
Copyright (c) Huawei Technologies Co., Ltd. 2024. All rights reserved
*/

// Package controller provides controller
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/computility/app"
	"github.com/opensourceways/xihe-server/computility/domain"
	commondomain "github.com/opensourceways/xihe-server/domain"
)

// AddRouterForComputilityWebController adds routes to
// the given router group for the AddRouterForComputilityWebController.
func AddRouterForComputilityWebController(
	rg *gin.RouterGroup,
	s app.ComputilityAppService,
) {
	ctl := ComputilityWebController{
		appService: s,
	}

	rg.GET("/v1/computility/account/:type", ctl.GetComputilityAccountDetail)
}

// ComputilityWebController is a struct that holds the necessary dependencies for
// handling computility-related operations.
type ComputilityWebController struct {
	baseController
	appService app.ComputilityAppService
}

// @Summary  GetComputilityAccountDetail
// @Description  get user computility account detail
// @Tags     ComputilityWeb
// @Param    type   path  string  true  "computility type"
// @Accept   json
// @Success  200  {object} commonctl.ResponseData{data=app.AccountQuotaDetailDTO,msg=string,code=string}
// @Failure  400  {object} commonctl.ResponseData{data=error,msg=string,code=string}
// @Router   /v1/computility/account/{type} [get]
func (ctl *ComputilityWebController) GetComputilityAccountDetail(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	compute, err := commondomain.NewComputilityType(ctx.Param("type"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	r, err := ctl.appService.GetAccountDetail(domain.ComputilityAccountIndex{
		UserName:    pl.DomainAccount(),
		ComputeType: compute,
	})
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, r)
	}
}
