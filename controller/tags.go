package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func AddRouterForTagsController(
	rg *gin.RouterGroup,
	repo repository.Tags,
) {
	ctl := TagsController{
		s: app.NewTagsService(repo),
	}

	rg.GET("/v1/tags/:type", ctl.List)
}

type TagsController struct {
	baseController

	s app.TagsService
}

// @Title			List
// @Description	list tags
// @Tags			Tags
// @Accept			json
// @Success		200	{object}		app.DomainTagsDTO
// @Failure		500	system_error	system	error
// @Router			/v1/tags/{type} [get]
func (ctl *TagsController) List(ctx *gin.Context) {
	names := apiConfig.Tags.getDomains(ctx.Param("type"))
	if len(names) == 0 {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "unknown type",
		))

		return
	}

	if data, err := ctl.s.List(names); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))
	} else {
		ctx.JSON(http.StatusOK, newResponseData(data))
	}
}
