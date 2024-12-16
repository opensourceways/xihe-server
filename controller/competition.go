package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/competition/app"
	cc "github.com/opensourceways/xihe-server/competition/controller"
	"github.com/opensourceways/xihe-server/competition/domain"
	types "github.com/opensourceways/xihe-server/domain"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	userapp "github.com/opensourceways/xihe-server/user/app"
)

func AddRouterForCompetitionController(
	rg *gin.RouterGroup,
	s app.CompetitionService,
	us userapp.RegService,
	projectPg spacerepo.ProjectPg,
) {
	ctl := CompetitionController{
		s:         s,
		us:        us,
		projectPg: projectPg,
	}

	rg.GET("/v1/competition", ctl.List)
	rg.GET("/v1/competition/:id", ctl.Get)
	rg.GET("/v1/competition/:id/team", ctl.GetMyTeam)
	rg.GET("/v1/competition/:id/ranking", ctl.GetRankingList)
	rg.GET("/v1/competition/:id/submissions", ctl.GetSubmissions)
	rg.GET("/v1/competition/reginfo", ctl.GetRegisterInfo)
	rg.POST("/v1/competition/:id/team", ctl.CreateTeam)
	rg.POST("/v1/competition/:id/submissions", ctl.Submit)
	rg.POST("/v1/competition/:id/competitor", ctl.Apply)
	rg.PUT("/v1/competition/:id/team", ctl.JoinTeam)
	rg.PUT("/v1/competition/:id/related_project", checkUserEmailMiddleware(&ctl.baseController), ctl.AddRelatedProject)
	rg.PUT("/v1/competition/:id/team/action/change_name", ctl.ChangeName)
	rg.PUT("/v1/competition/:id/team/action/transfer_leader", ctl.TransferLeader)
	rg.PUT("/v1/competition/:id/team/action/quit", ctl.QuitTeam)
	rg.PUT("/v1/competition/:id/team/action/delete_member", ctl.DeleteMember)
	rg.PUT("/v1/competition/:id/team/action/dissolve", ctl.Dissolve)
}

type CompetitionController struct {
	baseController

	s         app.CompetitionService
	us        userapp.RegService
	projectPg spacerepo.ProjectPg
}

// @Summary		Apply
// @Description	apply the competition
// @Tags			Competition
// @Param			id		path	string						true	"competition id"
// @Param			body	body	cc.CompetitorApplyRequest	true	"body of applying"
// @Accept			json
// @Success		201
// @Failure		500	system_error	system	error
// @Router			/v1/competition/{id}/competitor [post]
func (ctl *CompetitionController) Apply(ctx *gin.Context) {
	req := cc.CompetitorApplyRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "apply the competition")

	if !req.Agreement {
		ctl.sendBadRequestParamWithMsg(ctx, "do not sign the agreement")
		return
	}

	cmd, err := req.ToCmd(pl.DomainAccount())
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

// @Summary		Get
// @Description	get detail of competition
// @Tags			Competition
// @Param			id	path	string	true	"competition id"
// @Accept			json
// @Success		200	{object}		app.UserCompetitionDTO
// @Failure		500	system_error	system	error
// @Router			/v1/competition/{id} [get]
func (ctl *CompetitionController) Get(ctx *gin.Context) {
	pl, visitor, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	cmd := app.CompetitionGetCmd{}
	var user types.Account
	if !visitor {
		user = pl.DomainAccount()
		cmd.User = user
	}

	var err error
	if str := ctl.getQueryParameter(ctx, "lang"); str != "" {
		cmd.Lang, err = domain.NewLanguage(str)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}
	}

	cmd.CompetitionId = ctx.Param("id")
	data, err := ctl.s.Get(&cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}

// @Summary		List
// @Description	list competitions
// @Tags			Competition
// @Param			status	query	string	false	"competition status, such as over, preparing, in-progress"
// @Param			mine	query	string	false	"just list competitions of competitor, if it is set"
// @Accept			json
// @Success		200	{object}		app.CompetitionSummaryDTO
// @Failure		500	system_error	system	error
// @Router			/v1/competition [get]
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

	if s := ctl.getQueryParameter(ctx, "tag"); s != "" {
		tag, err := domain.NewCompetitionTag(s)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))
			return
		}
		cmd.Tag = tag
	}

	if str := ctl.getQueryParameter(ctx, "lang"); str != "" {
		cmd.Lang, err = domain.NewLanguage(str)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}
	}

	pl, v, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	if ctl.getQueryParameter(ctx, "mine") != "" {
		if v {
			return
		}

		cmd.User = pl.DomainAccount()
	}

	if data, err := ctl.s.List(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}

