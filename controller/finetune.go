package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/opensourceways/community-robot-lib/utils"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/finetune"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForFinetuneController(
	rg *gin.RouterGroup,
	fs finetune.Finetune,
	repo repository.Finetune,
	sender message.Sender,
) {
	ctl := FinetuneController{
		fs: app.NewFinetuneService(
			fs, repo, sender,
		),
	}

	rg.POST("/v1/finetune", ctl.Create)
	rg.GET("/v1/finetune", ctl.List)
	rg.GET("/v1/finetune/ws", ctl.WatchFinetunes)
	rg.GET("/v1/finetune/:id/log", ctl.Log)
	rg.GET("/v1/finetune/:id/log/ws", ctl.WatchSingle)
	rg.PUT("/v1/finetune/:id", ctl.Terminate)
	rg.DELETE("v1/finetune/:id", ctl.Delete)
}

type FinetuneController struct {
	baseController

	fs app.FinetuneService
}

// @Summary		Create
// @Description	create finetune
// @Tags			Finetune
// @Param			body	body	FinetuneCreateRequest	true	"body of creating finetune"
// @Accept			json
// @Success		201	{object}		finetuneCreateResp
// @Failure		500	system_error	system	error
// @Router			/v1/finetune [post]
func (ctl *FinetuneController) Create(ctx *gin.Context) {
	req := FinetuneCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		prepareOperateLog(ctx, "anonymous", OPERATE_TYPE_SYSTEM, "create finetune")
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "create finetune")

	cmd, err := req.toCmd(pl.DomainAccount())
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if v, code, err := ctl.fs.Create(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, finetuneCreateResp{v})
	}
}

// @Summary		Delete
// @Description	delete finetune
// @Tags			Finetune
// @Param			id	path	string	true	"finetune id"
// @Accept			json
// @Success		204
// @Failure		500	system_error	system	error
// @Router			/v1/finetune/{id} [delete]
func (ctl *FinetuneController) Delete(ctx *gin.Context) {
	index, ok := ctl.finetuneIndex(ctx)
	if !ok {
		prepareOperateLog(ctx, "anonymous", OPERATE_TYPE_SYSTEM, "delete finetune")
		return
	}

	pl, _, _ := ctl.checkUserApiToken(ctx, false)
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "delete finetune")

	if err := ctl.fs.Delete(&index); err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfDelete(ctx)
	}
}

// @Summary		Terminate
// @Description	terminate finetune
// @Tags			Finetune
// @Param			id	path	string	true	"finetune id"
// @Accept			json
// @Success		202
// @Failure		500	system_error	system	error
// @Router			/v1/finetune/{id} [put]
func (ctl *FinetuneController) Terminate(ctx *gin.Context) {
	index, ok := ctl.finetuneIndex(ctx)
	if !ok {
		prepareOperateLog(ctx, "anonymous", OPERATE_TYPE_SYSTEM, "terminate finetune")
		return
	}

	pl, _, _ := ctl.checkUserApiToken(ctx, false)
	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "terminate finetune")

	if err := ctl.fs.Terminate(&index); err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfPut(ctx, "success")
	}
}

func (ctl *FinetuneController) finetuneIndex(ctx *gin.Context) (
	index domain.FinetuneIndex, ok bool,
) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if ok {
		index.Owner = pl.DomainAccount()
		index.Id = ctx.Param("id")
	}

	return
}

// @Summary		List
// @Description	list finetunes
// @Tags			Finetune
// @Accept			json
// @Success		200	{object}		app.UserFinetunesDTO
// @Failure		500	system_error	system	error
// @Router			/v1/finetune [get]
func (ctl *FinetuneController) List(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if v, code, err := ctl.fs.List(pl.DomainAccount()); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfGet(ctx, v)
	}
}

// @Summary		WatchFinetunes
// @Description	watch finetunes
// @Tags			Finetune
// @Accept			json
// @Success		200	{object}		app.FinetuneSummaryDTO
// @Failure		500	system_error	system	error
// @Router			/v1/finetune/ws [get]
func (ctl *FinetuneController) WatchFinetunes(ctx *gin.Context) {
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
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	defer ws.Close()

	ctl.watchFinetunes(ws, pl.DomainAccount())
}

func (ctl *FinetuneController) watchFinetunes(ws *websocket.Conn, user domain.Account) {
	finished := func(v []app.FinetuneSummaryDTO) (b bool, i int) {
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
	var v []app.FinetuneSummaryDTO
	var running *app.FinetuneSummaryDTO

	start, end := 4, 5
	i := start
	for {
		if i++; i == end {
			dto, code, err := ctl.fs.List(user)
			if err != nil {
				if code == app.ErrorFinetuneNoPermission {
					break
				}

				i = start
				sleep()

				continue
			}

			if len(dto.Datas) == 0 {
				break
			}
			v = dto.Datas

			done, index := finished(v)
			if done {
				ws.WriteJSON(newResponseData(v))

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

				if err := ws.WriteJSON(newResponseData(v)); err != nil {
					break
				}
			}
		}

		sleep()
	}
}

// @Summary		WatchSingle
// @Description	watch single finetune
// @Tags			Finetune
// @Param			id	path	string	true	"finetune id"
// @Accept			json
// @Success		200	{object}		finetuneLog
// @Failure		500	system_error	system	error
// @Router			/v1/finetune/{id}/log/ws [get]
func (ctl *FinetuneController) WatchSingle(ctx *gin.Context) {
	pl, csrftoken, _, ok := ctl.checkTokenForWebsocket(ctx, false)
	if !ok {
		return
	}

	index := domain.FinetuneIndex{
		Owner: pl.DomainAccount(),
		Id:    ctx.Param("id"),
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

	for {
		v, code, err := ctl.fs.GetJobInfo(&index)
		if err != nil {
			logrus.Errorf(
				"get finetune job failed, code=%s, err=%s",
				code, err.Error(),
			)
			if code == app.ErrorFinetuneNotFound {
				break
			}

			time.Sleep(time.Second)

			continue
		}

		content, err := downloadLog(v.LogPreviewURL)
		if err != nil {
			logrus.Errorf(
				"download finetune log failed, err=%s", err.Error(),
			)
			time.Sleep(time.Second)

			continue
		}

		if len(content) > 0 {
			data := newResponseData(finetuneLog{string(content)})

			if err = ws.WriteJSON(data); err != nil {
				break
			}
		}

		if v.IsDone {
			break
		}

		time.Sleep(5 * time.Second)
	}
}

// @Summary		Log
// @Description	download finetune log
// @Tags			Finetune
// @Param			id	path	string	true	"finetune id"
// @Accept			json
// @Success		200	{object}		finetuneLog
// @Failure		500	system_error	system	error
// @Router			/v1/finetune/{id}/log [get]
func (ctl *FinetuneController) Log(ctx *gin.Context) {
	index, ok := ctl.finetuneIndex(ctx)
	if !ok {
		return
	}

	v, code, err := ctl.fs.GetJobInfo(&index)
	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)

		return
	}

	content, err := downloadLog(v.LogPreviewURL)
	if err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfGet(ctx, finetuneLog{string(content)})
	}
}

func downloadLog(link string) ([]byte, error) {
	if link == "" {
		return nil, nil
	}

	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, err
	}

	cli := utils.NewHttpClient(3)
	v, _, err := cli.Download(req)

	return v, err
}
