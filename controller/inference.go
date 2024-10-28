package controller

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	spaceappApp "github.com/opensourceways/xihe-server/spaceapp/app"
	spaceappdomain "github.com/opensourceways/xihe-server/spaceapp/domain"
	spacemesage "github.com/opensourceways/xihe-server/spaceapp/domain/message"
	spaceappApprepo "github.com/opensourceways/xihe-server/spaceapp/domain/repository"
	userapp "github.com/opensourceways/xihe-server/user/app"
)

func AddRouterForInferenceController(
	rg *gin.RouterGroup,
	p platform.RepoFile,
	repo spaceappApprepo.Inference,
	project spacerepo.Project,
	sender message.Sender,
	whitelist userapp.WhiteListService,
	spacesender spacemesage.SpaceAppMessageProducer,
	appService spaceappApp.SpaceappAppService,
	spaceappRepo spaceappApprepo.SpaceAppRepository,
) {
	ctl := InferenceController{
		s: spaceappApp.NewInferenceService(
			p, repo, sender, apiConfig.MinSurvivalTimeOfInference, spacesender, spaceappRepo, project,
		),
		project:    project,
		whitelist:  whitelist,
		appService: appService,
	}

	ctl.inferenceDir, _ = domain.NewDirectory(apiConfig.InferenceDir)
	ctl.inferenceBootFile, _ = domain.NewFilePath(apiConfig.InferenceBootFile)

	rg.POST("/v1/inference", internalApiCheckMiddleware(&ctl.baseController), ctl.Create)
	rg.GET("/v1/inference/:owner/:name", internalApiCheckMiddleware(&ctl.baseController), ctl.Get)
	rg.GET("/v1/inference/:owner/:name/buildlog/complete", internalApiCheckMiddleware(&ctl.baseController), ctl.GetBuildLogs)
	rg.GET("/v1/inference/:owner/:name/buildlog/realtime", internalApiCheckMiddleware(&ctl.baseController), ctl.GetRealTimeBuildLog)
	rg.GET("/v1/inference/:owner/:name/spacelog/realtime", internalApiCheckMiddleware(&ctl.baseController), ctl.GetRealTimeSpaceLog)
	rg.GET("/v1/space-app/:owner/:name/read", ctl.CanRead)
}

type InferenceController struct {
	baseController

	s          spaceappApp.InferenceService
	appService spaceappApp.SpaceappAppService

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
func (ctl *InferenceController) Create(ctx *gin.Context) {
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

	if err := ctl.appService.Create(ctx, cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData("success"))
	}
}

func (ctl *InferenceController) getResourceLevel(owner domain.Account, pid string) (level string, err error) {
	resources, err := ctl.project.FindUserProjects(
		[]repository.UserResourceListOption{
			{
				Owner: owner,
				Ids: []string{
					pid,
				},
			},
		},
	)

	if err != nil || len(resources) < 1 {
		return
	}

	if resources[0].Level != nil {
		level = resources[0].Level.ResourceLevel()
	}

	return
}

// @Summary  Get
// @Description  get space app
// @Tags     SpaceApp
// @Param    owner  path  string  true  "owner of space" MaxLength(40)
// @Param    name   path  string  true  "name of space" MaxLength(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=app.SpaceAppDTO,msg=string,code=string}
// @Router   /v1/space-app/{owner}/{name} [get]
func (ctl *InferenceController) Get(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	cmd, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

	if dto, err := ctl.appService.GetByName(ctx.Request.Context(), pl.DomainAccount(), &cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, &dto)
	}
}

// parseIndex parses the index from the request.
func (ctl *InferenceController) parseIndex(ctx *gin.Context) (cmd spaceappApp.GetSpaceAppCmd, err error) {
	cmd.Owner, err = domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	cmd.Name, err = domain.NewResourceName(ctx.Param("name"))
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)
	}

	return
}

// @Summary  GetBuildLogs
// @Description  get space app complete build logs
// @Tags     SpaceApp
// @Param    id  path  string  true  "space app id"
// @Accept   json
// @Success  200  {object}  app.BuildLogsDTO
// @Router   /v1/space-app/{owner}/{name}/buildlog/complete [get]
func (ctl *InferenceController) GetBuildLogs(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

	if dto, err := ctl.appService.GetBuildLogs(ctx.Request.Context(), pl.DomainAccount(), &index); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, &dto)
	}
}

