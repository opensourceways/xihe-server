package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/domain/training"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

func AddRouterForTrainingController(
	rg *gin.RouterGroup,
	ts training.Training,
	repo repository.Training,
	model repository.Model,
	projectPg spacerepo.ProjectPg,
	dataset repository.Dataset,
	sender message.MessageProducer,
) {
	ctl := TrainingController{
		ts: app.NewTrainingService(
			ts, repo, sender, apiConfig.MaxTrainingRecordNum,
		),
		model:     model,
		projectPg: projectPg,
		dataset:   dataset,
	}

	rg.POST("/v1/train/project/:pid/training", checkUserEmailMiddleware(&ctl.baseController), ctl.Create)
	rg.POST("/v1/train/project/:pid/training/:id", ctl.Recreate)
	rg.PUT("/v1/train/project/:pid/training/:id", ctl.Terminate)
	rg.GET("/v1/train/project/:pid/training", checkUserEmailMiddleware(&ctl.baseController), ctl.List)
	rg.GET("/v1/train/project/:pid/training/ws", ctl.ListByWS)
	rg.GET(
		"/v1/train/project/:pid/training/:id/result/:type", checkUserEmailMiddleware(&ctl.baseController),
		ctl.GetResultDownloadURL,
	)
	rg.GET("/v1/train/project/:pid/training/:id", ctl.Get)
	rg.GET("/v1/train/project/:pid/config", ctl.GetLastTrainingConfig)
	rg.DELETE("v1/train/project/:pid/training/:id", ctl.Delete)
}

type TrainingController struct {
	baseController

	ts app.TrainingService

	model     repository.Model
	projectPg spacerepo.ProjectPg
	dataset   repository.Dataset
}

// @Summary		Create
// @Description	create training
// @Tags			Training
// @Param			pid		path	string					true	"project id"
// @Param			body	body	TrainingCreateRequest	true	"body of creating training"
// @Accept			json
// @Success		201	{object}			trainingCreateResp
// @Failure		400	bad_request_body	can't	parse		request	body
// @Failure		401	bad_request_param	some	parameter	of		body	is	invalid
// @Failure		500	system_error		system	error
// @Router			/v1/train/project/{pid}/training [post]
func (ctl *TrainingController) Create(ctx *gin.Context) {
	req := TrainingCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "create training")

	cmd := new(app.TrainingCreateCmd)

	if err := req.toCmd(cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if !ctl.setProjectInfo(ctx, cmd, pl.DomainAccount(), ctx.Param("pid")) {
		return
	}

	if !ctl.setModelsInput(ctx, cmd, pl.DomainAccount(), req.Models) {
		return
	}

	if !ctl.setDatasetsInput(ctx, cmd, pl.DomainAccount(), req.Datasets) {
		return
	}

	if err := cmd.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	v, err := ctl.ts.Create(cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	utils.DoLog("", pl.Account, "create training",
		fmt.Sprintf("projectid: %s, trainingid: %s", ctx.Param("pid"), v), "success")

	ctx.JSON(http.StatusCreated, newResponseData(trainingCreateResp{v}))
}

// @Summary		Recreate
// @Description	recreate training
// @Tags			Training
// @Param			pid	path	string	true	"project id"
// @Param			id	path	string	true	"training id"
// @Accept			json
// @Success		201	{object}			trainingCreateResp
// @Failure		400	bad_request_body	can't	parse		request	body
// @Failure		401	bad_request_param	some	parameter	of		body	is	invalid
// @Failure		500	system_error		system	error
// @Router			/v1/train/project/{pid}/training/{id} [post]
func (ctl *TrainingController) Recreate(ctx *gin.Context) {
	info, ok := ctl.getTrainingInfo(ctx)
	if !ok {
		return
	}

	pl, _, _ := ctl.checkUserApiToken(ctx, false)
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "recreate training")

	v, err := ctl.ts.Recreate(&info)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	utils.DoLog("", info.Project.Owner.Account(), "recreate training",
		fmt.Sprintf("projectid: %s, trainingid: %s", info.Project.Id, info.TrainingId), "success")

	ctx.JSON(http.StatusCreated, newResponseData(trainingCreateResp{v}))
}

