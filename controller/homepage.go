package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/app"
	compapp "github.com/opensourceways/xihe-server/competition/app"
	compdomain "github.com/opensourceways/xihe-server/competition/domain"
	courseapp "github.com/opensourceways/xihe-server/course/app"
	coursedomain "github.com/opensourceways/xihe-server/course/domain"
	promapp "github.com/opensourceways/xihe-server/promotion/app"
	spaceapp "github.com/opensourceways/xihe-server/space/app"
)

func AddRouterForHomeController(
	rg *gin.RouterGroup,

	course courseapp.CourseService,
	comp compapp.CompetitionService,
	project spaceapp.ProjectService,
	model app.ModelService,
	dataset app.DatasetService,
	promotion promapp.PromotionService,

) {
	ctl := HomeController{
		course:    course,
		comp:      comp,
		project:   project,
		model:     model,
		dataset:   dataset,
		promotion: promotion,
	}
	rg.GET("/v1/homepage", ctl.ListAll)
	rg.GET("/v1/homepage/:industry", ctl.Get)
}

type HomeController struct {
	baseController

	course    courseapp.CourseService
	comp      compapp.CompetitionService
	project   spaceapp.ProjectService
	model     app.ModelService
	dataset   app.DatasetService
	promotion promapp.PromotionService
}

// @Summary		ListAll
// @Description	list the courses and competitions
// @Tags			HomePage
// @Accept			json
// @Success		200	{object}		homeInfo
// @Failure		500	system_error	system	error
// @Router			/v1/homepage [get]
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

// @Summary		Get
// @Description	get the project dataset model courses and competitions
// @Tags			HomePage
// @Accept			json
// @Success		200	{object}		IndustryDTO
// @Failure		500	system_error	system	error
// @Router			/v1/homepage/{industry} [get]
func (ctl *HomeController) Get(ctx *gin.Context) {
	_, _, ok := ctl.checkUserApiToken(ctx, true)
	if !ok {
		return
	}

	industry := ctx.Param("industry")

	compCmd := compapp.CompetitionListCMD{}
	t, _ := compdomain.NewCompetitionTag(industry)
	compCmd.Tag = t
	compRes, err := ctl.comp.List(&compCmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	}

	ct, _ := coursedomain.NewCourseType(industry)
	courseCmd := courseapp.CourseListCmd{Type: ct}
	courseRes, err := ctl.course.List(&courseCmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	}

	promotionsDTO, err := ctl.promotion.List(
		&promapp.ListPromotionsCmd{
			Tags:     []string{industry},
			PageNo:   1,
			PageSize: 12,
		},
	)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	}

	cmd, err := ctl.getListGlobalResourceParameter(ctx)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	}

	cmd.Tags = append(cmd.Tags, industry)

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

	dto := IndustryDTO{
		Comp:       compRes,
		Course:     courseRes,
		Promotions: promotionsDTO.Items,
		Peoject:    p,
		Model:      m,
		Dataset:    d,
	}

	ctl.sendRespOfGet(ctx, dto)
}
