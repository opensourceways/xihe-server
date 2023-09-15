package main

import (
	"flag"
	"os"

	liboptions "github.com/opensourceways/community-robot-lib/options"
	"github.com/opensourceways/server-common-lib/logrusutil"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/async-server/app"
	"github.com/opensourceways/xihe-server/async-server/config"
	"github.com/opensourceways/xihe-server/async-server/infrastructure/bigmodelimpl"
	"github.com/opensourceways/xihe-server/async-server/infrastructure/poolimpl"
	"github.com/opensourceways/xihe-server/async-server/infrastructure/repositoryimpl"
	"github.com/opensourceways/xihe-server/async-server/infrastructure/watchimpl"
	bigmodelapp "github.com/opensourceways/xihe-server/bigmodel/app"
	"github.com/opensourceways/xihe-server/bigmodel/infrastructure/bigmodels"
	bmmsgadapter "github.com/opensourceways/xihe-server/bigmodel/infrastructure/messageadapter"
	"github.com/opensourceways/xihe-server/common/infrastructure/kafka"
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
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
	cfg := new(config.Config)
	if err := config.LoadConfig(o.service.ConfigFile, cfg); err != nil {
		logrus.Fatalf("load config, err:%s", err.Error())
	}

	// bigmodel
	if err := bigmodels.Init(&cfg.BigModel.Config); err != nil {
		logrus.Fatalf("initialize big model failed, err:%s", err.Error())
	}

	// mq
	if err := kafka.Init(&cfg.MQ, log, nil); err != nil {
		log.Fatalf("initialize mq failed, err:%v", err)
	}

	defer kafka.Exit()

	// postgresql
	if err := pgsql.Init(&cfg.Postgresql.DB); err != nil {
		logrus.Fatalf("init db, err:%s", err.Error())
	}

	// pool
	if err := poolimpl.Init(&cfg.Pool); err != nil {
		logrus.Fatalf("init pool, err:%s", err.Error())
	}

	// bigmodel & sender
	bm := bigmodels.NewBigModelService()
	publisher := kafka.PublisherAdapter()

	// aysnc.bigmodel.bigmodel
	bigmodel := bigmodelimpl.NewBigModelImpl(
		bigmodelapp.NewAsyncBigModelService(bm, bmmsgadapter.NewMessageAdapter(
			&cfg.BigModel.Message, publisher,
		)),
	)

	// repo
	asyncWuKongRepo := repositoryimpl.NewAsyncTaskRepo(&cfg.Postgresql.Config)

	// async app
	asyncAppService := app.NewAsyncService(
		bigmodel,
		poolimpl.NewPoolImpl(),
		asyncWuKongRepo,
	)

	// watch
	w := watchimpl.NewWather(
		cfg.Watcher,
		asyncWuKongRepo,
		map[string]func(string, int64) error{
			"wukong":      asyncAppService.AsyncWuKong,
			"wukong_4img": asyncAppService.AsyncWuKong4Img,
		},
	)

	w.Run()
	defer w.Exit()
}
