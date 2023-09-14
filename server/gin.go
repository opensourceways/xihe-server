package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/community-robot-lib/interrupts"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/opensourceways/xihe-server/app"
	asyncapp "github.com/opensourceways/xihe-server/async-server/app"
	asyncrepoimpl "github.com/opensourceways/xihe-server/async-server/infrastructure/repositoryimpl"
	bigmodelapp "github.com/opensourceways/xihe-server/bigmodel/app"
	bigmodelasynccli "github.com/opensourceways/xihe-server/bigmodel/infrastructure/asynccli"
	"github.com/opensourceways/xihe-server/bigmodel/infrastructure/bigmodels"
	bigmodelrepo "github.com/opensourceways/xihe-server/bigmodel/infrastructure/repositoryimpl"
	cloudapp "github.com/opensourceways/xihe-server/cloud/app"
	cloudmsg "github.com/opensourceways/xihe-server/cloud/infrastructure/messageadapter"
	cloudrepo "github.com/opensourceways/xihe-server/cloud/infrastructure/repositoryimpl"
	"github.com/opensourceways/xihe-server/common/infrastructure/kafka"
	competitionapp "github.com/opensourceways/xihe-server/competition/app"
	competitionmsg "github.com/opensourceways/xihe-server/competition/infrastructure/messageadapter"
	competitionrepo "github.com/opensourceways/xihe-server/competition/infrastructure/repositoryimpl"
	competitionusercli "github.com/opensourceways/xihe-server/competition/infrastructure/usercli"
	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/controller"
	courseapp "github.com/opensourceways/xihe-server/course/app"
	courserepo "github.com/opensourceways/xihe-server/course/infrastructure/repositoryimpl"
	courseusercli "github.com/opensourceways/xihe-server/course/infrastructure/usercli"
	"github.com/opensourceways/xihe-server/docs"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/infrastructure/authingimpl"
	"github.com/opensourceways/xihe-server/infrastructure/challengeimpl"
	"github.com/opensourceways/xihe-server/infrastructure/competitionimpl"
	"github.com/opensourceways/xihe-server/infrastructure/finetuneimpl"
	"github.com/opensourceways/xihe-server/infrastructure/gitlab"
	"github.com/opensourceways/xihe-server/infrastructure/messages"
	"github.com/opensourceways/xihe-server/infrastructure/mongodb"
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
	"github.com/opensourceways/xihe-server/infrastructure/trainingimpl"
	pointsapp "github.com/opensourceways/xihe-server/points/app"
	pointsrepo "github.com/opensourceways/xihe-server/points/infrastructure/repositoryadapter"
	userapp "github.com/opensourceways/xihe-server/user/app"
	userrepoimpl "github.com/opensourceways/xihe-server/user/infrastructure/repositoryimpl"
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

	interrupts.ListenAndServe(srv, timeout)
}