// @Summary		Delete
// @Description	delete training
// @Tags			Training
// @Param			pid	path	string	true	"project id"
// @Param			id	path	string	true	"training id"
// @Accept			json
// @Success		204
// @Failure		500	system_error	system	error
// @Router			/v1/train/project/{pid}/training/{id} [delete]
func (ctl *TrainingController) Delete(ctx *gin.Context) {
	info, ok := ctl.getTrainingInfo(ctx)
	if !ok {
		return
	}

	pl, _, _ := ctl.checkUserApiToken(ctx, false)
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "delete training")

	if err := ctl.ts.Delete(&info); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	utils.DoLog("", info.Project.Owner.Account(), "delete training",
		fmt.Sprintf("projectid: %s, trainingid: %s", info.Project.Id, info.TrainingId), "success")

	ctx.JSON(http.StatusNoContent, newResponseData("success"))
}

// @Summary		Terminate
// @Description	terminate training
// @Tags			Training
// @Param			pid	path	string	true	"project id"
// @Param			id	path	string	true	"training id"
// @Accept			json
// @Success		202
// @Failure		500	system_error	system	error
// @Router			/v1/train/project/{pid}/training/{id} [put]
func (ctl *TrainingController) Terminate(ctx *gin.Context) {
	info, ok := ctl.getTrainingInfo(ctx)
	if !ok {
		return
	}

	pl, _, _ := ctl.checkUserApiToken(ctx, false)
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "terminate training")

	if err := ctl.ts.Terminate(&info); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	utils.DoLog("", info.Project.Owner.Account(), "terminate training",
		fmt.Sprintf("projectid: %s, trainingid: %s", info.Project.Id, info.TrainingId), "success")

	ctx.JSON(http.StatusAccepted, newResponseData("success"))
}

// @Summary		Get
// @Description	get training info
// @Tags			Training
// @Param			pid	path	string	true	"project id"
// @Param			id	path	string	true	"training id"
// @Accept			json
// @Success		200	{object}		trainingDetail
// @Failure		500	system_error	system	error
// @Router			/v1/train/project/{pid}/training/{id} [get]
func (ctl *TrainingController) Get(ctx *gin.Context) {
	pl, csrftoken, _, ok := ctl.checkTokenForWebsocket(ctx, false)
	if !ok {
		return
	}

	index := domain.TrainingIndex{
		Project: domain.ResourceIndex{
			Owner: pl.DomainAccount(),
			Id:    ctx.Param("pid"),
		},
		TrainingId: ctx.Param("id"),
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
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	defer ws.Close()

	ctl.watchTraining(ws, &index)
}

func (ctl *TrainingController) watchTraining(ws *websocket.Conn, index *domain.TrainingIndex) {
	duration := 0
	sleep := func() {
		time.Sleep(time.Second)

		if duration > 0 {
			duration++
		}
	}

	data := &trainingDetail{}

	start, end := 4, 5
	i := start
	for {
		if i++; i == end {
			v, code, err := ctl.ts.Get(index)
			if err != nil {
				if code == app.ErrorTrainNotFound {
					break
				}

				i = start
				sleep()

				continue
			}

			data.TrainingDTO = v

			if duration == 0 {
				duration = v.Duration
			} else {
				data.Duration = duration
			}

			log, err := downloadLog(v.LogPreviewURL)
			if err == nil && len(log) > 0 {
				data.Log = string(log)
			}

			if err = ws.WriteJSON(newResponseData(data)); err != nil {
				break
			}

			if v.IsDone {
				break
			}

			i = 0
		} else {
			if data.Duration > 0 {
				data.Duration++

				if err := ws.WriteJSON(newResponseData(data)); err != nil {
					break
				}
			}
		}

		sleep()
	}
}

// @Summary		List
// @Description	get trainings
// @Tags			Training
// @Param			pid	path	string	true	"project id"
// @Accept			json
// @Success		200	{object}		app.TrainingSummaryDTO
// @Failure		500	system_error	system	error
// @Router			/v1/train/project/{pid}/training [get]
func (ctl *TrainingController) List(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	v, err := ctl.ts.List(pl.DomainAccount(), ctx.Param("pid"))
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(v))
}

