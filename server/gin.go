package server

import (
	"fmt"

	"github.com/opensourceways/xihe-server/docs"
	"github.com/opensourceways/xihe-server/interfaces"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func StartWebServer() {
	r := setRouter()
	address := fmt.Sprintf(":%d", 8080)
	if err := r.Run(address); err != nil {
		fmt.Errorf("startup meta  http service failed, err:%v\n", err)
	}
}

//setRouter init router
func setRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "xihe"
	docs.SwaggerInfo.Description = "set token name: 'Authorization' at header "

	v1 := r.Group(docs.SwaggerInfo.BasePath)
	{
		v1.GET("/v1/helloworld", interfaces.HelloWorld)

	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