// @Summary  GetBuildLog
// @Description  get space app real-time build log
// @Tags     SpaceApp
// @Param    owner  path  string  true  "owner of space" MaxLength(40)
// @Param    name   path  string  true  "name of space" MaxLength(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=app.SpaceAppDTO,msg=string,code=string}
// @Router   /v1/space-app/{owner}/{name}/buildlog/realtime [get]
func (ctl *InferenceController) GetRealTimeBuildLog(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	index, err := ctl.parseIndex(ctx)
	if err != nil {
		ctx.SSEvent("error", err.Error())
		return
	}

	buildLog, err := ctl.appService.GetBuildLog(ctx.Request.Context(), pl.DomainAccount(), &index)
	if err != nil {
		logrus.Errorf("get build log err:%s", err)
		ctx.SSEvent("error", "get build log failed")
		return
	}

	streamWrite := func(doOnce func() ([]byte, error)) {
		ctx.Stream(func(w io.Writer) bool {
			done, err := doOnce()
			if err != nil {
				if err.Error() == "finish" {
					ctx.SSEvent("message", "")
				} else {
					logrus.Errorf("request build log err:%s", err)
					ctx.SSEvent("error", "request build log failed")
				}
				return false
			}
			if done != nil {
				ctx.SSEvent("message", string(done))
			}
			return true
		})
	}

	params := spaceappdomain.StreamParameter{
		StreamUrl: buildLog,
	}
	cmd := &spaceappdomain.SeverSentStream{
		Parameter:   params,
		Ctx:         ctx,
		StreamWrite: streamWrite,
	}

	if err := ctl.appService.GetRequestDataStream(cmd); err != nil {
		ctx.SSEvent("error", err.Error())
	}
}

// @Summary  GetSpaceLog
// @Description  get space app real-time space log
// @Tags     SpaceApp
// @Param    owner  path  string  true  "owner of space" MaxLength(40)
// @Param    name   path  string  true  "name of space" MaxLength(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData{data=app.SpaceAppDTO,msg=string,code=string}
// @Router   /v1/space-app/:owner/:name/spacelog/realtime [get]
func (ctl *InferenceController) GetRealTimeSpaceLog(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	index, err := ctl.parseIndex(ctx)
	if err != nil {
		ctx.SSEvent("error", err.Error())
		return
	}

	spaceLog, err := ctl.appService.GetSpaceLog(ctx.Request.Context(), pl.DomainAccount(), &index)
	if err != nil {
		logrus.Errorf("get space log err:%s", err)
		ctx.SSEvent("error", "get space log failed")
		return
	}

	streamWrite := func(doOnce func() ([]byte, error)) {
		ctx.Stream(func(w io.Writer) bool {
			done, err := doOnce()
			if err != nil {
				if err.Error() == "finish" {
					ctx.SSEvent("message", "")
				} else {
					logrus.Errorf("request space log err:%s", err)
					ctx.SSEvent("error", "request space log failed")
				}
				return false
			}
			if done != nil {
				ctx.SSEvent("message", string(done))
			}
			return true
		})
	}

	params := spaceappdomain.StreamParameter{
		StreamUrl: spaceLog,
	}
	cmd := &spaceappdomain.SeverSentStream{
		Parameter:   params,
		Ctx:         ctx,
		StreamWrite: streamWrite,
	}

	if err := ctl.appService.GetRequestDataStream(cmd); err != nil {
		ctx.SSEvent("error", err.Error())
	}
}

// @Summary  CanRead
// @Description  check permission for read space app
// @Tags     SpaceAppWeb
// @Param    owner  path  string  true  "owner of space" MaxLength(40)
// @Param    name   path  string  true  "name of space" MaxLength(100)
// @Accept   json
// @Success  200  {object}  commonctl.ResponseData
// @x-example {"data": "successfully"}
// @Router   /v1/space-app/{owner}/{name}/read [get]
func (ctl *InferenceController) CanRead(ctx *gin.Context) {
	index, err := ctl.parseIndex(ctx)
	if err != nil {
		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if err := ctl.appService.CheckPermissionRead(ctx.Request.Context(), pl.DomainAccount(), &index); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, "successfully")
	}
}
