package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/opensourceways/community-robot-lib/logrusutil"
	liboptions "github.com/opensourceways/community-robot-lib/options"
	"github.com/sirupsen/logrus"

	aiccapp "github.com/opensourceways/xihe-server/aiccfinetune/app"
	aiccimpl "github.com/opensourceways/xihe-server/aiccfinetune/infrastructure/aiccfinetuneimpl"
	aiccrepo "github.com/opensourceways/xihe-server/aiccfinetune/infrastructure/repositoryimpl"
	aiccmq "github.com/opensourceways/xihe-server/aiccfinetune/messagequeue"
	"github.com/opensourceways/xihe-server/app"
	asyncapp "github.com/opensourceways/xihe-server/async-server/app"
	asyncrepo "github.com/opensourceways/xihe-server/async-server/infrastructure/repositoryimpl"
	bigmodelmq "github.com/opensourceways/xihe-server/bigmodel/messagequeue"
	cloudapp "github.com/opensourceways/xihe-server/cloud/app"
	"github.com/opensourceways/xihe-server/cloud/infrastructure/cloudimpl"
	cloudrepo "github.com/opensourceways/xihe-server/cloud/infrastructure/repositoryimpl"
	"github.com/opensourceways/xihe-server/common/infrastructure/kafka"
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	"github.com/opensourceways/xihe-server/infrastructure/evaluateimpl"
	"github.com/opensourceways/xihe-server/infrastructure/finetuneimpl"
	"github.com/opensourceways/xihe-server/infrastructure/inferenceimpl"
	"github.com/opensourceways/xihe-server/infrastructure/messages"
	"github.com/opensourceways/xihe-server/infrastructure/mongodb"
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
	"github.com/opensourceways/xihe-server/infrastructure/trainingimpl"
	"github.com/opensourceways/xihe-server/messagequeue"
	pointsapp "github.com/opensourceways/xihe-server/points/app"
	pointsrepo "github.com/opensourceways/xihe-server/points/infrastructure/repositoryadapter"
	pointsmq "github.com/opensourceways/xihe-server/points/messagequeue"
	userapp "github.com/opensourceways/xihe-server/user/app"
	userrepo "github.com/opensourceways/xihe-server/user/infrastructure/repositoryimpl"
	usermq "github.com/opensourceways/xihe-server/user/messagequeue"
)

type options struct {
	service     liboptions.ServiceOptions
	enableDebug bool
}

func (o *options) Validate() error {
	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) (options, error) {
	var o options

	o.service.AddFlags(fs)

	fs.BoolVar(
		&o.enableDebug, "enable_debug", false,
		"whether to enable debug model.",
	)

	err := fs.Parse(args)
	return o, err
}

func main() {
	logrusutil.ComponentInit("xihe")
	log := logrus.NewEntry(logrus.StandardLogger())

	o, err := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if err != nil {
		logrus.Fatalf("new options failed, err:%s", err.Error())
	}

	if err := o.Validate(); err != nil {
		logrus.Fatalf("Invalid options, err:%s", err.Error())
	}

	if o.enableDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug enabled.")
	}

	// cfg
	cfg := new(configuration)
	if err := loadConfig(o.service.ConfigFile, cfg); err != nil {
		logrus.Fatalf("load config, err:%s", err.Error())
	}

	if err := os.Remove(o.service.ConfigFile); err != nil {
		logrus.Fatalf("config file delete failed, err:%s", err.Error())
	}

	// mq
	if err = kafka.Init(&cfg.MQ, log, nil); err != nil {
		log.Fatalf("initialize mq failed, err:%v", err)
	}

	defer kafka.Exit()

	// mongo
	m := &cfg.Mongodb
	if err := mongodb.Initialize(m.DBConn, m.DBName, m.DBCert); err != nil {
		logrus.Fatalf("initialize mongodb failed, err:%s", err.Error())
	}

	defer mongodb.Close()

	// postgresql
	if err := pgsql.Init(&cfg.Postgresql.DB); err != nil {
		logrus.Fatalf("init db, err:%s", err.Error())
	}

	if err := os.Remove(cfg.Postgresql.DB.DBCert); err != nil {
		logrus.Fatalf("postgresql dbcert file delete failed, err:%s", err.Error())
	}

	// cfg
	cfg.initDomainConfig()

	// points
	if err = pointsSubscribesMessage(cfg, &cfg.MQTopics); err != nil {
		logrus.Errorf("points subscribes message failed, err:%s", err.Error())

		return
	}

	// bigmodel
	if err = bigmodelSubscribesMessage(cfg, &cfg.MQTopics); err != nil {
		logrus.Errorf("bigmodel subscribes message failed, err:%s", err.Error())

		return
	}

	// training
	if err = trainingSubscribesMessage(log, cfg); err != nil {
		logrus.Errorf("training subscribes message failed, err:%s", err.Error())

		return
	}

	// aicc finetune
	if err = aiccFinetuneSubscribesMessage(cfg); err != nil {
		logrus.Errorf("aicc finetune subscribes message failed, err:%s", err.Error())
		return
	}

	// user
	if err = userSubscribesMessage(cfg, &cfg.MQTopics.User); err != nil {
		logrus.Errorf("user subscribes message failed, err:%s", err.Error())

		return
	}

	// run
	run(newHandler(cfg, log), log, &cfg.MQTopics)
}

