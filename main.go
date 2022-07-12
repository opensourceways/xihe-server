package main

import (
	"flag"
	"os"

	"github.com/opensourceways/community-robot-lib/logrusutil"
	liboptions "github.com/opensourceways/community-robot-lib/options"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/infrastructure/authing"
	"github.com/opensourceways/xihe-server/infrastructure/mq"
	"github.com/opensourceways/xihe-server/infrastructure/redis"
	"github.com/opensourceways/xihe-server/server"
)

type options struct {
	service liboptions.ServiceOptions
}

func (o *options) Validate() error {
	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options
	o.service.AddFlags(fs)
	fs.Parse(args)
	return o
}

func main() {
	logrusutil.ComponentInit("xihe")
	o := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if o.service.ConfigFile == "" {
		o.service.ConfigFile = "./conf/app.conf.yaml"
	}
	if err := o.Validate(); err != nil {
		logrus.Fatalf("Invalid options, err:%s", err.Error())
	}
	cfg, err := config.LoadConfig(o.service.ConfigFile)
	if err != nil {
		logrus.Fatalf("load config, err:%s", err.Error())
	}
	authing.Init(&cfg.Authing)
	redis.InitRedis(cfg)
	mq.InitMQ(cfg)
	go mq.StartEventLinsten(mq.ProjectLikeCountIncreaseEvent, "test", app.ReceiveFunction)
	server.StartWebServer(o.service.Port, o.service.GracePeriod, cfg)
}
