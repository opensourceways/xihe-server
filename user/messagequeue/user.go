package messagequeue

import (
	"encoding/json"

	kfk "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/xihe-server/user/app"
	"github.com/opensourceways/xihe-server/user/domain"
)

const (
	retryNum = 3

	handleNameUserFollowingAdd    = "user_following_add"
	handleNameUserFollowingRemove = "user_following_remove"
)

func Subscribe(s app.RegService, topics *TopicConfig) (err error) {
	c := &consumer{s}

	if err = kfk.SubscribeWithStrategyOfRetry(
		handleNameUserFollowingAdd,
		c.HandleEventAddFollowing,
		[]string{topics.FollowingAdd}, retryNum); err != nil {
		return
	}
	err = kfk.SubscribeWithStrategyOfRetry(
		handleNameUserFollowingRemove,
		c.HandleEventRemoveFollowing,
		[]string{topics.FollowingRemove}, retryNum)

	return

}

type consumer struct {
	s app.RegService
}

func (c *consumer) HandleEventAddFollowing(body []byte, h map[string]string) error {
	msg := &domain.UserRegInfo{}

	if err := json.Unmarshal(body, msg); err != nil {
		return err
	}

	cmd, err := toCmd(msg)
	if err != nil {
		return nil
	}
	return c.s.UpsertUserRegInfo(&cmd)

}

func (c *consumer) HandleEventRemoveFollowing(body []byte, h map[string]string) error {
	msg := &domain.UserRegInfo{}

	if err := json.Unmarshal(body, msg); err != nil {
		return err
	}
	cmd, err := toCmd(msg)
	if err != nil {
		return nil
	}
	return c.s.UpsertUserRegInfo(&cmd)

}

func toCmd(msg *domain.UserRegInfo) (cmd app.UserRegisterInfoCmd, err error) {

	cmd.Account = msg.Account
	cmd.Email = msg.Email
	cmd.Name = msg.Name
	cmd.Identity = msg.Identity
	cmd.City = msg.City
	cmd.Detail = msg.Detail
	cmd.Phone = msg.Phone
	cmd.Province = msg.Province
	cmd.Version = msg.Version

	return
}

type TopicConfig struct {
	FollowingAdd    string `json:"following_add"`
	FollowingRemove string `json:"following_remove"`
}
