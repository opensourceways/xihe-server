package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/opensourceways/xihe-server/aiccfinetune/app"
	"github.com/opensourceways/xihe-server/aiccfinetune/domain"
	appout "github.com/opensourceways/xihe-server/app"
	types "github.com/opensourceways/xihe-server/domain"
	papp "github.com/opensourceways/xihe-server/promotion/app"

	"github.com/opensourceways/xihe-server/utils"
)

func AddRouterForAICCFinetuneController(
	rg *gin.RouterGroup,
	as app.AICCFinetuneService,
	pros papp.PromotionService,
) {
	ctl := AICCFinetuneController{
		as:   as,
		pros: pros,
	}

	rg.POST("/v1/aiccfinetune/:model", checkUserEmailMiddleware(&ctl.baseController), ctl.Create)
	rg.GET("/v1/aiccfinetune/:model", ctl.List)
	rg.GET("/v1/aiccfinetune/:model/ws", ctl.ListByWS)
	rg.GET(
		"/v1/aiccfinetune/:model/:id/result/:type", checkUserEmailMiddleware(&ctl.baseController),
		ctl.GetResultDownloadURL,
	)
	rg.PUT("/v1/aiccfinetune/:model/:id", ctl.Terminate)
	rg.GET("/v1/aiccfinetune/:model/:id", ctl.Get)
	rg.DELETE("/v1/aiccfinetune/:model/:id", ctl.Delete)
	rg.POST("/v1/aiccfinetune/:model/:task/data", ctl.UploadData)
}

type AICCFinetuneController struct {
	baseController

	as   app.AICCFinetuneService
	pros papp.PromotionService
}

// @Summary		Create
// @Description	create aicc finetune
// @Tags			AICC Finetune
// @Param			model	path	string						true	"model name"
// @Param			body	body	aiccFinetuneCreateRequest	true	"body of creating aicc finetune"
// @Accept			json
// @Success		201	{object}			aiccFinetuneCreateResp
// @Failure		400	bad_request_body	can't	parse		request	body
// @Failure		401	bad_request_param	some	parameter	of		body	is	invalid
// @Failure		500	system_error		system	error
// @Router			/v1/aiccfinetune/{model} [post]
func (ctl *AICCFinetuneController) Create(ctx *gin.Context) {
	req := aiccFinetuneCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	// check mindcon apply
	u := pl.DomainAccount()
	if dto, code, err := ctl.pros.GetUserRegisterPromotion(u); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
		return
	} else {
		if !dto[0].IsRegister {
			ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
				errorBadRequestParam, "not registered yet",
			))
			return
		}
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "create aicc finetune")

	cmd := new(app.AICCFinetuneCreateCmd)
	cmd.User = pl.DomainAccount()

	if err := req.toCmd(cmd, ctx.Param("model")); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if err := cmd.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	v, err := ctl.as.Create(cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData(aiccFinetuneCreateResp{v}))
}

// @Summary		Delete
// @Description	delete AICC Finetune
// @Tags			AICC Finetune
// @Param			model	path	string	true	"model name"
// @Param			id		path	string	true	"finetune id"
// @Accept			json
// @Success		201
// @Failure		500	system_error	system	error
// @Router			/v1/aiccfinetune/{model}/{id} [delete]
func (ctl *AICCFinetuneController) Delete(ctx *gin.Context) {
	info, ok := ctl.getAICCFinetuneInfo(ctx)
	if !ok {
		return
	}

	pl, _, _ := ctl.checkUserApiToken(ctx, false)
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "delete aicc finetune")

	if err := ctl.as.Delete(&info); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	utils.DoLog("", info.User.Account(), "delete aicc finetune",
		fmt.Sprintf(" finetune id: %s", info.FinetuneId), "success")

	ctx.JSON(http.StatusNoContent, newResponseData("success"))
}

// @Summary		Terminate
// @Description	terminate AICC Finetune
// @Tags			AICC Finetune
// @Param			model	path	string	true	"model name"
// @Param			id		path	string	true	"finetune id"
// @Accept			json
// @Success		202
// @Failure		500	system_error	system	error
// @Router			/v1/aiccfinetune/{model}/{id} [put]
func (ctl *AICCFinetuneController) Terminate(ctx *gin.Context) {
	info, ok := ctl.getAICCFinetuneInfo(ctx)
	if !ok {
		return
	}

	pl, _, _ := ctl.checkUserApiToken(ctx, false)
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "terminate aicc finetune")

	if err := ctl.as.Terminate(&info); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	utils.DoLog("", info.User.Account(), "terminate finetune",
		fmt.Sprintf("finetune id: %s", info.FinetuneId), "success")

	ctx.JSON(http.StatusAccepted, newResponseData("success"))
}

