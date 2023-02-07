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

	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/controller"
	"github.com/opensourceways/xihe-server/docs"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/infrastructure/authing"
	"github.com/opensourceways/xihe-server/infrastructure/bigmodels"
	"github.com/opensourceways/xihe-server/infrastructure/challengeimpl"
	"github.com/opensourceways/xihe-server/infrastructure/competitionimpl"
	"github.com/opensourceways/xihe-server/infrastructure/finetuneimpl"
	"github.com/opensourceways/xihe-server/infrastructure/gitlab"
	"github.com/opensourceways/xihe-server/infrastructure/messages"
	"github.com/opensourceways/xihe-server/infrastructure/mongodb"
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
	"github.com/opensourceways/xihe-server/infrastructure/trainingimpl"
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

//setRouter init router
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

	user := repositories.NewUserRepository(
		mongodb.NewUserMapper(collections.User),
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

	evaluate := repositories.NewEvaluateRepository(
		mongodb.NewEvaluateMapper(
			collections.Evaluate,
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

	luojia := repositories.NewLuoJiaRepository(
		mongodb.NewLuoJiaMapper(collections.LuoJia),
	)

	wukong := repositories.NewWuKongRepository(
		mongodb.NewWuKongMapper(collections.WuKong),
	)

	wukongPicture := repositories.NewWuKongPictureRepository(
		mongodb.NewWuKongPictureMapper(collections.WuKongPicture),
	)

	bigmodel := bigmodels.NewBigModelService()
	gitlabUser := gitlab.NewUserSerivce()
	gitlabRepo := gitlab.NewRepoFile()
	authingUser := authing.NewAuthingUser()
	sender := messages.NewMessageSender()
	trainingAdapter := trainingimpl.NewTraining(&cfg.Training)
	finetuneImpl := finetuneimpl.NewFinetune(&cfg.Finetune)
	uploader := competitionimpl.NewCompetitionService()
	challengeHelper := challengeimpl.NewChallenge(&cfg.Challenge)

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
			v1, user, gitlabUser,
			authingUser, sender,
		)

		controller.AddRouterForLoginController(
			v1, user, gitlabUser, authingUser, login, sender,
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
			v1, user, bigmodel, luojia, wukong, wukongPicture, sender,
		)

		controller.AddRouterForTrainingController(
			v1, trainingAdapter, training, model, proj, dataset, sender,
		)

		controller.AddRouterForFinetuneController(
			v1, finetuneImpl, finetune, sender,
		)

		controller.AddRouterForRepoFileController(
			v1, gitlabRepo, model, proj, dataset, sender,
		)

		controller.AddRouterForInferenceController(
			v1, gitlabRepo, inference, proj, sender,
		)

		controller.AddRouterForEvaluateController(
			v1, evaluate, training, sender,
		)

		controller.AddRouterForSearchController(
			v1, user, proj, model, dataset,
		)

		controller.AddRouterForCompetitionController(
			v1, competition, proj, sender, uploader,
		)

		controller.AddRouterForChallengeController(
			v1, competition, aiquestion, challengeHelper,
		)
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
