package messagequeue

import (
	"encoding/json"

	"github.com/opensourceways/xihe-server/common/domain/message"
	common "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/user/app"
	"github.com/opensourceways/xihe-server/user/domain"
)

const (
	retryNum = 3

	handleNameUserFollowingAdd    = "user_following_add"
	handleNameUserFollowingRemove = "user_following_remove"
)

func Subscribe(s app.UserService, subscriber message.Subscriber, topics *TopicConfig) (err error) {
	c := &consumer{s}

	err = subscriber.SubscribeWithStrategyOfRetry(
		handleNameUserFollowingAdd,
		c.HandleEventAddFollowing,
		[]string{topics.FollowingAdded},
		retryNum,
	)
	if err != nil {
		return err
	}

	err = subscriber.SubscribeWithStrategyOfRetry(
		handleNameUserFollowingRemove,
		c.HandleEventRemoveFollowing,
		[]string{topics.FollowingRemoved},
		retryNum,
	)

	return
}

type consumer struct {
	s app.UserService
}

type MsgFollowing struct {
	MsgNormal common.MsgNormal
	Follower  string `json:"follower"`
}

func (c *consumer) HandleEventAddFollowing(body []byte, h map[string]string) (err error) {
	msg := MsgFollowing{}

	if err := json.Unmarshal(body, &msg); err != nil {
		return err
	}

	user, err := domain.NewAccount(msg.MsgNormal.User)
	if err != nil {
		return
	}

	follower, err := domain.NewAccount(msg.Follower)
	if err != nil {
		return
	}

	v := domain.FollowerInfo{
		User:     user,
		Follower: follower,
	}

	if err = c.s.AddFollower(&v); err != nil {
		_, ok := err.(repository.ErrorDuplicateCreating)
		if ok {
			err = nil
		}
	}

	return
}

func (c *consumer) HandleEventRemoveFollowing(body []byte, h map[string]string) (err error) {
	msg := MsgFollowing{}

	if err := json.Unmarshal(body, &msg); err != nil {
		return
	}

	user, err := domain.NewAccount(msg.MsgNormal.User)
	if err != nil {
		return
	}

	follower, err := domain.NewAccount(msg.Follower)
	if err != nil {
		return
	}

	return c.s.RemoveFollower(&domain.FollowerInfo{
		User:     user,
		Follower: follower,
	})
}

type TopicConfig struct {
	FollowingAdded   string `json:"following_added"       required:"true"`
	FollowingRemoved string `json:"following_removed"    required:"true"`
}
