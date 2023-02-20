package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/competition/app"
	"github.com/opensourceways/xihe-server/competition/domain"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForCompetitionController(
	rg *gin.RouterGroup,
	s app.CompetitionService,
	project repository.Project,
) {
	ctl := CompetitionController{
		s:       s,
		project: project,
	}

	rg.GET("/v1/competition", ctl.List)
	rg.GET("/v1/competition/:id", ctl.Get)
	rg.GET("/v1/competition/:id/team", ctl.GetTeam)
	rg.GET("/v1/competition/:id/ranking", ctl.GetRankingList)
	rg.GET("/v1/competition/:id/submissions", ctl.GetSubmissions)
	rg.POST("/v1/competition/:id/submissions", ctl.Submit)
	rg.POST("/v1/competition/:id/competitor", ctl.Apply)
	rg.PUT("/v1/competition/:id/realted_project", ctl.AddRelatedProject)
}

type CompetitionController struct {
	baseController

	s       app.CompetitionService
	project repository.Project
}

// @Summary Apply
// @Description apply the competition
// @Tags  Competition
// @Param	body	body	competitorApplyRequest	true	"body of applying"
// @Accept json
// @Success 201
// @Failure 500 system_error        system error
// @Router /v1/competition/{id}/competitor [post]
func (ctl *CompetitionController) Apply(ctx *gin.Context) {
	req := competitorApplyRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	cmd, err := req.toCmd(pl.DomainAccount())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestParam(err))

		return
	}

	if code, err := ctl.s.Apply(ctx.Param("id"), &cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, "success")
	}
}

// @Summary Get
// @Description get detail of competition
// @Tags  Competition
// @Param	id	path	string	true	"competition id"
// @Accept json
// @Success 200 {object} app.UserCompetitionDTO
// @Failure 500 system_error        system error
// @Router /v1/competition/{id} [get]
func (ctl *CompetitionController) Get(ctx *gin.Context) {
	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	var user types.Account
	if !visitor {
		user = pl.DomainAccount()
	}

	data, err := ctl.s.Get(ctx.Param("id"), user)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}

// @Summary List
// @Description list competitions
// @Tags  Competition
// @Param	status	query	string	false	"competition status, such as done, preparing, in-progress"
// @Param	mine	query	string	false	"just list competitions of competitor, if it is set"
// @Accept json
// @Success 200 {object} app.CompetitionSummaryDTO
// @Failure 500 system_error        system error
// @Router /v1/competition [get]
func (ctl *CompetitionController) List(ctx *gin.Context) {
	cmd := app.CompetitionListCMD{}
	var err error

	if str := ctl.getQueryParameter(ctx, "status"); str != "" {
		cmd.Status, err = domain.NewCompetitionStatus(str)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}
	}

	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	if !visitor && ctl.getQueryParameter(ctx, "mine") != "" {
		cmd.User = pl.DomainAccount()
	}

	if data, err := ctl.s.List(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}

// @Summary GetTeam
// @Description get team of competition
// @Tags  Competition
// @Param	id	path	string	true	"competition id"
// @Accept json
// @Success 200 {object} app.CompetitionTeamDTO
// @Failure 500 system_error        system error
// @Router /v1/competition/{id}/team [get]
func (ctl *CompetitionController) GetTeam(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	data, code, err := ctl.s.GetTeam(ctx.Param("id"), pl.DomainAccount())
	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}

// @Summary GetRankingList
// @Description get ranking list of competition
// @Tags  Competition
// @Param	id	path	string	true	"competition id"
// @Accept json
// @Success 200 {object} app.CompetitonRankingDTO
// @Failure 500 system_error        system error
// @Router /v1/competition/{id}/ranking [get]
func (ctl *CompetitionController) GetRankingList(ctx *gin.Context) {
	data, err := ctl.s.GetRankingList(ctx.Param("id"))
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}

// @Summary GetSubmissions
// @Description get submissions
// @Tags  Competition
// @Param	id	path	string	true	"competition id"
// @Accept json
// @Success 200 {object} app.CompetitionSubmissionsDTO
// @Failure 500 system_error        system error
// @Router /v1/competition/{id}/submissions [get]
func (ctl *CompetitionController) GetSubmissions(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	data, err := ctl.s.GetSubmissions(ctx.Param("id"), pl.DomainAccount())
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}

// @Summary Submit
// @Description submit
// @Tags  Competition
// @Param	id	path		string	true	"competition id"
// @Param	file	formData	file	true	"result file"
// @Accept json
// @Success 201 {object} app.CompetitionSubmissionDTO
// @Failure 500 system_error        system error
// @Router /v1/competition/{id}/submissions [post]
func (ctl *CompetitionController) Submit(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	f, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody, err.Error(),
		))

		return
	}

	p, err := f.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "can't get picture",
		))

		return
	}

	defer p.Close()

	cmd := &app.CompetitionSubmitCMD{
		CompetitionId: ctx.Param("id"),
		FileName:      f.Filename,
		Data:          p,
		User:          pl.DomainAccount(),
	}
	if v, code, err := ctl.s.Submit(cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, v)
	}
}

// @Summary AddRelatedProject
// @Description add related project
// @Tags  Competition
// @Param	id	path	string					true	"competition id"
// @Param	body	body	competitionAddRelatedProjectRequest	true	"project info"
// @Accept json
// @Success 202
// @Failure 500 system_error        system error
// @Router /v1/competition/{id}/realted_project [put]
func (ctl *CompetitionController) AddRelatedProject(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	req := competitionAddRelatedProjectRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	owner, name, err := req.toInfo()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	p, err := ctl.project.GetSummaryByName(owner, name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	cmd := app.CompetitionAddReleatedProjectCMD{
		Id:      ctx.Param("id"),
		User:    pl.DomainAccount(),
		Project: p,
	}

	if err = ctl.s.AddRelatedProject(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfPut(ctx, "success")
	}
}
