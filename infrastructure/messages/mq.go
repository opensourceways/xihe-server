package messages

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	kfkmq "github.com/opensourceways/kafka-lib/mq"
	redislib "github.com/opensourceways/redis-lib"
)

const (
	kfkQueueName = "xihe-kafka-queue"
)

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
