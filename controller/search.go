package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain/repository"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

func AddRouterForSearchController(
	rg *gin.RouterGroup,
	user userrepo.User,
	proj repository.Project,
	model repository.Model,
	dataset repository.Dataset,
) {
	ctl := SearchController{
		s: app.NewSearchService(user, model, proj, dataset),
	}

	rg.GET("/v1/search", ctl.List)
}

type SearchController struct {
	baseController

	s app.SearchService
}

//	@Title			Search
//	@Description	search resource and user
//	@Tags			Search
//	@Param			name	query	string	true	"name of resource or user"
//	@Accept			json
//	@Success		200	{object}		app.SearchDTO
//	@Failure		500	system_error	system	error
//	@Router			/v1/search [get]
func (ctl *SearchController) List(ctx *gin.Context) {
	name := ctl.getQueryParameter(ctx, "name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, newResponseCodeMsg(
			errorBadRequestParam, "no search object",
		))

	}

	name = utils.XSSFilter(name)

	data := ctl.s.Search(name)
	ctx.JSON(http.StatusOK, newResponseData(data))
}
