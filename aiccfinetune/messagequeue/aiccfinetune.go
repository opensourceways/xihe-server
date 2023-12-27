package messagequeue

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/aiccfinetune/app"
	"github.com/opensourceways/xihe-server/aiccfinetune/domain"
	"github.com/opensourceways/xihe-server/common/domain/message"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	retryNum = 3

	handleNameAICCFinetuneCreated = "aicc_finetune_created_test"
)

func Subscribe(
	cfg AICCFinetuneConfig,
	tcTopic string,
	s app.AICCFinetuneService,
	subscriber message.Subscriber,
) (err error) {
	c := &consumer{cfg: cfg, as: s}

	// aicc finetune created
	err = subscriber.SubscribeWithStrategyOfRetry(
		handleNameAICCFinetuneCreated,
		c.handleEventAICCFinetuneCreated,
		[]string{tcTopic}, retryNum,
	)
	return
}

type consumer struct {
	cfg AICCFinetuneConfig
	as  app.AICCFinetuneService
}

func (c *consumer) handleEventAICCFinetuneCreated(body []byte, h map[string]string) (err error) {
	b := message.MsgNormal{}
	if err = json.Unmarshal(body, &b); err != nil {
		return
	}

	if b.Details["id"] == "" || b.Details["model"] == "" {
		err = errors.New("invalid message of aicc finetune")

		return
	}

	v := domain.AICCFinetuneIndex{}

	v.FinetuneId = b.Details["id"]
	v.Model, err = domain.NewModelName(b.Details["model"])
	if err != nil {
		return
	}

	v.User, err = types.NewAccount(b.User)
	if err != nil {
		return
	}

	return c.createJob(&v)
}

func (c *consumer) createJob(info *domain.AICCFinetuneIndex) error {
	// wait for the sync of model and dataset
	time.Sleep(10 * time.Second)
	utils.RetryThreeTimes(func() error {
		retry, err := c.as.CreateAICCFinetuneJob(
			info, c.cfg.AICCFinetuneEndpoint, true,
		)
		if err != nil {
			logrus.Errorf(
				"handle aicc-finetune(%s %s) failed, err:%s", info.User.Account(),
				info.FinetuneId, err.Error(),
			)

			if !retry {
				return nil
			}
		}

		return err
	})

	return nil
}

type AICCFinetuneConfig struct {
	AICCFinetuneEndpoint string `json:"aiccfinetune_endpoint" required:"true"`
}