// @Summary		CreateTeam
// @Description	create team of competition
// @Tags			Competition
// @Param			id		path	string					true	"competition id"
// @Param			body	body	cc.CreateTeamRequest	true	"body of creating team"
// @Accept			json
// @Success		201
// @Failure		500	system_error	system	error
// @Router			/v1/competition/{id}/team [post]
func (ctl *CompetitionController) CreateTeam(ctx *gin.Context) {
	req := cc.CreateTeamRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "create a team of competition")

	cmd, err := req.ToCmd(pl.DomainAccount())
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if code, err := ctl.s.CreateTeam(ctx.Param("id"), &cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, "success")
	}
}

// @Summary		JoinTeam
// @Description	join a team of competition
// @Tags			Competition
// @Param			id		path	string				true	"competition id"
// @Param			body	body	cc.JoinTeamRequest	true	"body of joining team"
// @Accept			json
// @Success		202
// @Failure		500	system_error	system	error
// @Router			/v1/competition/{id}/team [put]
func (ctl *CompetitionController) JoinTeam(ctx *gin.Context) {
	req := cc.JoinTeamRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "join a team of competition")

	cmd, err := req.ToCmd(pl.DomainAccount())
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if code, err := ctl.s.JoinTeam(ctx.Param("id"), &cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPut(ctx, "success")
	}
}

// @Summary		GetMyTeam
// @Description	get team of competition
// @Tags			Competition
// @Param			id	path	string	true	"competition id"
// @Accept			json
// @Success		200	{object}		app.CompetitionTeamDTO
// @Failure		500	system_error	system	error
// @Router			/v1/competition/{id}/team [get]
func (ctl *CompetitionController) GetMyTeam(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	data, code, err := ctl.s.GetMyTeam(ctx.Param("id"), pl.DomainAccount())
	if err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}

// @Summary		GetRankingList
// @Description	get ranking list of competition
// @Tags			Competition
// @Param			id	path	string	true	"competition id"
// @Accept			json
// @Success		200	{object}		app.CompetitionRankingDTO
// @Failure		500	system_error	system	error
// @Router			/v1/competition/{id}/ranking [get]
func (ctl *CompetitionController) GetRankingList(ctx *gin.Context) {
	data, err := ctl.s.GetRankingList(ctx.Param("id"))
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}

// @Summary		GetSubmissions
// @Description	get submissions
// @Tags			Competition
// @Param			id	path	string	true	"competition id"
// @Accept			json
// @Success		200	{object}		app.CompetitionSubmissionsDTO
// @Failure		500	system_error	system	error
// @Router			/v1/competition/{id}/submissions [get]
func (ctl *CompetitionController) GetSubmissions(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	cmd := app.CompetitionGetCmd{}
	cmd.User = pl.DomainAccount()
	cmd.CompetitionId = ctx.Param("id")

	data, err := ctl.s.GetSubmissions(&cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}

// @Summary		Submit
// @Description	submit
// @Tags			Competition
// @Param			id		path		string	true	"competition id"
// @Param			file	formData	file	true	"result file"
// @Accept			json
// @Success		201	{object}		app.CompetitionSubmissionDTO
// @Failure		500	system_error	system	error
// @Router			/v1/competition/{id}/submissions [post]
func (ctl *CompetitionController) Submit(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "submit")

	f, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody, err.Error(),
		))

		return
	}

	if f.Size > apiConfig.MaxCompetitionSubmmitFileSize {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "too big picture",
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

	if err = cmd.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody, err.Error(),
		))
	}

	if v, code, err := ctl.s.Submit(cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, v)
	}
}

