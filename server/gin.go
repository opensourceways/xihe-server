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

	asyncapp "github.com/opensourceways/xihe-extra-services/async-server/app"
	asyncrepoimpl "github.com/opensourceways/xihe-extra-services/async-server/infrastructure/repositoryimpl"
	"github.com/opensourceways/xihe-server/app"
	bigmodelapp "github.com/opensourceways/xihe-server/bigmodel/app"
	bigmodelasynccli "github.com/opensourceways/xihe-server/bigmodel/infrastructure/asynccli"
	"github.com/opensourceways/xihe-server/bigmodel/infrastructure/bigmodels"
	bigmodelmsg "github.com/opensourceways/xihe-server/bigmodel/infrastructure/messageadapter"
	bigmodelrepo "github.com/opensourceways/xihe-server/bigmodel/infrastructure/repositoryimpl"
	cloudapp "github.com/opensourceways/xihe-server/cloud/app"
	cloudrepo "github.com/opensourceways/xihe-server/cloud/infrastructure/repositoryimpl"
	"github.com/opensourceways/xihe-server/common/infrastructure/kafka"
	competitionapp "github.com/opensourceways/xihe-server/competition/app"
	competitionmsg "github.com/opensourceways/xihe-server/competition/infrastructure/messageadapter"
	competitionrepo "github.com/opensourceways/xihe-server/competition/infrastructure/repositoryimpl"
	competitionusercli "github.com/opensourceways/xihe-server/competition/infrastructure/usercli"
	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/controller"
	courseapp "github.com/opensourceways/xihe-server/course/app"
	coursemsg "github.com/opensourceways/xihe-server/course/infrastructure/messageadapter"
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
	pointsservice "github.com/opensourceways/xihe-server/points/domain/service"
	pointsrepo "github.com/opensourceways/xihe-server/points/infrastructure/repositoryadapter"
	"github.com/opensourceways/xihe-server/points/infrastructure/taskdocimpl"
	promotionapp "github.com/opensourceways/xihe-server/promotion/app"
	prmotionservice "github.com/opensourceways/xihe-server/promotion/domain/service"
	promotionadapter "github.com/opensourceways/xihe-server/promotion/infrastructure/repositoryadapter"
	promotionuseradapter "github.com/opensourceways/xihe-server/promotion/infrastructure/useradapter"
	userapp "github.com/opensourceways/xihe-server/user/app"
	usermsg "github.com/opensourceways/xihe-server/user/infrastructure/messageadapter"
	userrepoimpl "github.com/opensourceways/xihe-server/user/infrastructure/repositoryimpl"
)

func StartWebServer(port int, timeout time.Duration, cfg *config.Config) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logRequest())

	if err := setRouter(r, cfg); err != nil {
		logrus.Error(err)

		return
	}

	r.Use(controller.ClearSenstiveInfoMiddleware())

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: time.Duration(cfg.ReadHeaderTimeout) * time.Second,
		Handler:           r,
	}

	defer interrupts.WaitForGracefulShutdown()

	interrupts.ListenAndServe(srv, timeout)
}

