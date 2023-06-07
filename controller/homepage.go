package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	compapp "github.com/opensourceways/xihe-server/competition/app"
	compdomain "github.com/opensourceways/xihe-server/competition/domain"
	courseapp "github.com/opensourceways/xihe-server/course/app"
	coursedomain "github.com/opensourceways/xihe-server/course/domain"
)

func AddRouterForHomeController(
	rg *gin.RouterGroup,

	course courseapp.CourseService,
	comp compapp.CompetitionService,
	project app.ProjectService,
	model app.ModelService,
	dataset app.DatasetService,

) {
	ctl := HomeController{
		course:  course,
		comp:    comp,
		project: project,
		model:   model,
		dataset: dataset,
	}
	rg.GET("/v1/homepage", ctl.ListAll)
	rg.GET("/v1/homepage/electricity", ctl.ListAllElectricity)
}

type HomeController struct {
	baseController

	course  courseapp.CourseService
	comp    compapp.CompetitionService
	project app.ProjectService
	model   app.ModelService
	dataset app.DatasetService
}

//	@Summary		ListAll
//	@Description	list the courses and competitions
//	@Tags			HomePage
//	@Accept			json
//	@Success		200	{object}		homeInfo
//	@Failure		500	system_error	system	error
//	@Router			/v1/homepage [get]
func (ctl *HomeController) ListAll(ctx *gin.Context) {
	_, _, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	compCmd := compapp.CompetitionListCMD{}
	compRes, err := ctl.comp.List(&compCmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	}

	courseCmd := courseapp.CourseListCmd{}
	courseRes, err := ctl.course.List(&courseCmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	}

	info := homeInfo{
		Comp:   compRes,
		Course: courseRes,
	}

	ctl.sendRespOfGet(ctx, info)
}

//	@Summary		ListAllElectricity
//	@Description	list the project dataset model courses and competitions
//	@Tags			HomePage
//	@Accept			json
//	@Success		200	{object}		homeElectricityInfo
//	@Failure		500	system_error	system	error
//	@Router			/v1/homepage/electricity [get]
func (ctl *HomeController) ListAllElectricity(ctx *gin.Context) {
	_, _, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	compCmd := compapp.CompetitionListCMD{}
	t, _ := compdomain.NewCompetitionTag("electricity")
	compCmd.Tag = t
	compRes, err := ctl.comp.List(&compCmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	}

	ct, _ := coursedomain.NewCourseType("electricity")
	courseCmd := courseapp.CourseListCmd{Type: ct}
	courseRes, err := ctl.course.List(&courseCmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	}

	cmd := app.GlobalResourceListCmd{}

	cmd.TagKinds = append(cmd.Tags, "electricity")

	p, err := ctl.project.ListGlobal(&cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	}

	m, err := ctl.model.ListGlobal(&cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	}

	d, err := ctl.dataset.ListGlobal(&cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	}

	info := homeElectricityInfo{
		Comp:    compRes,
		Course:  courseRes,
		Peoject: p,
		Model:   m,
		Dataset: d,
	}

	ctl.sendRespOfGet(ctx, info)
}
