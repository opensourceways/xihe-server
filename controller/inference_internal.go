package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	spaceappApp "github.com/opensourceways/xihe-server/spaceapp/app"
	spacemesage "github.com/opensourceways/xihe-server/spaceapp/domain/message"
	spaceappApprepo "github.com/opensourceways/xihe-server/spaceapp/domain/repository"
	userapp "github.com/opensourceways/xihe-server/user/app"
)

func AddRouterForInferenceInternalController(
	rg *gin.RouterGroup,
	p platform.RepoFile,
	repo spaceappApprepo.Inference,
	project spacerepo.Project,
	sender message.Sender,
	whitelist userapp.WhiteListService,
	spacesender spacemesage.SpaceAppMessageProducer,
	spaceappRepo spaceappApprepo.SpaceAppRepository,
) {
	ctl := InferenceInternalController{
		s: spaceappApp.NewInferenceService(
			p, repo, sender, apiConfig.MinSurvivalTimeOfInference, spacesender, spaceappRepo, project,
		),
		project:   project,
		whitelist: whitelist,
	}

	ctl.inferenceDir, _ = domain.NewDirectory(apiConfig.InferenceDir)
	ctl.inferenceBootFile, _ = domain.NewFilePath(apiConfig.InferenceBootFile)

	rg.POST("/v1/inference", internalApiCheckMiddleware(&ctl.baseController), ctl.Create)
	rg.PUT("/v1/inference/serving", internalApiCheckMiddleware(&ctl.baseController), ctl.NotifySpaceAppServing)
}

type InferenceInternalController struct {
	baseController

	s spaceappApp.InferenceService

	project spacerepo.Project

	inferenceDir      domain.Directory
	inferenceBootFile domain.FilePath
	whitelist         userapp.WhiteListService

	spacesender spacemesage.SpaceAppMessageProducer
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

	if err := ctl.s.CreateSpaceApp(cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData("string"))
	}
}

// @Summary  NotifySpaceAppServing
// @Description  notify space app service is started
// @Tags     SpaceApp
// @Param    body  body  reqToUpdateServiceInfo  true  "body"
// @Accept   json
// @Success  202   {object}  commonctl.ResponseData{data=nil,msg=string,code=string}
// @Security Internal
// @Router   /v1/space-app/serving [put]
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
