package messagequeue

import (
	"encoding/json"

	kfk "github.com/opensourceways/kafka-lib/agent"

	"github.com/opensourceways/xihe-server/common/domain/message"
	commondomain "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/points/app"
)

const (
	group    = "xihe-user-points"
	retryNum = 3
)

func Subscribe(s app.UserPointsAppService, topics []string) error {
	c := &consumer{s}

	return kfk.SubscribeWithStrategyOfRetry(group, c.handle, topics, retryNum)
}

type consumer struct {
	s app.UserPointsAppService
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

	cmd.Type = msg.Type
	cmd.Time = msg.CreatedAt
	cmd.Desc = msg.Desc

	return
}