// @Summary		Get
// @Description	get AICC finetune info
// @Tags			AICC Finetune
//
// @Param			model	path	string	true	"model name"
//
// @Param			id		path	string	true	"finetune id"
// @Accept			json
// @Success		200	{object}		aiccFinetuneDetail
// @Failure		500	system_error	system	error
// @Router			/v1/aiccfinetune/{model}/{id} [get]
func (ctl *AICCFinetuneController) Get(ctx *gin.Context) {
	pl, csrftoken, _, ok := ctl.checkTokenForWebsocket(ctx, false)
	if !ok {
		return
	}

	model, err := domain.NewModelName(ctx.Param("model"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	index := domain.AICCFinetuneIndex{
		User:       pl.DomainAccount(),
		FinetuneId: ctx.Param("id"),
		Model:      model,
	}

	//setup websocket
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

	ctl.watchFinetune(ws, &index)
}

func (ctl *AICCFinetuneController) watchFinetune(ws *websocket.Conn, index *domain.AICCFinetuneIndex) {
	duration := 0
	sleep := func() {
		time.Sleep(time.Second)

		if duration > 0 {
			duration++
		}
	}

	data := &aiccFinetuneDetail{}

	start, end := 4, 5
	i := start
	for {
		if i++; i == end {
			v, code, err := ctl.as.Get(index)
			if err != nil {
				if code == appout.ErrorAICCFinetuneNotFound {
					break
				}

				i = start
				sleep()

				continue
			}

			data.AICCFinetuneDTO = v

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
// @Description	get AICC finetunes
// @Tags			AICC Finetune
// @Param			pid	path	string	true	"project id"
// @Accept			json
// @Success		200	{object}		app.AICCFinetuneSummaryDTO
// @Failure		500	system_error	system	error
// @Router			/v1/aiccfinetune/{model} [get]
func (ctl *AICCFinetuneController) List(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	model, err := domain.NewModelName(ctx.Param("model"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))
		return
	}

	v, err := ctl.as.List(pl.DomainAccount(), model)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(v))
}

// @Summary		List
// @Description	get AICC Finetunes
// @Tags			AICC Finetune
// @Param			pid	path	string	true	"project id"
// @Accept			json
// @Success		200	{object}		app.AICCFinetuneSummaryDTO
// @Failure		500	system_error	system	error
// @Router			/v1/aiccfinetune/{model}/ws [get]
func (ctl *AICCFinetuneController) ListByWS(ctx *gin.Context) {
	pl, csrftoken, _, ok := ctl.checkTokenForWebsocket(ctx, false)
	if !ok {
		return
	}

	model, err := domain.NewModelName(ctx.Param("model"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

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
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	defer ws.Close()

	ctl.watchFinetunes(ws, pl.DomainAccount(), model)
}

func (ctl *AICCFinetuneController) watchFinetunes(ws *websocket.Conn, user types.Account, model domain.ModelName) {
	finished := func(v []app.AICCFinetuneSummaryDTO) (b bool, i int) {
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
	var v []app.AICCFinetuneSummaryDTO
	var running *app.AICCFinetuneSummaryDTO

	start, end := 4, 5
	i := start
	for {
		if i++; i == end {
			v, err = ctl.as.List(user, model)
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
					log.Error("create pod timeout | web socket write error:", err)
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

// @Summary		GetLog
// @Description	get log url of AICC Finetune for downloading
// @Tags			AICC Finetune
// @Param			pid		path	string	true	"project id"
// @Param			id		path	string	true	"finetune id"
// @Param			type	path	string	true	"aicc finetune result: log, output"
// @Accept			json
// @Success		200	{object}		aiccFinetuneLogResp
// @Failure		500	system_error	system	error
// @Router			/v1/aiccfinetune/{model}/{id}/result/{type} [get]
func (ctl *AICCFinetuneController) GetResultDownloadURL(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	var err error
	model, err := domain.NewModelName(ctx.Param("model"))
	if err != nil {
		ctl.sendBadRequest(ctx, newResponseCodeMsg(
			errorBadRequestParam, "unknown model name",
		))

	}

	info := domain.AICCFinetuneIndex{
		User:       pl.DomainAccount(),
		Model:      model,
		FinetuneId: ctx.Param("id"),
	}

	v, code := "", ""

	switch ctx.Param("type") {
	case "log":
		v, code, err = ctl.as.GetLogDownloadURL(&info)

	case "output":
		v, code, err = ctl.as.GetOutputDownloadURL(&info)

	default:
		ctl.sendBadRequest(ctx, newResponseCodeMsg(
			errorBadRequestParam, "unknown result type",
		))

		return
	}

	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfGet(ctx, aiccFinetuneLogResp{v})
	}
}

func (ctl *AICCFinetuneController) getAICCFinetuneInfo(ctx *gin.Context) (domain.AICCFinetuneIndex, bool) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return domain.AICCFinetuneIndex{}, ok
	}

	model, err := domain.NewModelName(ctx.Param("model"))
	if err != nil {
		return domain.AICCFinetuneIndex{}, false
	}

	return domain.AICCFinetuneIndex{
		User:       pl.DomainAccount(),
		FinetuneId: ctx.Param("id"),
		Model:      model,
	}, true
}

// @Summary		UploadData
// @Description	Upload Data
// @Tags			AICC Finetune
// @Param			model	path		string	true	"model name"
// @Param			file	formData	file	true	"result file"
// @Accept			json
// @Success		201	{object}		app.UploadDataDTO
// @Failure		500	system_error	system	error
// @Router			/v1/aiccfinetune/{model}/{task}/data [post]
func (ctl *AICCFinetuneController) UploadData(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "upload data")

	model, err := domain.NewModelName(ctx.Param("model"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	task, err := domain.NewFinetuneTask(ctx.Param("task"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	f, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody, err.Error(),
		))

		return
	}

	if f.Size > apiConfig.MaxFinetuneSubmmitFileSize {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "too big file",
		))

		return
	}

	p, err := f.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "can't get file",
		))

		return
	}

	defer p.Close()

	err = domain.NewFileName(f.Filename, model.ModelName(), task.FinetuneTask())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	cmd := &app.UploadDataCmd{
		FileName: f.Filename,
		Data:     p,
		User:     pl.DomainAccount(),
		Model:    model,
		Task:     task,
	}

	if v, err := ctl.as.UploadData(cmd); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	} else {
		ctl.sendRespOfPost(ctx, v)
	}
}
