package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

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
		log.Error("create pod timeout | web socket write error:", wsErr)
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
