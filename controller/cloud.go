package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/cloud/app"
	"github.com/opensourceways/xihe-server/domain"
	userapp "github.com/opensourceways/xihe-server/user/app"
	"github.com/opensourceways/xihe-server/utils"
)

func AddRouterForCloudController(
	rg *gin.RouterGroup,
	s app.CloudService,
	us userapp.WhiteListService,
) {
	ctl := CloudController{
		s:  s,
		us: us,
	}

	rg.GET("/v1/cloud", ctl.List)
	rg.POST("/v1/cloud/subscribe", checkUserEmailMiddleware(&ctl.baseController), ctl.Subscribe)
	rg.GET("/v1/cloud/:cid", ctl.Get)
	rg.GET("/v1/cloud/pod/:cid", ctl.GetHttp)
	rg.GET("/v1/cloud/read/:owner", ctl.CanRead)
	rg.DELETE("/v1/cloud/pod/:id", ctl.ReleasePod)
	rg.GET("/v1/ws/cloud/pod/:id", ctl.WsSendReleasedPod)
}

type CloudController struct {
	baseController

	s  app.CloudService
	us userapp.WhiteListService
}

// @Summary		List
// @Description	list cloud config
// @Tags			Cloud
// @Accept			json
// @Success		200	{object}		[]app.CloudDTO
// @Failure		500	system_error	system	error
// @Router			/v1/cloud [get]
func (ctl *CloudController) List(ctx *gin.Context) {
	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	cmd := new(app.GetCloudConfCmd)
	if visitor {
		cmd.ToCmd(nil, visitor)
	} else {
		cmd.ToCmd(pl.DomainAccount(), visitor)
	}

	data, err := ctl.s.ListCloud(cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}

// @Summary		Subscribe
// @Description	subscribe cloud
// @Tags			Cloud
// @Param			body	body	cloudSubscribeRequest	true	"body of subscribe cloud"
// @Accept			json
// @Success		201
// @Failure		500	system_error	system	error
// @Router			/v1/cloud/subscribe [post]
func (ctl *CloudController) Subscribe(ctx *gin.Context) {
	req := cloudSubscribeRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "subscribe cloud")

	cmd := req.toCmd(pl.DomainAccount())
	if err := cmd.Validate(); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	if code, err := ctl.s.SubscribeCloud(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		utils.DoLog("", pl.Account, "create jupyter", cmd.CloudId, "success")

		ctl.sendRespOfPost(ctx, "success")
	}
}

// @Summary		Get
// @Description	get cloud pod
// @Tags			Cloud
// @Param			cid	path	string	true	"cloud config id"
// @Accept			json
// @Success		201	{object}			app.InferenceDTO
// @Failure		400	bad_request_body	can't	parse	request	body
// @Failure		500	system_error		system	error
// @Router			/v1/cloud/{cid} [get]
func (ctl *CloudController) Get(ctx *gin.Context) {
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

	cmd := app.PodInfoCmd{
		User:    pl.DomainAccount(),
		CloudId: ctx.Param("cid"),
	}
	if err := cmd.Validate(); err != nil {
		if wsErr := ws.WriteJSON(newResponseCodeError(errorBadRequestParam, err)); wsErr != nil {
			log.Errorf("create pod failed | web socket err:%s", wsErr.Error())
		}

		log.Errorf("create pod failed: new cmd, err:%s", err.Error())

		return
	}

	for i := 0; i < apiConfig.PodTimeout; i++ {
		dto, err := ctl.s.Get(&cmd)
		if err != nil {
			if wsErr := ws.WriteJSON(newResponseError(err)); wsErr != nil {
				log.Errorf("create pod failed: web socket write err:%s", wsErr.Error())
			}

			log.Errorf("create pod failed: get status, err:%s", err.Error())

			return
		}

		log.Debugf("info dto:%v", dto)

		if dto.Error != "" || dto.AccessURL != "" {
			if wsErr := ws.WriteJSON(newResponseData(dto)); wsErr != nil {
				log.Errorf("create pod failed: web socket write err:%s", wsErr.Error())
			}

			log.Debug("create pod done")

			return
		}

		time.Sleep(time.Second)
	}

	log.Error("create pod timeout")

	if wsErr := ws.WriteJSON(newResponseCodeMsg(errorSystemError, "timeout")); wsErr != nil {
		log.Errorf("create pod timeout | web socket write error:%s", wsErr.Error())
	}
}

