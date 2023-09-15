package messagequeue

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain"
)

const (
	retryNum = 3

	handleNameTrainingCreated = "training_created"
)

func Subscribe(
	log *logrus.Entry,
	cfg TrainingConfig,
	s app.TrainingService,
	subscriber message.Subscriber,
) (err error) {
	c := &consumer{log: log, cfg: cfg, s: s}

	// training created
	err = subscriber.SubscribeWithStrategyOfRetry(
		handleNameTrainingCreated,
		c.handleEventTrainingCreated,
		[]string{cfg.Topics.TrainingCreated}, retryNum,
	)

	return
}

type consumer struct {
	log *logrus.Entry
	cfg TrainingConfig
	s   app.TrainingService
}

func (c *consumer) handleEventTrainingCreated(body []byte, h map[string]string) (err error) {
	b := message.MsgNormal{}
	if err = json.Unmarshal(body, &b); err != nil {
		return
	}

	if b.Details["project_id"] == "" || b.Details["training_id"] == "" {
		err = errors.New("invalid message of training")

		return
	}

	v := domain.TrainingIndex{}
	if v.Project.Owner, err = domain.NewAccount(b.Details["project_owner"]); err != nil {
		return
	}

	v.Project.Id = b.Details["project_id"]
	v.TrainingId = b.Details["training_id"]

	return c.handleEventTrainCreated(&v)
}

func (c *consumer) handleEventTrainCreated(info *domain.TrainingIndex) error {
	// wait for the sync of model and dataset
	time.Sleep(10 * time.Second)

	return c.retry(
		func(lastChance bool) error {
			retry, err := c.s.CreateTrainingJob(
				info, c.cfg.TrainingEndpoint, lastChance,
			)
			if err != nil {
				c.log.Errorf(
					"handle training(%s/%s/%s) failed, err:%s",
					info.Project.Owner.Account(), info.Project.Id,
					info.TrainingId, err.Error(),
				)

				if !retry {
					return nil
				}
			}

			return err
		},
		10*time.Second,
	)
}

func (c *consumer) retry(f func(bool) error, interval time.Duration) (err error) {
	n := c.cfg.MaxRetry - 1

	if err = f(n <= 0); err == nil || n <= 0 {
		return
	}

	for i := 1; i < n; i++ {
		time.Sleep(interval)

		if err = f(false); err == nil {
			return
		}
	}

	time.Sleep(interval)

	return f(true)
}

type TrainingConfig struct {
	MaxRetry         int    `json:"max_retry"         required:"true"`
	TrainingEndpoint string `json:"training_endpoint" required:"true"`

	Topics TopicConfig `json:"topics"`
}

type TopicConfig struct {
	TrainingCreated string `json:"training_created" required:"true"`
}
