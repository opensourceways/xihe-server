package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/community-robot-lib/interrupts"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/controller"
	"github.com/opensourceways/xihe-server/docs"
	"github.com/opensourceways/xihe-server/infrastructure/authing"
	"github.com/opensourceways/xihe-server/infrastructure/mongodb"
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func StartWebServer(port int, timeout time.Duration, cfg *config.Config) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logRequest())

	setRouter(r, cfg)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: r,
	}

	defer interrupts.WaitForGracefulShutdown()
	log.Println("server start at :" + srv.Addr)
	interrupts.ListenAndServe(srv, timeout)
}

//setRouter init router
func setRouter(engine *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "xihe"
	docs.SwaggerInfo.Description = "set token name: 'Authorization' at header "

	v1 := engine.Group(docs.SwaggerInfo.BasePath)
	{
		controller.AddRouterForProjectController(
			v1,
			mongodb.NewProjectRepository(cfg.Mongodb.ProjectCollection),
		)

		controller.AddRouterForUserController(
			v1,
			repositories.NewUserRepository(authing.NewUserMapper()),
		)
		controller.AddRouterForGitGroupController(
			v1,
		)
		controller.AddRouterForGitProjectController(
			v1,
		)
		controller.AddRouterForGitUserController(
			v1,
		)
	}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()

		logrus.Infof(
			"| %d | %d | %s | %s |",
			c.Writer.Status(),
			endTime.Sub(startTime),
			c.Request.Method,
			c.Request.RequestURI,
		)
	}
}
