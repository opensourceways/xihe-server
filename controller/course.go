package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/course/app"
	"github.com/opensourceways/xihe-server/course/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	commonuser "github.com/opensourceways/xihe-server/user/app"
)

func AddRouterForCourseController(
	rg *gin.RouterGroup,

	s app.CourseService,
	us commonuser.UserService,
	project repository.Project,
	user repository.User,
) {
	ctl := CourseController{
		s:       s,
		us:      us,
		project: project,
		user:    user,
	}

	rg.POST("/v1/course/:id/player", ctl.Apply)
	rg.GET("/v1/course", ctl.List)
	rg.GET("/v1/course/:id", ctl.Get)
	rg.PUT("/v1/course/:id/realted_project", ctl.AddCourseRelatedProject)
	rg.GET("/v1/course/:id/asg/list", ctl.ListAssignments)
	rg.GET("/v1/course/:id/asg/result", ctl.GetSubmissions)
	rg.GET("/v1/course/:id/cert", ctl.GetCertification)
	rg.GET("/v1/course/reginfo", ctl.GetRegisterInfo)
	rg.GET("/v1/course/:id/asg/:asgid", ctl.GetAssignment)
	rg.PUT("/v1/course/:id/record", ctl.AddPlayRecord)
}

type CourseController struct {
	baseController

	s       app.CourseService
	us      commonuser.UserService
	project repository.Project
	user    repository.User
}

// @Summary Apply
// @Description apply the course
// @Tags  Course
// @Param	id	path	string				true	"course id"
// @Param	body	body	StudentApplyRequest	true	"body of applying"
// @Accept json
// @Success 201
// @Failure 500 system_error        system error
// @Router /v1/course/{id}/player [post]
func (ctl *CourseController) Apply(ctx *gin.Context) {
	req := StudentApplyRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	cmd, err := req.toCmd(ctx.Param("id"), pl.DomainAccount())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestParam(err))

		return
	}

	if code, err := ctl.s.Apply(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPost(ctx, "success")
	}
}

// @Summary List
// @Description list the course
// @Tags  Course
// @Param	status	query	string	false	"course status, such as over, preparing, in-progress"
// @Param	type	query	string	false	"course type, such as ai, mindspore, foundation"
// @Param	mine	query	string	false	"just list courses of player, if it is set"
// @Accept json
// @Success 200
// @Failure 500 system_error        system error
// @Router /v1/course [get]
func (ctl *CourseController) List(ctx *gin.Context) {
	var cmd app.CourseListCmd
	var err error

	if str := ctl.getQueryParameter(ctx, "status"); str != "" {
		cmd.Status, err = domain.NewCourseStatus(str)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}
	}

	if str := ctl.getQueryParameter(ctx, "type"); str != "" {
		cmd.Type, err = domain.NewCourseType(str)
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

	if !visitor && ctl.getQueryParameter(ctx, "mine") == "true" {
		cmd.User = pl.DomainAccount()
	}

	if data, err := ctl.s.List(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}

// @Summary Get
// @Description get course infomation
// @Tags  Course
// @Param	id	path	string				true	"course id"
// @Accept json
// @Success 201
// @Failure 500 system_error        system error
// @Router /v1/course/{id} [get]
func (ctl *CourseController) Get(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	cmd := toGetCmd(ctx.Param("id"), pl.DomainAccount())

	if data, err := ctl.s.Get(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}

// @Summary AddCourseRelatedProject
// @Description add related project
// @Tags  Course
// @Param	id	path	string					true	"course id"
// @Param	body	body	AddCourseRelatedProjectRequest	true	"project info"
// @Accept json
// @Success 202
// @Failure 500 system_error        system error
// @Router /v1/course/{id}/realted_project [put]
func (ctl *CourseController) AddCourseRelatedProject(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	req := AddCourseRelatedProjectRequest{}
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
	p, err := ctl.project.GetSummaryByName(owner, name)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	cmd := app.CourseAddReleatedProjectCmd{
		Cid:     ctx.Param("id"),
		User:    pl.DomainAccount(),
		Project: p,
	}

	if code, err := ctl.s.AddReleatedProject(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPut(ctx, "success")
	}
}

// @Summary ListAssignments
// @Description list assignments
// @Tags  Course
// @Param	id	path	string					true	"course id"
// @Param	status	query	string	false	"assignments status, such as finish"
// @Accept json
// @Success 201
// @Failure 500 system_error        system error
// @Router /v1/course/{id}/asg/list [get]
func (ctl *CourseController) ListAssignments(ctx *gin.Context) {

	pl, visitor, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	var cmd app.AsgListCmd
	var err error
	if str := ctl.getQueryParameter(ctx, "status"); str != "" {
		cmd.Status, err = domain.NewWorkStatus(str)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, newResponseCodeError(
				errorBadRequestParam, err,
			))

			return
		}
	}

	if !visitor {
		cmd.User = pl.DomainAccount()
	}
	cmd.Cid = ctx.Param("id")

	if data, err := ctl.s.ListAssignments(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}

// @Summary GetSubmissions
// @Description get submissions
// @Tags  Course
// @Param	id	path	string					true	"course id"
// @Accept json
// @Success 200
// @Failure 500 system_error        system error
// @Router /v1/course/{id}/asg/result [get]
func (ctl *CourseController) GetSubmissions(ctx *gin.Context) {

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	var cmd app.GetSubmissionCmd
	cmd.User = pl.DomainAccount()
	cmd.Cid = ctx.Param("id")

	avatar, err := ctl.user.GetUserAvatarId(pl.DomainAccount())
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
		return
	}

	data, err := ctl.s.GetSubmissions(&cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
		return
	}

	detail := submissionDetail{
		RelateProjectDTO: &data,
	}

	if avatar != nil {
		detail.AvatarId = avatar.AvatarId()
	}

	ctl.sendRespOfGet(ctx, detail)

}

