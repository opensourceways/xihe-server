package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/challenge"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForChallengeController(
	rg *gin.RouterGroup,
	crepo repository.Competition,
	qrepo repository.AIQuestion,
	h challenge.Challenge,

) {
	ctl := ChallengeController{
		s: app.NewChallengeService(crepo, qrepo, h, encryptHelper),
	}

	rg.GET("/v1/challenge", ctl.Get)
	rg.GET("/v1/challenge/ranking", ctl.GetRankingList)
	rg.GET("/v1/challenge/aiquestions", ctl.GetAIQuestions)
	rg.POST("/v1/challenge/aiquestions", ctl.Submit)
	rg.POST("/v1/challenge/competitor", ctl.Apply)
}

type ChallengeController struct {
	baseController

	s app.ChallengeService
}

// @Summary		Get
// @Description	get detail of challenge
// @Tags			Challenge
// @Accept			json
// @Success		200	{object}		app.ChallengeCompetitorInfoDTO
// @Failure		500	system_error	system	error
// @Router			/v1/challenge [get]
func (ctl *ChallengeController) Get(ctx *gin.Context) {
	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	var user domain.Account
	if !visitor {
		user = pl.DomainAccount()
	}

	data, err := ctl.s.GetCompetitor(user)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}

// @Summary		GetRankingList
// @Description	get ranking list of challenge
// @Tags			Challenge
// @Accept			json
// @Success		200	{object}		app.ChallengeRankingDTO
// @Failure		500	system_error	system	error
// @Router			/v1/challenge/ranking [get]
func (ctl *ChallengeController) GetRankingList(ctx *gin.Context) {
	data, err := ctl.s.GetRankingList()
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}

// @Summary		Apply
// @Description	apply the challenge
// @Tags			Challenge
// @Param			body	body	competitorApplyRequest	true	"body of applying"
// @Accept			json
// @Success		201
// @Failure		500	system_error	system	error
// @Router			/v1/challenge/competitor [post]
func (ctl *ChallengeController) Apply(ctx *gin.Context) {
	req := competitorApplyRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "apply the challenge")

	if !req.Agreement {
		ctl.sendBadRequestParamWithMsg(ctx, "do not sign the agreement")
		return
	}

	cmd, err := req.toCmd(pl.DomainAccount())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if err := ctl.s.Apply(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusCreated, newResponseData("success"))
	}
}

// @Summary		GetAIQuestions
// @Description	get ai questions
// @Tags			Challenge
// @Accept			json
// @Success		200	{object}		app.AIQuestionDTO
// @Failure		500	system_error	system	error
// @Router			/v1/challenge/aiquestions [get]
func (ctl *ChallengeController) GetAIQuestions(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if data, code, err := ctl.s.GetAIQuestions(pl.DomainAccount()); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}

// @Summary		Submit
// @Description	submit answer of ai question
// @Tags			Challenge
// @Param			body	body	aiQuestionAnswerSubmitRequest	true	"body of ai question answer"
// @Accept			json
// @Success		201	{object}		aiQuestionAnswerSubmitResp
// @Failure		500	system_error	system	error
// @Router			/v1/challenge/aiquestions [post]
func (ctl *ChallengeController) Submit(ctx *gin.Context) {
	req := aiQuestionAnswerSubmitRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "submit answer of ai question")

	score, code, err := ctl.s.SubmitAIQuestionAnswer(pl.DomainAccount(), &cmd)
	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctx.JSON(
			http.StatusCreated,
			newResponseData(aiQuestionAnswerSubmitResp{score}),
		)
	}
}
