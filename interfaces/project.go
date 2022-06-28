package interfaces

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary helloworld
// @Description helloworld
// @Tags  Project
// @Accept json
// @Produce json
// @Router /v1/helloworld [get]
func HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, "hello world")
}
