package main

import (
	"flag"
	"os"

	"github.com/opensourceways/community-robot-lib/logrusutil"
	liboptions "github.com/opensourceways/community-robot-lib/options"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/controller"
	"github.com/opensourceways/xihe-server/infrastructure/authingimpl"
	"github.com/opensourceways/xihe-server/infrastructure/bigmodels"
	"github.com/opensourceways/xihe-server/infrastructure/competitionimpl"
	"github.com/opensourceways/xihe-server/infrastructure/gitlab"
	"github.com/opensourceways/xihe-server/infrastructure/messages"
	"github.com/opensourceways/xihe-server/infrastructure/mongodb"
	"github.com/opensourceways/xihe-server/server"
)

type options struct {
	service     liboptions.ServiceOptions
	enableDebug bool
}

func (o *options) Validate() error {
	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.service.AddFlags(fs)

	fs.BoolVar(
		&o.enableDebug, "enable_debug", false,
		"whether to enable debug model.",
	)

	fs.Parse(args)
	return o
}

func main() {
	logrusutil.ComponentInit("xihe")
	log := logrus.NewEntry(logrus.StandardLogger())

	o := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
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
	if err := bigmodels.Init(&cfg.BigModel); err != nil {
		logrus.Fatalf("initialize big model failed, err:%s", err.Error())
	}

	// gitlab
	if err := gitlab.Init(&cfg.Gitlab); err != nil {
		logrus.Fatalf("initialize gitlab failed, err:%s", err.Error())
	}

	// competition
	if err := competitionimpl.Init(&cfg.Competition); err != nil {
		logrus.Fatalf("initialize competition failed, err:%s", err.Error())
	}

	// authing
	authingimpl.Init(&cfg.Authing)

	// controller
	api := &cfg.API
	api.MaxPictureSizeToVQA = cfg.BigModel.MaxPictureSizeToVQA
	api.MaxPictureSizeToDescribe = cfg.BigModel.MaxPictureSizeToDescribe

	if err := controller.Init(api, log); err != nil {
		logrus.Fatalf("initialize api controller failed, err:%s", err.Error())
	}

	// mq
	if err := messages.Init(cfg.GetMQConfig(), log, cfg.MQ.Topics); err != nil {
		log.Fatalf("initialize mq failed, err:%v", err)
	}

	defer messages.Exit(log)

	// mongo
	m := &cfg.Mongodb
	if err := mongodb.Initialize(m.DBConn, m.DBName); err != nil {
		logrus.Fatalf("initialize mongodb failed, err:%s", err.Error())
	}

	defer mongodb.Close()

	// postgresql
	if err := pgsql.Init(&cfg.Postgresql.DB); err != nil {
		logrus.Fatalf("init db, err:%s", err.Error())
	}

	// cfg
	cfg.InitDomainConfig()
	cfg.InitAppConfig()

	// run
	server.StartWebServer(o.service.Port, o.service.GracePeriod, cfg)
}
