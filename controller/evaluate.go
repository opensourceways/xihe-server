package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForEvaluateController(
	rg *gin.RouterGroup,
	repo repository.Evaluate,
	train repository.Training,
	sender message.Sender,
) {
	ctl := EvaluateController{
		s: app.NewEvaluateService(
			repo, sender, apiConfig.MinSurvivalTimeOfEvaluate,
		),
		train: train,
	}

	rg.POST("/v1/evaluate/project/:pid/training/:tid/evaluate", checkUserEmailMiddleware(&ctl.baseController), ctl.Create)
	rg.GET("/v1/evaluate/project/:pid/training/:tid/evaluate/:id", ctl.Watch)
}

type EvaluateController struct {
	baseController

	s     app.EvaluateService
	train repository.Training
}

// @Summary Create
// @Description create evaluate
// @Tags  Evaluate
// @Param	pid	path 	string			true	"project id"
// @Param	tid	path 	string			true	"training id"
// @Param	body	body 	EvaluateCreateRequest	true	"body of creating inference"
// @Accept json
// @Success 201 {object} app.EvaluateDTO
// @Failure 400 bad_request_body    can't parse request body
// @Failure 401 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Router /v1/evaluate/project/{pid}/training/{tid}/evaluate [post]
func (ctl *EvaluateController) Create(ctx *gin.Context) {
	req := EvaluateCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	index := domain.TrainingIndex{
		Project: domain.ResourceIndex{
			Owner: pl.DomainAccount(),
			Id:    ctx.Param("pid"),
		},
		TrainingId: ctx.Param("tid"),
	}

	switch req.Type {
	case domain.EvaluateTypeCustom:
		ctl.createCustom(ctx, &index)

	case domain.EvaluateTypeStandard:
		ctl.createStandard(ctx, &index, &req)

	default:
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "unknown evaluate type",
		))

		return
	}
}

func (ctl *EvaluateController) createCustom(ctx *gin.Context, index *domain.TrainingIndex) {
	job, err := ctl.train.GetJob(index)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	dto, err := ctl.s.CreateCustom(&app.CustomEvaluateCreateCmd{
		TrainingIndex: *index,
		AimPath:       job.AimDir,
	})
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(dto))
	}
}

func (ctl *EvaluateController) createStandard(
	ctx *gin.Context, index *domain.TrainingIndex, req *EvaluateCreateRequest,
) {
	job, _, err := ctl.train.GetJobDetail(index)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}
	dto, err := ctl.s.CreateStandard(&app.StandardEvaluateCreateCmd{
		TrainingIndex: *index,
		LogPath:       job.LogPath,

		MomentumScope:     req.MomentumScope,
		BatchSizeScope:    req.BatchSizeScope,
		LearningRateScope: req.LearningRateScope,
	})
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData(dto))

	}
}

// @Summary Watch
// @Description watch evaluate
// @Tags  Evaluate
// @Param	pid	path 	string		true	"project id"
// @Param	tid	path 	string		true	"training id"
// @Param	id	path 	string		true	"evaluate id"
// @Accept json
// @Success 201 {object} app.EvaluateDTO
// @Failure 400 bad_request_body    can't parse request body
// @Failure 401 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Router /v1/evaluate/project/{pid}/training/{tid}/evaluate/{id} [get]
func (ctl *EvaluateController) Watch(ctx *gin.Context) {
	pl, token, ok := ctl.checkTokenForWebsocket(ctx)
	if !ok {
		return
	}

	index := app.EvaluateIndex{
		Id: ctx.Param("id"),
	}
	index.TrainingId = ctx.Param("tid")
	index.Project.Id = ctx.Param("pid")
	index.Project.Owner = pl.DomainAccount()

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

	// start
	for i := 0; i < apiConfig.EvaluateTimeout; i++ {
		dto, err := ctl.s.Get(&index)
		if err != nil {
			ws.WriteJSON(newResponseError(err))

			return
		}

		if dto.Error != "" || dto.AccessURL != "" {
			ws.WriteJSON(newResponseData(dto))

			return
		}

		time.Sleep(time.Second)
	}

	ws.WriteJSON(newResponseCodeMsg(errorSystemError, "timeout"))
}
