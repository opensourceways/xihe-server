package messages

import (
	kfklib "github.com/opensourceways/kafka-lib/agent"
	kfkmq "github.com/opensourceways/kafka-lib/mq"
	redislib "github.com/opensourceways/redis-lib"
)

const (
	kfkQueueName = "xihe-kafka-queue"
)

func InitKfkLib(kfkCfg kfklib.Config, log kfkmq.Logger, topic Topics) (err error) {
	topics = topic

	return kfklib.Init(&kfkCfg, log, redislib.DAO(), kfkQueueName)
}

func KfkLibExit() {
	kfklib.Exit()
}
