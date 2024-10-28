package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/common/domain/allerror"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	spaceappApp "github.com/opensourceways/xihe-server/spaceapp/app"
	"github.com/opensourceways/xihe-server/spaceapp/domain"
	spacemesage "github.com/opensourceways/xihe-server/spaceapp/domain/message"
	spaceappApprepo "github.com/opensourceways/xihe-server/spaceapp/domain/repository"
	userapp "github.com/opensourceways/xihe-server/user/app"
)

func AddRouterForInferenceInternalController(
	rg *gin.RouterGroup,
	p platform.RepoFile,
	project spacerepo.Project,
	sender message.Sender,
	whitelist userapp.WhiteListService,
	spacesender spacemesage.SpaceAppMessageProducer,
	spaceappRepo spaceappApprepo.SpaceAppRepository,
) {
	ctl := InferenceInternalController{
		s: spaceappApp.NewInferenceService(
			p, sender, apiConfig.MinSurvivalTimeOfInference, spacesender, spaceappRepo, project,
		),
		project:   project,
		whitelist: whitelist,
	}

	ctl.inferenceDir, _ = types.NewDirectory(apiConfig.InferenceDir)
	ctl.inferenceBootFile, _ = types.NewFilePath(apiConfig.InferenceBootFile)

	rg.POST("/v1/inference", internalApiCheckMiddleware(&ctl.baseController), ctl.Create)
	rg.PUT("/v1/inference/serving", internalApiCheckMiddleware(&ctl.baseController), ctl.NotifySpaceAppServing)
	rg.PUT("/v1/inference/building", internalApiCheckMiddleware(&ctl.baseController), ctl.NotifySpaceAppBuilding)
	rg.PUT("/v1/inference/starting", internalApiCheckMiddleware(&ctl.baseController), ctl.NotifySpaceAppStarting)
	rg.PUT("/v1/inference/failed_status", internalApiCheckMiddleware(&ctl.baseController),
		ctl.NotifySpaceAppFailedStatus)

}

type InferenceInternalController struct {
	baseController

	s spaceappApp.InferenceService

	project spacerepo.Project

	inferenceDir      types.Directory
	inferenceBootFile types.FilePath
	whitelist         userapp.WhiteListService
}

// @Summary		Create
// @Description	create inference
// @Tags			Inference
// @Param			owner	path	string	true	"project owner"
// @Param			pid		path	string	true	"project id"
// @Accept			json
// @Success		201	{object}			app.InferenceDTO
// @Failure		400	bad_request_body	can't	parse		request	body
// @Failure		401	bad_request_param	some	parameter	of		body	is	invalid
// @Failure		500	system_error		system	error
// @Router			/v1/inference/project/{owner}/{pid} [get]
func (ctl *InferenceInternalController) Create(ctx *gin.Context) {
	req := reqToCreateSpaceApp{}
	if err := ctx.BindJSON(&req); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	fmt.Printf("=====================cmd: %+v\n", cmd)

	if err := ctl.s.Create(ctx, cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData("success"))
	}
}

// @Summary  NotifySpaceAppServing
// @Description  notify space app service is started
// @Tags     SpaceApp
// @Param    body  body  reqToUpdateServiceInfo  true  "body"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Security Internal
// @Router   /v1/inference/serving [put]
func (ctl *InferenceInternalController) NotifySpaceAppServing(ctx *gin.Context) {
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

	if err := ctl.s.NotifyIsServing(ctx.Request.Context(), &cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfPut(ctx, nil)
	}
}

// @Summary  NotifySpaceAppBuilding
// @Description  notify space app building is started
// @Tags     SpaceApp
// @Param    body  body  reqToUpdateBuildInfo  true  "body"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Security Internal
// @Router   /v1/infernce/building [put]
func (ctl *InferenceInternalController) NotifySpaceAppBuilding(ctx *gin.Context) {
	req := reqToUpdateBuildInfo{}

	if err := ctx.BindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.s.NotifyIsBuilding(ctx.Request.Context(), &cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfPut(ctx, nil)
	}
}

// @Summary  NotifySpaceAppStarting
// @Description  notify space app build is starting
// @Tags     SpaceApp
// @Param    body  body  reqToNotifyStarting  true  "body"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Security Internal
// @Router   /v1/inference/starting [put]
func (ctl *InferenceInternalController) NotifySpaceAppStarting(ctx *gin.Context) {
	req := reqToNotifyStarting{}

	if err := ctx.BindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.s.NotifyStarting(ctx.Request.Context(), &cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfPut(ctx, nil)
	}
}

// @Summary  NotifySpaceAppFailedStatus
// @Description  notify space app failed status
// @Tags     SpaceApp
// @Param    body  body  reqToFailedStatus  true  "body"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Security Internal
// @Router   /v1/inference/failed_status [put]
func (ctl *InferenceInternalController) NotifySpaceAppFailedStatus(ctx *gin.Context) {
	req := reqToFailedStatus{}

	if err := ctx.BindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	switch cmd.Status {
	case domain.AppStatusBuildFailed:
		if err := ctl.s.NotifyIsBuildFailed(ctx.Request.Context(), &cmd); err != nil {
			ctl.sendRespWithInternalError(ctx, newResponseError(err))
			return
		}
	case domain.AppStatusStartFailed:
		if err := ctl.s.NotifyIsStartFailed(ctx.Request.Context(), &cmd); err != nil {
			ctl.sendRespWithInternalError(ctx, newResponseError(err))
			return
		}
	default:
		e := fmt.Errorf("old status not %s, can not set", cmd.Status.AppStatus())
		err = allerror.New(allerror.ErrorCodeSpaceAppUnmatchedStatus, e.Error(), e)
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
		return
	}
	ctl.sendRespOfPut(ctx, nil)
}
