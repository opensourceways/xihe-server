package routers

import (
	"github.com/opensourceways/xihe-server/controllers"
	"github.com/opensourceways/xihe-server/docs"
	"github.com/opensourceways/xihe-server/models"
	"github.com/opensourceways/xihe-server/util"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

//InitRouter init router
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(util.LoggerToFile())
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = util.GetConfig().AppName
	docs.SwaggerInfo.Description = "set token name: 'Authorization' at header "
	auth := r.Group(docs.SwaggerInfo.BasePath)
	{
		auth.GET("/v1/auth/loginok", controllers.AuthingLoginOk)
		auth.GET("/v1/auth/getDetail/:authingUserId", controllers.AuthingGetUserDetail)
		auth.Use(models.Authorize()) //
		auth.POST("/v1/auth/createUser", controllers.AuthingCreateUser)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
