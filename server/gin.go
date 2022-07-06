package server

import (
	"fmt"

	"github.com/opensourceways/xihe-server/docs"
	"github.com/opensourceways/xihe-server/infrastructure"
	"github.com/opensourceways/xihe-server/interfaces"
	"github.com/opensourceways/xihe-server/util"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func StartWebServer() {
	r := setRouter()
	address := fmt.Sprintf(":%d", util.GetConfig().AppPort)
	util.Log.Infof(
		" startup meta http service at port %s .and %s mode \n",
		address, util.GetConfig().AppModel,
	)
	if err := r.Run(address); err != nil {
		util.Log.Fatalf("startup meta  http service failed, err:%v\n", err.Error())
	}
}

//setRouter init router
func setRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(util.LoggerToFile())
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = util.GetConfig().AppName
	docs.SwaggerInfo.Description = "请在请求的header中加入 'Authorization'  "
	repo, err := infrastructure.NewRepositories()
	if err != nil {
		util.Log.Fatalf(" start repository failed, err:%v\n", err.Error())

		return nil
	}
	projectController := interfaces.NewProject(repo.ProjectRepo)
	userController := interfaces.NewUser(repo.UserRepo)

	v1 := r.Group(docs.SwaggerInfo.BasePath + "/v1")
	{
		user := v1.Group("/user")
		{
			user.GET("/checkLogin", userController.CheckLogin)
			user.GET("/callback", userController.AuthingCallback)
			user.GET("/findUser", userController.FindUser)
			user.Use(infrastructure.Authorize())
			user.GET("/getCurrentUser", userController.GetCurrentUser)
			user.PUT("/updatePhone", userController.UpdatePhone)
			user.GET("/sendSmsCode", userController.SendSmsCode)
			user.PUT("/bindPhone", userController.BindPhone)
			user.GET("/sendEmailToResetPswd", userController.SendEmailToResetPswd)
			user.GET("/sendEmailToVerifyEmail", userController.SendEmailToVerifyEmail)
			user.PUT("/resetPasswordByEmailCode", userController.ResetPasswordByEmailCode)
			user.PUT("/updateProfile/:id", userController.UpdateProfile)
		}
		git := v1.Group("/git")
		{
			git.Use(infrastructure.Authorize())
		}
		project := v1.Group("/project")
		{
			project.Use(infrastructure.Authorize())
			project.POST("/save", projectController.Save)
			project.PUT("/update/:id", projectController.Update)
			project.GET("/getSingleOne/:id", projectController.GetSingleOne)
			project.GET("/query", projectController.Query)
		}
		model := v1.Group("/model")
		{
			model.Use(infrastructure.Authorize())
		}
		dataset := v1.Group("/dataset")
		{
			dataset.Use(infrastructure.Authorize())
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
