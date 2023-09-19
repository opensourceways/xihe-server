package messagequeue

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	retryNum = 3

	handleNameTrainingCreated = "training_created"
)

func Subscribe(
	cfg TrainingConfig,
	tcTopic string,
	s app.TrainingService,
	subscriber message.Subscriber,
) (err error) {
	c := &consumer{cfg: cfg, s: s}

	// training created
	err = subscriber.SubscribeWithStrategyOfRetry(
		handleNameTrainingCreated,
		c.handleEventTrainingCreated,
		[]string{tcTopic}, retryNum,
	)

	return
}

type consumer struct {
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

	return c.createJob(&v)
}

func (c *consumer) createJob(info *domain.TrainingIndex) error {
	// wait for the sync of model and dataset
	time.Sleep(10 * time.Second)

	utils.RetryThreeTimes(func() error {
		retry, err := c.s.CreateTrainingJob(
			info, c.cfg.TrainingEndpoint, true,
		)
		if err != nil {
			logrus.Errorf(
				"handle training(%s/%s/%s) failed, err:%s",
				info.Project.Owner.Account(), info.Project.Id,
				info.TrainingId, err.Error(),
			)

			if !retry {
				return nil
			}
		}

		return err
	})

	return nil
}

type TrainingConfig struct {
	MaxRetry         int    `json:"max_retry"`
	TrainingEndpoint string `json:"training_endpoint" required:"true"`
}

func (cfg *TrainingConfig) SetDefault() {
	if cfg.MaxRetry == 0 {
		cfg.MaxRetry = 3
	}
}
