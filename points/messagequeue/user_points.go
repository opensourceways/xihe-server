package messagequeue

import (
	"encoding/json"

	"github.com/opensourceways/xihe-server/common/domain/message"
	commondomain "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/points/app"
)

const (
	group    = "xihe-user-points"
	retryNum = 3
)

func Subscribe(s app.UserPointsAppMessageService, topics []string, subscriber message.Subscriber) error {
	c := &consumer{s}

	return subscriber.SubscribeWithStrategyOfRetry(group, c.handle, topics, retryNum)
}

type consumer struct {
	s app.UserPointsAppMessageService
}

func (c *consumer) handle(body []byte, h map[string]string) error {
	msg := &message.MsgNormal{}

	if err := json.Unmarshal(body, msg); err != nil {
		return err
	}

	cmd, err := toCmd(msg)
	if err != nil {
		// no need retry
		return nil
	}

	return c.s.AddPointsItem(&cmd)
}

func toCmd(msg *message.MsgNormal) (cmd app.CmdToAddPointsItem, err error) {
	if cmd.Account, err = commondomain.NewAccount(msg.User); err != nil {
		return
	}

	cmd.TaskId = msg.Type
	cmd.Time = msg.CreatedAt
	cmd.Desc = msg.Desc

	return
}
