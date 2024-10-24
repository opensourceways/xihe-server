package controller

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	"github.com/opensourceways/xihe-server/utils"
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
) {
	ctl := InferenceController{
		s: spaceappApp.NewInferenceService(
			p, repo, sender, apiConfig.MinSurvivalTimeOfInference, spacesender,
		),
		project:    project,
		whitelist:  whitelist,
		appService: appService,
	}

	ctl.inferenceDir, _ = domain.NewDirectory(apiConfig.InferenceDir)
	ctl.inferenceBootFile, _ = domain.NewFilePath(apiConfig.InferenceBootFile)

	rg.GET("/v1/inference/project/:owner/:pid", ctl.Create)
	rg.GET("/v1/space-app/:owner/:name", ctl.Get)
	rg.GET("/v1/space-app/:owner/:name/buildlog/complete", ctl.GetBuildLogs)
	rg.GET("/v1/space-app/:owner/:name/buildlog/realtime", ctl.GetRealTimeBuildLog)
	rg.GET("/v1/space-app/:owner/:name/spacelog/realtime", ctl.GetRealTimeSpaceLog)
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
	pl, csrftoken, _, ok := ctl.checkTokenForWebsocket(ctx, false)
	if !ok {
		return
	}

	// setup websocket
	upgrader := websocket.Upgrader{
		Subprotocols: []string{csrftoken},
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get(headerSecWebsocket) == csrftoken
		},
	}

	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		//TODO delete
		log.Errorf("update ws failed, err:%s", err.Error())

		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	defer ws.Close()

	// start
	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		if wsErr := ws.WriteJSON(newResponseCodeError(errorBadRequestParam, err)); wsErr != nil {
			log.Errorf("inference failed: web socket write err:%s", wsErr.Error())
		}

		log.Errorf("inference failed: new account, err:%s", err.Error())

		return
	}

	projectId := ctx.Param("pid")
	v, err := ctl.project.GetSummary(owner, projectId)
	if err != nil {
		if wsErr := ws.WriteJSON(newResponseError(err)); wsErr != nil {
			log.Errorf("inference get | web socket write err:%s", wsErr.Error())
		}

		log.Errorf("inference failed: get summary, err:%s", err.Error())

		return
	}

	viewOther := pl.isNotMe(owner)

	if v.IsPrivate() {
		wsErr := ws.WriteJSON(
			newResponseCodeMsg(
				errorNotAllowed,
				"project is not found",
			),
		)
		if wsErr != nil {
			log.Errorf("inference get | web socket write err:%s", wsErr.Error())
		}

		log.Debug("inference failed: project is private")

		return
	}

	var level string
	if level, err = ctl.getResourceLevel(owner, projectId); err != nil {
		if wsErr := ws.WriteJSON(newResponseError(err)); wsErr != nil {
			log.Errorf("inference get | web socket write err:%s", wsErr.Error())
		}

		log.Errorf("inference failed: get resource, err:%s", err.Error())

		return
	}

	u := platform.UserInfo{}
	if viewOther {
		u.User = owner
	} else {
		u = pl.PlatformUserInfo()
	}

	cmd := spaceappApp.InferenceCreateCmd{
		ProjectId:     v.Id,
		ProjectName:   v.Name,
		ProjectOwner:  owner,
		ResourceLevel: level,
		InferenceDir:  ctl.inferenceDir,
		BootFile:      ctl.inferenceBootFile,
	}

	dto, lastCommit, err := ctl.s.Create(pl.Account, &u, &cmd)
	if err != nil {
		if wsErr := ws.WriteJSON(newResponseError(err)); wsErr != nil {
			log.Errorf("inference get | web socket write err:%s", wsErr.Error())
		}

		log.Errorf("inference failed: create, err:%s", err.Error())

		return
	}

	utils.DoLog("", pl.Account, "create gradio",
		fmt.Sprintf("projectid: %s", v.Id), "success")

	if dto.Error != "" || dto.AccessURL != "" {
		if wsErr := ws.WriteJSON(newResponseData(dto)); wsErr != nil {
			log.Errorf("inference get | web socket write err:%s", wsErr.Error())
		}

		return
	}

	time.Sleep(10 * time.Second)

	info := spaceappApp.InferenceIndex{
		Id:         dto.InstanceId,
		LastCommit: lastCommit,
	}
	info.Project.Id = projectId
	info.Project.Owner = owner

	for i := 0; i < apiConfig.InferenceTimeout; i++ {
		dto, err = ctl.s.Get(&info)
		if err != nil {
			if wsErr := ws.WriteJSON(newResponseError(err)); wsErr != nil {
				log.Errorf("inference create | web socket write err:%s", wsErr.Error())
			}

			log.Errorf("inference failed: get status, err:%s", err.Error())

			return
		}

		log.Debugf("info dto:%v", dto)

		if dto.Error != "" || dto.AccessURL != "" {
			if wsErr := ws.WriteJSON(newResponseData(dto)); wsErr != nil {
				log.Errorf("inference create | web socket write err:%s", wsErr.Error())
			}

			log.Debug("inference done")

			return
		}

		time.Sleep(time.Second)
	}

	log.Error("inference timeout")

	if wsErr := ws.WriteJSON(newResponseCodeMsg(errorSystemError, "timeout")); wsErr != nil {
		logrus.Errorf("inference | web socket write error: %s", wsErr.Error())
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
