package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/infrastructure/authing"
)

func AddRouterForGitUserController(
	rg *gin.RouterGroup,
) {
	ctl := GitUserController{}
	rg.Use(authing.Authorize())
	_ = ctl
}

type GitUserController struct {
}