// @Summary		GetHttp
// @Description	get cloud pod
// @Tags			Cloud
// @Param			cid	path	string	true	"cloud config id"
// @Accept			json
// @Success		201	{object}			app.InferenceDTO
// @Failure		400	bad_request_body	can't	parse	request	body
// @Failure		500	system_error		system	error
// @Router			/v1/cloud/{cid} [get]
func (ctl *CloudController) GetHttp(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	cmd := app.PodInfoCmd{
		User:    pl.DomainAccount(),
		CloudId: ctx.Param("cid"),
	}
	if err := cmd.Validate(); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	if dto, err := ctl.s.Get(&cmd); err != nil {
		ctl.sendBadRequestParam(ctx, err)
	} else {
		ctl.sendRespOfGet(ctx, dto)
	}
}

// @Summary		CanRead
// @Description	get cloud pod
// @Tags			Cloud
// @Param			owner	path	string	true	""
// @Accept			json
// @Success		201	{object}			app.InferenceDTO
// @Failure		400	bad_request	not allowed
// @Failure		500	system_error		system	error
// @Router			/v1/cloud/read/{owner} [get]
func (ctl *CloudController) CanRead(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	if pl.Account == "" {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, errors.New("not identified"),
		))

		return
	}

	owner, err := domain.NewAccount(ctx.Param("owner"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if pl.isNotMe(owner) {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorNotAllowed, "not allowed",
		))

		return
	}
	type Res struct {
		Status string `json:"status"`
	}

	res := Res{
		Status: "success",
	}

	ctl.sendRespOfGet(ctx, res)
}

// @Summary		Release
// @Description	release cloud resource
// @Tags			Cloud
// @Param			id	path	string	true	""
// @Param			cloud_id	query	string	true	""
// @Accept			json
// @Success		204
// @Failure		404	{string} string "not found"
// @Failure		500	{object} responseData "system error"
// @Router			/v1/cloud/pod/{id} [delete]
func (ctl *CloudController) ReleasePod(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "release cloud")

	cmd := &app.ReleaseCloudCmd{
		PodId: ctx.Param("id"),
		User:  pl.DomainAccount(),
	}

	if err := cmd.Validate(); err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.s.ReleaseCloud(cmd); err != nil {
		if errors.Is(err, app.ErrCloudReleased) {
			ctx.JSON(http.StatusNotFound, newResponseError(err))
		} else if errors.Is(err, app.ErrCloudNotAllowed) {
			ctx.JSON(http.StatusForbidden, newResponseError(err))
		} else {
			logrus.Errorf("fail to release cloud, pod id: %s, err: %s", cmd.PodId, err.Error())

			ctl.sendRespWithInternalError(ctx, newResponseError(ErrSystemFault))
		}

		return
	}

	ctl.sendRespOfDelete(ctx)
}

// @Summary		WsSendReleasedPod
// @Description	WsSendReleasedPod is a websocket api which sends a released pod data to client.
// @Tags			Cloud
// @Param			id	    path	string	true	"pod id"
// @Accept			json
// @Success		200 {object} app.PodInfoDTO
// @Failure		404	{string} string "not found"
// @Failure		500	{object} responseData "system error"
// @Router			/v1/cloud/pod/{id} [get]
func (ctl *CloudController) WsSendReleasedPod(ctx *gin.Context) {
	_, csrftoken, _, ok := ctl.checkTokenForWebsocket(ctx, false)
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

	cmd := app.GetReleasedPodCmd{
		PodId: ctx.Param("id"),
	}

	for i := 0; i < apiConfig.PodTimeout; i++ {
		dto, err := ctl.s.GetReleasedPod(&cmd)
		if errors.Is(err, app.ErrPodNotFound) {
			time.Sleep(time.Second)
			continue
		} else if err != nil {
			log.Errorf("[RELEASE] fail to get pod %s, err:%s", cmd.PodId, err.Error())

			if wsErr := ws.WriteJSON(newResponseError(err)); wsErr != nil {
				log.Errorf("[RELEASE] fail to get pod | web socket write error:%s", wsErr.Error())
			}

			return
		}

		if wsErr := ws.WriteJSON(newResponseData(dto)); wsErr != nil {
			log.Errorf("[RELEASE] fail to get pod | web socket write error:%s", wsErr.Error())
		}

		return
	}

	log.Errorf("release pod %s timeout", cmd.PodId)

	if wsErr := ws.WriteJSON(newResponseCodeMsg(errorSystemError, "timeout")); wsErr != nil {
		log.Errorf("[RELEASE] fail to get pod | web socket write error:%s", wsErr.Error())
	}
}
