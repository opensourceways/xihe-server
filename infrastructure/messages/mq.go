package messages

import (
	"github.com/opensourceways/community-robot-lib/kafka"
	"github.com/opensourceways/community-robot-lib/mq"
	kfklib "github.com/opensourceways/kafka-lib/agent"
	kfkmq "github.com/opensourceways/kafka-lib/mq"
	redislib "github.com/opensourceways/redis-lib"
	"github.com/sirupsen/logrus"
)

const (
	kfkQueueName = "xihe-kafka-queue"
)

func Init(cfg mq.MQConfig, log *logrus.Entry, topic Topics) error {
	topics = topic

	err := kafka.Init(
		mq.Addresses(cfg.Addresses...),
		mq.Log(log),
	)
	if err != nil {
		return err
	}

	return kafka.Connect()
}

func Exit(log *logrus.Entry) {
	if err := kafka.Disconnect(); err != nil {
		log.Errorf("exit kafka, err:%v", err)
	}
}

func InitKfkLib(kfkCfg kfklib.Config, redisCfg redislib.Config, log kfkmq.Logger, topic Topics) (err error) {
	topics = topic

	if err = redislib.Init(&redisCfg); err != nil {
		return
	}

	if err = kfklib.Init(&kfkCfg, log, redislib.DAO(), kfkQueueName); err != nil {
		return
	}

	return
}

func KfkLibExit() {
	kfklib.Exit()
}
