package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/course/app"
)

func AddRouterForCourseController(
	rg *gin.RouterGroup,
	s app.CourseService,
) {
	ctl := CourseController{
		s: s,
	}

	rg.POST("/v1/course/:id/player", ctl.Apply)
}

type CourseController struct {
	baseController

	s app.CourseService
}

// @Summary Apply
// @Description apply the competition
// @Tags  Competition
// @Param	id	path	string				true	"competition id"
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

	if err := ctl.s.Apply(&cmd); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctl.sendRespOfPost(ctx, "success")
	}
}
