package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/opensourceways/xihe-server/cloud/app"
	"github.com/opensourceways/xihe-server/cloud/domain"
)

func AddRouterForCloudController(
	rg *gin.RouterGroup,
	s app.CloudService,
) {
	ctl := CloudController{
		s: s,
	}

	rg.GET("/v1/cloud", ctl.List)
	rg.POST("/v1/cloud/subscribe", ctl.Subscribe)
	rg.GET("/v1/cloud/:cid", ctl.Get)
	rg.GET("/v1/cloud/pod/:cid", ctl.GetHttp)
}

type CloudController struct {
	baseController

	s app.CloudService
}

// @Summary List
// @Description list cloud config
// @Tags  Cloud
// @Accept json
// @Success 200 {object} app.UserCompetitionDTO
// @Failure 500 system_error        system error
// @Router /v1/cloud [get]
func (ctl *CloudController) List(ctx *gin.Context) {
	_, _, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	data, err := ctl.s.ListCloud()
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}

// @Summary Subscribe
// @Description subscribe cloud
// @Tags  Cloud
// @Param	body	body	cloudSubscribeRequest	true	"body of subscribe cloud"
// @Accept json
// @Success 201
// @Failure 500 system_error        system error
// @Router /v1/cloud/subscribe [post]
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

	cmd := req.toCmd(pl.DomainAccount())
	if err := cmd.Validate(); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	if code, err := ctl.s.SubscribeCloud(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, "success")
	}
}

// @Summary Get
// @Description get cloud pod
// @Tags  Cloud
// @Param	cid	path 	string			true	"cloud config id"
// @Accept json
// @Success 201 {object} app.InferenceDTO
// @Failure 400 bad_request_body    can't parse request body
// @Failure 500 system_error        system error
// @Router /v1/cloud/{cid} [get]
func (ctl *CloudController) Get(ctx *gin.Context) {
	token := ctx.GetHeader(headerSecWebsocket)
	if token == "" {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "no token",
		))

		return
	}

	pl := oldUserTokenPayload{}
	if ok := ctl.checkApiToken(ctx, token, &pl, false); !ok {
		return
	}

	// setup websocket
	upgrader := websocket.Upgrader{
		Subprotocols: []string{token},
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get(headerSecWebsocket) == token
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
		Pod: domain.Pod{
			Owner:   pl.DomainAccount(),
			CloudId: ctx.Param("cid"),
		},
	}
	if err := cmd.Validate(); err != nil {
		ws.WriteJSON(
			newResponseCodeError(errorBadRequestParam, err),
		)

		log.Errorf("create pod failed: new cmd, err:%s", err.Error())

		return
	}

	for i := 0; i < apiConfig.PodTimeout; i++ {
		dto, err := ctl.s.Get(&cmd)
		if err != nil {
			ws.WriteJSON(newResponseError(err))

			log.Errorf("create pod failed: get status, err:%s", err.Error())

			return
		}

		log.Debugf("info dto:%v", dto)

		if dto.Error != "" || dto.AccessURL != "" {
			ws.WriteJSON(newResponseData(dto))

			log.Debug("create pod done")

			return
		}

		time.Sleep(time.Second)
	}

	log.Error("create pod timeout")

	ws.WriteJSON(newResponseCodeMsg(errorSystemError, "timeout"))
}

// @Summary GetHttp
// @Description get cloud pod
// @Tags  Cloud
// @Param	cid	path 	string			true	"cloud config id"
// @Accept json
// @Success 201 {object} app.InferenceDTO
// @Failure 400 bad_request_body    can't parse request body
// @Failure 500 system_error        system error
// @Router /v1/cloud/{cid} [get]
func (ctl *CloudController) GetHttp(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	cmd := app.PodInfoCmd{
		Pod: domain.Pod{
			Owner:   pl.DomainAccount(),
			CloudId: ctx.Param("cid"),
		},
	}
	if err := cmd.Validate(); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	if dto, err := ctl.s.Get(&cmd); err != nil {
		ctl.sendBadRequestParam(ctx, err)
	} else {
		ctl.sendRespOfGet(ctx, newResponseData(dto))
	}
}
