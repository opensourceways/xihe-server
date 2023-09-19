package messagequeue

import (
	"encoding/json"

	kfk "github.com/opensourceways/kafka-lib/agent"
	"github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/user/app"
	"github.com/opensourceways/xihe-server/user/domain"
)

const (
	retryNum = 3

	handleNameUserFollowingAdd    = "user_following_add"
	handleNameUserFollowingRemove = "user_following_remove"
)

func Subscribe(s app.UserService, topics *TopicConfig) (err error) {
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
	s app.UserService
}

func (c *consumer) HandleEventAddFollowing(body []byte, h map[string]string) error {
	msg := message.MsgNormal{}

	if err := json.Unmarshal(body, msg); err != nil {
		return err
	}
	user, err := domain.NewAccount(msg.User)
	if err != nil {
		return nil
	}

	v := domain.FollowerInfo{
		User:     user,
		Follower: user,
	}
	return c.s.AddFollower(&v)

}

func (c *consumer) HandleEventRemoveFollowing(body []byte, h map[string]string) error {
	msg := message.MsgNormal{}

	if err := json.Unmarshal(body, msg); err != nil {
		return err
	}

	user, err := domain.NewAccount(msg.User)
	if err != nil {
		return nil
	}

	v := domain.FollowerInfo{
		User:     user,
		Follower: user,
	}
	return c.s.RemoveFollower(&v)

}

type TopicConfig struct {
	FollowingAdd    string `json:"following_add"`
	FollowingRemove string `json:"following_remove"`
}