// @Summary GetCertification
// @Description get certification
// @Tags  Course
// @Param	id	path	string					true	"course id"
// @Accept json
// @Success 200
// @Failure 500 system_error        system error
// @Router /v1/course/{id}/cert [get]
func (ctl *CourseController) GetCertification(ctx *gin.Context) {

	pl, visitor, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	var cmd app.CourseGetCmd

	if !visitor {
		cmd.User = pl.DomainAccount()
	}
	cmd.Cid = ctx.Param("id")

	if data, err := ctl.s.GetCertification(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}

// @Summary GetRegisterInfo
// @Description get register info
// @Tags  Course
// @Accept json
// @Success 200
// @Failure 500 system_error        system error
// @Router /v1/course/reginfo [get]
func (ctl *CourseController) GetRegisterInfo(ctx *gin.Context) {

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

// @Summary GetAssignment
// @Description Get assignment
// @Tags  Course
// @Param	id	path	string					true	"course id"
// @Param	asgid	path	string				true	"asg id"
// @Accept json
// @Success 200
// @Failure 500 system_error        system error
// @Router /v1/course/:id/asg/:asgid [get]
func (ctl *CourseController) GetAssignment(ctx *gin.Context) {

	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	var cmd app.AsgGetCmd
	cmd.User = pl.DomainAccount()
	cmd.Cid = ctx.Param("id")
	cmd.AsgId = ctx.Param("asgid")

	if data, err := ctl.s.GetAssignment(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfGet(ctx, data)
	}
}

// @Summary AddPlayRecord
// @Description add play record
// @Tags  Course
// @Param	id	path	string					true	"course id"
// @Param	body	body	PlayRecordRequest	true	"record info"
// @Accept json
// @Success 202
// @Failure 500 system_error        system error
// @Router /v1/course/{id}/record [put]
func (ctl *CourseController) AddPlayRecord(ctx *gin.Context) {
	pl, _, ok := ctl.checkUserApiToken(ctx, false)
	if !ok {
		return
	}

	req := PlayRecordRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestBody,
			"can't fetch request body",
		))

		return
	}

	cmd, err := req.toRecordCmd(ctx.Param("id"), pl.DomainAccount())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))

		return
	}

	if code, err := ctl.s.AddPlayRecord(&cmd); err != nil {
		ctl.sendCodeMessage(ctx, code, err)
	} else {
		ctl.sendRespOfPut(ctx, "success")
	}
}