func pointsSubscribesMessage(cfg *configuration, topics *mqTopics) error {
	collections := &cfg.Mongodb.Collections

	return pointsmq.Subscribe(
		pointsapp.NewUserPointsAppMessageService(
			pointsrepo.TaskAdapter(
				mongodb.NewCollection(collections.PointsTask),
			),
			pointsrepo.UserPointsAdapter(
				mongodb.NewCollection(collections.UserPoints),
				&cfg.Points.Repo,
			),
		),
		[]string{
			topics.SignIn,
			topics.CompetitorApplied,
			topics.JupyterCreated,
			topics.PicturePublicized,
			topics.PictureLiked,
			topics.CourseApplied,
			topics.TrainingCreated,
			topics.ProjectCreated,
			topics.DatasetCreated,
			topics.ModelCreated,
			topics.ProjectDownloaded,
			topics.ModelDownloaded,
			topics.DatasetDownloaded,
			topics.User.UserSignedUp,
			topics.User.BioSet,
			topics.User.AvatarSet,
			// both add and remove action will send the event
			// but we only fullfill the MsgNormal
			// when the action need to be scoring
			topics.Like,
			topics.BigModelTopics.BigModelFinished,
			topics.BigModelTopics.InferenceAsyncFinish,
		},
		kafka.SubscriberAdapter(),
	)
}

func userSubscribesMessage(cfg *configuration, topics *userConfig) error {
	collections := &cfg.Mongodb.Collections

	return usermq.Subscribe(
		userapp.NewUserService(
			userrepo.NewUserRepo(
				mongodb.NewCollection(collections.User),
			),
			nil, nil, nil, nil,
		),
		kafka.SubscriberAdapter(),
		&topics.TopicConfig,
	)
}

func bigmodelSubscribesMessage(cfg *configuration, topics *mqTopics) error {
	return bigmodelmq.Subscribe(
		asyncapp.NewAsyncMessageService(
			asyncrepo.NewAsyncTaskRepo(&cfg.Postgresql.asyncconf),
		),
		&topics.BigModelTopics,
	)
}

func trainingSubscribesMessage(log *logrus.Entry, cfg *configuration) error {
	collections := &cfg.Mongodb.Collections

	return messagequeue.Subscribe(
		cfg.Training, cfg.MQTopics.TrainingCreated,
		app.NewTrainingService(
			trainingimpl.NewTraining(&trainingimpl.Config{}),
			repositories.NewTrainingRepository(
				mongodb.NewTrainingMapper(collections.Training),
			),
			nil, 0,
		),
		kafka.SubscriberAdapter(),
	)
}

func aiccFinetuneSubscribesMessage(cfg *configuration) error {
	collections := &cfg.Mongodb.Collections

	return aiccmq.Subscribe(
		cfg.AICCFinetune,
		cfg.MQTopics.AICCFinetuneCreated,
		aiccapp.NewAICCFinetuneService(
			aiccimpl.NewAICCFinetune(&aiccimpl.Config{}),
			nil,
			nil,
			aiccrepo.NewAICCFinetuneRepo(
				mongodb.NewCollection(collections.AICCFinetune),
			),
			0,
		),
		kafka.SubscriberAdapter(),
	)
}

func newHandler(cfg *configuration, log *logrus.Entry) *handler {
	collections := &cfg.Mongodb.Collections

	userRepo := userrepo.NewUserRepo(mongodb.NewCollection(collections.User))

	h := &handler{
		log:      log,
		maxRetry: cfg.MaxRetry,

		project: app.NewProjectMessageService(
			repositories.NewProjectRepository(
				mongodb.NewProjectMapper(collections.Project),
			),
		),

		dataset: app.NewDatasetMessageService(
			repositories.NewDatasetRepository(
				mongodb.NewDatasetMapper(collections.Dataset),
			),
		),

		model: app.NewModelMessageService(
			repositories.NewModelRepository(
				mongodb.NewModelMapper(collections.Model),
			),
		),

		inference: app.NewInferenceMessageService(
			repositories.NewInferenceRepository(
				mongodb.NewInferenceMapper(collections.Inference),
			),
			userRepo,
			inferenceimpl.NewInference(&cfg.Inference),
		),

		evaluate: app.NewEvaluateMessageService(
			repositories.NewEvaluateRepository(
				mongodb.NewEvaluateMapper(collections.Evaluate),
			),
			evaluateimpl.NewEvaluate(&cfg.Evaluate.Config),
			cfg.Evaluate.SurvivalTime,
		),

		cloud: cloudapp.NewCloudMessageService(
			cloudrepo.NewPodRepo(&cfg.Postgresql.cloudconf),
			cloudimpl.NewCloud(&cfg.Cloud.Config),
			int64(cfg.Cloud.SurvivalTime),
		),

		async: asyncapp.NewAsyncMessageService(
			asyncrepo.NewAsyncTaskRepo(&cfg.Postgresql.asyncconf),
		),
	}

	fc := cfg.getFinetuneConfig()
	h.finetune = app.NewFinetuneMessageService(
		finetuneimpl.NewFinetune(&fc),
		repositories.NewFinetuneRepository(
			mongodb.NewFinetuneMapper(collections.Finetune),
		),
	)

	return h
}

func run(h *handler, log *logrus.Entry, topics *mqTopics) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	defer wg.Wait()

	called := false
	ctx, done := context.WithCancel(context.Background())

	defer func() {
		if !called {
			called = true
			done()
		}
	}()

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		select {
		case <-ctx.Done():
			log.Info("receive done. exit normally")
			return

		case <-sig:
			log.Info("receive exit signal")
			done()
			called = true
			return
		}
	}(ctx)

	err := messages.Subscribe(ctx, h, log, &topics.Topics, kafka.SubscriberAdapter())
	if err != nil {
		log.Errorf("subscribe failed, err:%v", err)
	}
}
