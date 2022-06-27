package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/xihe-server/util"
)

// @Summary helloworld
// @Description helloworld
// @Tags  V1
// @Accept json
// @Produce json
// @Router /v1/helloworld [get]
func HelloWorld(c *gin.Context) {

	c.JSON(http.StatusOK, util.ExportData(util.CodeStatusNormal, "hello world", nil))
}