// @Summary		List
// @Description	get trainings
// @Tags			Training
// @Param			pid	path	string	true	"project id"
// @Accept			json
// @Success		200	{object}		app.TrainingSummaryDTO
// @Failure		500	system_error	system	error
// @Router			/v1/train/project/{pid}/training/ws [get]
func (ctl *TrainingController) ListByWS(ctx *gin.Context) {
	pl, csrftoken, _, ok := ctl.checkTokenForWebsocket(ctx, false)
	if !ok {
		return
	}

	pid := ctx.Param("pid")

	// setup websocket
	upgrader := websocket.Upgrader{
		Subprotocols: []string{csrftoken},
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get(headerSecWebsocket) == csrftoken
		},
	}

	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	defer ws.Close()

	ctl.watchTrainings(ws, pl.DomainAccount(), pid)
}

func (ctl *TrainingController) watchTrainings(ws *websocket.Conn, user domain.Account, pid string) {
	finished := func(v []app.TrainingSummaryDTO) (b bool, i int) {
		for i = range v {
			if !v[i].IsDone {
				return
			}
		}

		b = true

		return
	}

	duration := 0
	sleep := func() {
		time.Sleep(time.Second)

		if duration > 0 {
			duration++
		}
	}

	// start loop
	var err error
	var v []app.TrainingSummaryDTO
	var running *app.TrainingSummaryDTO

	start, end := 4, 5
	i := start
	for {
		if i++; i == end {
			v, err = ctl.ts.List(user, pid)
			if err != nil {
				i = start
				sleep()

				continue
			}

			if len(v) == 0 {
				break
			}

			done, index := finished(v)
			if done {
				if err = ws.WriteJSON(newResponseData(v)); err != nil {
					return
				}

				break
			}

			running = &v[index]

			if duration == 0 {
				duration = running.Duration
			} else {
				running.Duration = duration
			}

			if err = ws.WriteJSON(newResponseData(v)); err != nil {
				break
			}

			i = 0
		} else {
			if running.Duration > 0 {
				running.Duration++

				if err = ws.WriteJSON(newResponseData(v)); err != nil {
					break
				}
			}
		}

		sleep()
	}
}

// @Summary		GetLastTrainingConfig
// @Description	get user last preset training config
// @Tags			Training
// @Param			pid	path	string	true	"project id"
// @Accept			json
// @Success		200	{object}		app.TrainingConfigDTO
// @Failure		500	system_error	system	error
// @Router			/v1/train/project/{pid}/config [get]
func (ctl *TrainingController) GetLastTrainingConfig(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	dto, code, err := ctl.ts.GetLastTrainingConfig(
		&app.ResourceIndexCmd{
			Owner: pl.DomainAccount(),
			Id:    ctx.Param("pid"),
		},
	)
	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctx.JSON(http.StatusOK, newResponseData(dto))
	}
}

// @Summary		GetLog
// @Description	get log url of training for downloading
// @Tags			Training
// @Param			pid		path	string	true	"project id"
// @Param			id		path	string	true	"training id"
// @Param			type	path	string	true	"training result: log, output"
// @Accept			json
// @Success		200	{object}		trainingLogResp
// @Failure		500	system_error	system	error
// @Router			/v1/train/project/{pid}/training/{id}/result/{type} [get]
func (ctl *TrainingController) GetResultDownloadURL(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	info := domain.TrainingIndex{
		Project: domain.ResourceIndex{
			Owner: pl.DomainAccount(),
			Id:    ctx.Param("pid"),
		},
		TrainingId: ctx.Param("id"),
	}

	v, code := "", ""
	var err error

	switch ctx.Param("type") {
	case "log":
		v, code, err = ctl.ts.GetLogDownloadURL(&info)

	case "output":
		v, code, err = ctl.ts.GetOutputDownloadURL(&info)

	default:
		ctl.sendBadRequest(ctx, newResponseCodeMsg(
			errorBadRequestParam, "unknown result type",
		))

		return
	}

	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfGet(ctx, trainingLogResp{v})
	}
}

func (ctl *TrainingController) getTrainingInfo(ctx *gin.Context) (domain.TrainingIndex, bool) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return domain.TrainingIndex{}, ok
	}

	return domain.TrainingIndex{
		Project: domain.ResourceIndex{
			Owner: pl.DomainAccount(),
			Id:    ctx.Param("pid"),
		},
		TrainingId: ctx.Param("id"),
	}, true
}