// @Summary		AddRelatedProject
// @Description	add related project
// @Tags			Competition
// @Param			id		path	string						true	"competition id"
// @Param			body	body	cc.AddRelatedProjectRequest	true	"project info"
// @Accept			json
// @Success		202
// @Failure		500	system_error	system	error
// @Router			/v1/competition/{id}/related_project [put]
func (ctl *CompetitionController) AddRelatedProject(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "add related project")

	req := cc.AddRelatedProjectRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	owner, name, err := req.ToInfo()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	p, err := ctl.projectPg.GetSummaryByName(owner, name)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	cmd := app.CompetitionAddRelatedProjectCMD{
		Id:      ctx.Param("id"),
		User:    pl.DomainAccount(),
		Project: p,
	}

	if code, err := ctl.s.AddRelatedProject(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPut(ctx, "success")
	}
}

// @Summary		ChangeName
// @Description	change name of a team
// @Tags			Competition
// @Param			id		path	string						true	"competition id"
// @Param			body	body	cc.ChangeTeamNameRequest	true	"body of team name"
// @Accept			json
// @Success		202
// @Failure		500	system_error	system	error
// @Router			/v1/competition/{id}/team/action/change_name [put]
func (ctl *CompetitionController) ChangeName(ctx *gin.Context) {
	req := cc.ChangeTeamNameRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "change a team name")

	cmd, err := req.ToCmd(pl.DomainAccount())
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.s.ChangeTeamName(ctx.Param("id"), &cmd); err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfPut(ctx, "success")
	}
}

// @Summary		TransferLeader
// @Description	transfer leader to a member
// @Tags			Competition
// @Param			id		path	string						true	"competition id"
// @Param			body	body	cc.TransferLeaderRequest	true	"body of member"
// @Accept			json
// @Success		202
// @Failure		500	system_error	system	error
// @Router			/v1/competition/{id}/team/action/transfer_leader [put]
func (ctl *CompetitionController) TransferLeader(ctx *gin.Context) {
	req := cc.TransferLeaderRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "transfer leader to a member")

	cmd, err := req.ToCmd(pl.DomainAccount())
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.s.TransferLeader(ctx.Param("id"), &cmd); err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfPut(ctx, "success")
	}
}

// @Summary		QuitTeam
// @Description	quit team
// @Tags			Competition
// @Param			id	path	string	true	"competition id"
// @Accept			json
// @Success		202
// @Failure		500	system_error	system	error
// @Router			/v1/competition/{id}/team/action/quit [put]
func (ctl *CompetitionController) QuitTeam(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "quit team")

	err := ctl.s.QuitTeam(ctx.Param("id"), pl.DomainAccount())
	if err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfPut(ctx, "success")
	}
}

// @Summary		DeleteMember
// @Description	delete member of a team
// @Tags			Competition
// @Param			id		path	string					true	"competition id"
// @Param			body	body	cc.DeleteMemberRequest	true	"body of delete member"
// @Accept			json
// @Success		202
// @Failure		500	system_error	system	error
// @Router			/v1/competition/{id}/team/action/delete_member [put]
func (ctl *CompetitionController) DeleteMember(ctx *gin.Context) {
	req := cc.DeleteMemberRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctl.sendBadRequestBody(ctx)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "delete member of a team")

	cmd, err := req.ToCmd(pl.DomainAccount())
	if err != nil {
		ctl.sendBadRequestParam(ctx, err)

		return
	}

	if err := ctl.s.DeleteMember(ctx.Param("id"), &cmd); err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfPut(ctx, "success")
	}
}

// @Summary		Dissolve
// @Description	dissolve a team
// @Tags			Competition
// @Param			id	path	string	true	"competition id"
// @Accept			json
// @Success		202
// @Failure		500	system_error	system	error
// @Router			/v1/competition/{id}/team/action/dissolve [put]
func (ctl *CompetitionController) Dissolve(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	prepareOperateLog(ctx, pl.Account, OPERATE_TYPE_USER, "dissolve a team")

	err := ctl.s.DissolveTeam(ctx.Param("id"), pl.DomainAccount())
	if err != nil {
		ctl.sendCodeMessage(ctx, "", err)
	} else {
		ctl.sendRespOfPut(ctx, "success")
	}
}

// @Summary		GetRegisterInfo
// @Description	get register info
// @Tags			Competition
// @Accept			json
// @Success		200	{object}		app.UserRegisterInfoDTO
// @Failure		500	system_error	system	error
// @Router			/v1/competition/reginfo [get]
func (ctl *CompetitionController) GetRegisterInfo(ctx *gin.Context) {

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	if data, err := ctl.us.GetUserRegInfo(pl.DomainAccount()); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}
