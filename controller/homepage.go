package controller

import (
	"github.com/gin-gonic/gin"
	compapp "github.com/opensourceways/xihe-server/competition/app"
	"github.com/opensourceways/xihe-server/course/app"
)

func AddRouterForHomeController(
	rg *gin.RouterGroup,

	course app.CourseService,
	comp compapp.CompetitionService,

) {
	ctl := HomeController{
		course: course,
		comp:   comp,
	}
	rg.GET("/v1/homepage", ctl.ListAll)
}

type HomeController struct {
	baseController

	course app.CourseService
	comp   compapp.CompetitionService
}

// @Summary ListAll
// @Description list the courses and competitions
// @Tags  HomePage
// @Accept json
// @Success 200
// @Failure 500 system_error        system error
// @Router /v1/homepage [get]
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

	courseCmd := app.CourseListCmd{}
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
