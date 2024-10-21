package controller

import (
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
) {
	ctl := InferenceInternalController{
		s: spaceappApp.NewInferenceService(
			p, repo, sender, apiConfig.MinSurvivalTimeOfInference, spacesender,
		),
		project:   project,
		whitelist: whitelist,
	}

	ctl.inferenceDir, _ = domain.NewDirectory(apiConfig.InferenceDir)
	ctl.inferenceBootFile, _ = domain.NewFilePath(apiConfig.InferenceBootFile)

	rg.GET("/v1/inference/project/:owner/:pid", internalApiCheckMiddleware(&ctl.baseController), ctl.Create)
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

	if err := ctl.s.CreateSpaceApp(cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	}
}