// setRouter init router
func setRouter(engine *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Title = "xihe"
	docs.SwaggerInfo.Description = "set header: 'PRIVATE-TOKEN=xxx'"

	newPlatformRepository := func(token, namespace string) platform.Repository {
		return gitlab.NewRepositoryService(gitlab.UserInfo{
			Token:     token,
			Namespace: namespace,
		})
	}

	collections := &cfg.Mongodb.Collections

	proj := repositories.NewProjectRepository(
		mongodb.NewProjectMapper(collections.Project),
	)

	model := repositories.NewModelRepository(
		mongodb.NewModelMapper(collections.Model),
	)

	dataset := repositories.NewDatasetRepository(
		mongodb.NewDatasetMapper(collections.Dataset),
	)

	// user := repositories.NewUserRepository(
	// 	mongodb.NewUserMapper(collections.User),
	// )
	user := userrepoimpl.NewUserRepo(
		mongodb.NewCollection(collections.User),
	)

	login := repositories.NewLoginRepository(
		mongodb.NewLoginMapper(collections.Login),
	)

	like := repositories.NewLikeRepository(
		mongodb.NewLikeMapper(collections.Like),
	)

	activity := repositories.NewActivityRepository(
		mongodb.NewActivityMapper(
			collections.Activity,
			cfg.ActivityKeepNum,
		),
	)

	training := repositories.NewTrainingRepository(
		mongodb.NewTrainingMapper(
			collections.Training,
		),
	)

	finetune := repositories.NewFinetuneRepository(
		mongodb.NewFinetuneMapper(
			collections.Finetune,
		),
	)

	inference := repositories.NewInferenceRepository(
		mongodb.NewInferenceMapper(
			collections.Inference,
		),
	)

	tags := repositories.NewTagsRepository(
		mongodb.NewTagsMapper(collections.Tag),
	)

	competition := repositories.NewCompetitionRepository(
		mongodb.NewCompetitionMapper(collections.Competition),
	)

	aiquestion := repositories.NewAIQuestionRepository(
		mongodb.NewAIQuestionMapper(
			collections.AIQuestion, collections.QuestionPool,
		),
	)

	bigmodel := bigmodels.NewBigModelService()
	gitlabUser := gitlab.NewUserSerivce()
	gitlabRepo := gitlab.NewRepoFile()
	authingUser := authingimpl.NewAuthingUser()
	publisher := kafka.PublisherAdapter()
	sender := messages.NewMessageSender(&cfg.MQTopics, publisher)
	trainingAdapter := trainingimpl.NewTraining(&cfg.Training)
	finetuneImpl := finetuneimpl.NewFinetune(&cfg.Finetune)
	uploader := competitionimpl.NewCompetitionService()
	challengeHelper := challengeimpl.NewChallenge(&cfg.Challenge)

	userRegService := userapp.NewRegService(
		userrepoimpl.NewUserRegRepo(
			mongodb.NewCollection(collections.Registration),
		),
	)

	loginService := app.NewLoginService(login, sender)

	asyncAppService := asyncapp.NewTaskService(asyncrepoimpl.NewAsyncTaskRepo(&cfg.Postgresql.Async))

	competitionAppService := competitionapp.NewCompetitionService(
		competitionrepo.NewCompetitionRepo(mongodb.NewCollection(collections.Competition)),
		competitionrepo.NewWorkRepo(mongodb.NewCollection(collections.CompetitionWork)),
		competitionrepo.NewPlayerRepo(mongodb.NewCollection(collections.CompetitionPlayer)),
		competitionmsg.MessageAdapter(&cfg.Competition.Message, publisher), uploader,
		competitionusercli.NewUserCli(userRegService),
	)

	courseAppService := courseapp.NewCourseService(
		courseusercli.NewUserCli(userRegService),
		proj,
		courserepo.NewCourseRepo(mongodb.NewCollection(collections.Course)),
		courserepo.NewPlayerRepo(mongodb.NewCollection(collections.CoursePlayer)),
		courserepo.NewWorkRepo(mongodb.NewCollection(collections.CourseWork)),
		courserepo.NewRecordRepo(mongodb.NewCollection(collections.CourseRecord)),
	)

	cloudAppService := cloudapp.NewCloudService(
		cloudrepo.NewCloudRepo(mongodb.NewCollection(collections.CloudConf)),
		cloudrepo.NewPodRepo(&cfg.Postgresql.Cloud),
		cloudmsg.NewPublisher(&cfg.Cloud, publisher),
	)

	bigmodelAppService := bigmodelapp.NewBigModelService(
		bigmodel, user,
		bigmodelrepo.NewLuoJiaRepo(mongodb.NewCollection(collections.LuoJia)),
		bigmodelrepo.NewWuKongRepo(mongodb.NewCollection(collections.WuKong)),
		bigmodelrepo.NewWuKongPictureRepo(mongodb.NewCollection(collections.WuKongPicture)),
		bigmodelasynccli.NewAsyncCli(asyncAppService),
		sender,
		bigmodelrepo.NewApiService(mongodb.NewCollection(collections.ApiApply)),
		bigmodelrepo.NewApiInfo(mongodb.NewCollection(collections.ApiInfo)),
		userRegService,
	)

	projectService := app.NewProjectService(user, proj, model, dataset, activity, nil, sender)

	modelService := app.NewModelService(user, model, proj, dataset, activity, nil, sender)

	datasetService := app.NewDatasetService(user, dataset, proj, model, activity, nil, sender)

	pointsAppService := pointsapp.NewUserPointsAppService(
		pointsrepo.TaskAdapter(
			mongodb.NewCollection(collections.PointsTask),
		),
		pointsrepo.UserPointsAdapter(
			mongodb.NewCollection(collections.UserPoints), &cfg.Points.Repo,
		),
	)

	userAppService := userapp.NewUserService(
		user, gitlabUser, sender, pointsAppService, controller.EncryptHelperToken(),
	)

	v1 := engine.Group(docs.SwaggerInfo.BasePath)
	{
		controller.AddRouterForProjectController(
			v1, user, proj, model, dataset, activity, tags, like, sender,
			newPlatformRepository,
		)

		controller.AddRouterForModelController(
			v1, user, model, proj, dataset, activity, tags, like, sender,
			newPlatformRepository,
		)

		controller.AddRouterForDatasetController(
			v1, user, dataset, model, proj, activity, tags, like, sender,
			newPlatformRepository,
		)

		controller.AddRouterForUserController(
			v1, userAppService, user,
			authingUser, loginService, userRegService,
		)

		controller.AddRouterForLoginController(
			v1, userAppService, authingUser, loginService,
		)

		controller.AddRouterForLikeController(
			v1, like, user, proj, model, dataset, activity, sender,
		)

		controller.AddRouterForActivityController(
			v1, activity, user, proj, model, dataset,
		)

		controller.AddRouterForTagsController(
			v1, tags,
		)

		controller.AddRouterForBigModelController(
			v1, bigmodelAppService, userRegService,
		)

		controller.AddRouterForTrainingController(
			v1, trainingAdapter, training, model, proj, dataset, sender,
		)

		controller.AddRouterForFinetuneController(
			v1, finetuneImpl, finetune, sender,
		)

		controller.AddRouterForRepoFileController(
			v1, gitlabRepo, model, proj, dataset, sender, userAppService,
		)

		controller.AddRouterForInferenceController(
			v1, gitlabRepo, inference, proj, sender,
		)

		controller.AddRouterForSearchController(
			v1, user, proj, model, dataset,
		)

		controller.AddRouterForCompetitionController(
			v1, competitionAppService, userRegService, proj,
		)

		controller.AddRouterForChallengeController(
			v1, competition, aiquestion, challengeHelper,
		)

		controller.AddRouterForCourseController(
			v1, courseAppService, userRegService, proj, user,
		)

		controller.AddRouterForHomeController(
			v1, courseAppService, competitionAppService, projectService, modelService, datasetService,
		)

		controller.AddRouterForCloudController(
			v1, cloudAppService,
		)

		controller.AddRouterForUserPointsController(v1, pointsAppService)

	}

	engine.UseRawPath = true
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