// setRouter init router
func setRouter(engine *gin.Engine, cfg *config.Config) error {
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

	promotionRepo := promotionadapter.PromotionAdapter(mongodb.NewCollection(collections.Promotion))
	promotionPointRepo := promotionadapter.PointsAdapter(mongodb.NewCollection(collections.PromotionPoint))
	promotionTaskRepo := promotionadapter.TaskAdapter(mongodb.NewCollection(collections.PromotionTask))

	bigmodel := bigmodels.NewBigModelService()
	gitlabUser := gitlab.NewUserSerivce()
	gitlabRepo := gitlab.NewRepoFile()
	authingUser := authingimpl.NewAuthingUser()
	publisher := kafka.PublisherAdapter()
	operater := kafka.OperateLogPublisherAdapter(cfg.MQTopics.OperateLog, publisher)
	trainingAdapter := trainingimpl.NewTraining(&cfg.Training.Config)
	repoAdapter := messages.NewDownloadMessageAdapter(cfg.MQTopics.Download, &cfg.Download, publisher, operater)
	finetuneImpl := finetuneimpl.NewFinetune(&cfg.Finetune)
	uploader := competitionimpl.NewCompetitionService()
	challengeHelper := challengeimpl.NewChallenge(&cfg.Challenge)
	likeAdapter := messages.NewLikeMessageAdapter(cfg.MQTopics.Like, &cfg.Like, publisher)

	// sender
	sender := messages.NewMessageSender(&cfg.MQTopics, publisher)
	// resource producer
	resProducer := messages.NewResourceMessageAdapter(&cfg.Resource, publisher, operater)

	userRegService := userapp.NewRegService(
		userrepoimpl.NewUserRegRepo(
			mongodb.NewCollection(collections.Registration),
		),
	)

	loginService := app.NewLoginService(
		login, messages.NewSignInMessageAdapter(&cfg.SignIn, publisher),
	)

	promotionPointTaskService, err := prmotionservice.NewPointsTaskService(promotionPointRepo, promotionTaskRepo)
	if err != nil {
		return err
	}

	asyncAppService := asyncapp.NewTaskService(asyncrepoimpl.NewAsyncTaskRepo(&cfg.Postgresql.Async))

	competitionAppService := competitionapp.NewCompetitionService(
		competitionrepo.NewCompetitionRepo(mongodb.NewCollection(collections.Competition)),
		competitionrepo.NewWorkRepo(mongodb.NewCollection(collections.CompetitionWork)),
		competitionrepo.NewPlayerRepo(mongodb.NewCollection(collections.CompetitionPlayer)),
		competitionmsg.MessageAdapter(&cfg.Competition.Message, publisher), uploader,
		competitionusercli.NewUserCli(userRegService),
		user,
	)

	courseAppService := courseapp.NewCourseService(
		courseusercli.NewUserCli(userRegService),
		proj,
		courserepo.NewCourseRepo(mongodb.NewCollection(collections.Course)),
		courserepo.NewPlayerRepo(mongodb.NewCollection(collections.CoursePlayer)),
		courserepo.NewWorkRepo(mongodb.NewCollection(collections.CourseWork)),
		courserepo.NewRecordRepo(mongodb.NewCollection(collections.CourseRecord)),
		coursemsg.MessageAdapter(&cfg.Course.Message, publisher),
		user,
	)

	bigmodelAppService := bigmodelapp.NewBigModelService(
		bigmodel, user,
		bigmodelrepo.NewLuoJiaRepo(mongodb.NewCollection(collections.LuoJia)),
		bigmodelrepo.NewWuKongRepo(mongodb.NewCollection(collections.WuKong)),
		bigmodelrepo.NewWuKongPictureRepo(mongodb.NewCollection(collections.WuKongPicture)),
		bigmodelasynccli.NewAsyncCli(asyncAppService),
		bigmodelmsg.NewMessageAdapter(&cfg.BigModel.Message, publisher),
		bigmodelrepo.NewApiService(mongodb.NewCollection(collections.ApiApply)),
		bigmodelrepo.NewApiInfo(mongodb.NewCollection(collections.ApiInfo)),
		userRegService,
	)

	cloudAppService := cloudapp.NewCloudService(
		cloudrepo.NewCloudRepo(mongodb.NewCollection(collections.CloudConf)),
		cloudrepo.NewPodRepo(&cfg.Postgresql.Cloud),
	)

	projectService := app.NewProjectService(user, proj, model, dataset, activity, nil, resProducer)

	modelService := app.NewModelService(user, model, proj, dataset, activity, nil, resProducer)

	datasetService := app.NewDatasetService(user, dataset, proj, model, activity, nil, resProducer)

	v1 := engine.Group(docs.SwaggerInfo.BasePath)

	pointsAppService, err := addRouterForUserPointsController(v1, cfg)
	if err != nil {
		return err
	}

	userAppService := userapp.NewUserService(
		user, gitlabUser, usermsg.MessageAdapter(&cfg.User.Message, publisher),
		pointsAppService, controller.EncryptHelperToken(),
	)

	promotionAppService := promotionapp.NewPromotionService(
		prmotionservice.NewPromotionUserService(
			promotionuseradapter.NewUserAdapter(userRegService), promotionRepo),
		promotionPointTaskService,
		promotionRepo,
	)

	promotionpointsAppService := promotionapp.NewPointsService(promotionPointTaskService)

	{
		controller.AddRouterForProjectController(
			v1, user, proj, model, dataset, activity, tags, like, resProducer,
			newPlatformRepository,
		)

		controller.AddRouterForModelController(
			v1, user, model, proj, dataset, activity, tags, like, resProducer,
			newPlatformRepository,
		)

		controller.AddRouterForDatasetController(
			v1, user, dataset, model, proj, activity, tags, like, resProducer,
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
			v1, like, user, proj, model, dataset, activity, likeAdapter,
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
			v1, trainingAdapter, training, model, proj, dataset,
			messages.NewTrainingMessageAdapter(
				&cfg.Training.Message, publisher,
			),
		)

		controller.AddRouterForFinetuneController(
			v1, finetuneImpl, finetune, sender,
		)

		controller.AddRouterForRepoFileController(
			v1, gitlabRepo, model, proj, dataset, repoAdapter, userAppService,
		)

		controller.AddRouterForSearchController(
			v1, user, proj, model, dataset,
		)

		controller.AddRouterForCompetitionController(
			v1, competitionAppService, userRegService, proj,
		)

		controller.AddRouterForPromotionController(
			v1, promotionAppService, promotionpointsAppService,
		)

		controller.AddRouterForChallengeController(
			v1, competition, aiquestion, challengeHelper, user,
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
	}

	engine.UseRawPath = true
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return nil
}

func addRouterForUserPointsController(g *gin.RouterGroup, cfg *config.Config) (
	pointsapp.UserPointsAppService, error,
) {
	collections := &cfg.Mongodb.Collections

	taskRepo := pointsrepo.TaskAdapter(
		mongodb.NewCollection(collections.PointsTask),
	)

	taskdoc, err := taskdocimpl.Init(&cfg.Points.TaskDoc)
	if err != nil {
		return nil, err
	}

	taskService, err := pointsservice.InitTaskService(taskRepo, taskdoc)
	if err != nil {
		return nil, err
	}

	pointsAppService := pointsapp.NewUserPointsAppService(
		taskRepo,
		pointsrepo.UserPointsAdapter(
			mongodb.NewCollection(collections.UserPoints), &cfg.Points.Repo,
		),
	)

	controller.AddRouterForUserPointsController(
		g, pointsAppService,
		pointsapp.NewTaskAppService(taskService, taskRepo),
	)

	return pointsAppService, nil
}

func logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()

		l := controller.GetOperateLog(c)
		logrus.Infof(
			"| %d | %d | %s | %s | %s",
			c.Writer.Status(),
			endTime.Sub(startTime),
			c.Request.Method,
			c.Request.RequestURI,
			l,
		)
	}
}
