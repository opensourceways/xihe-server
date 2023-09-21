package messagequeue

import (
	"encoding/json"

	"github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/user/app"
	"github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/user/infrastructure/messageadapter"
)

const (
	retryNum = 3

	handleNameUserFollowingAdd    = "user_following_add"
	handleNameUserFollowingRemove = "user_following_remove"
)

func Subscribe(s app.UserService, subscriber message.Subscriber, topics *messageadapter.Config) (err error) {
	c := &consumer{s}

	err = subscriber.SubscribeWithStrategyOfRetry(
		handleNameUserFollowingAdd,
		c.HandleEventAddFollowing,
		[]string{topics.FollowingAdded.Topic},
		retryNum,
	)

	if err != nil {
		return err
	}

	err = subscriber.SubscribeWithStrategyOfRetry(
		handleNameUserFollowingRemove,
		c.HandleEventRemoveFollowing,
		[]string{topics.FollowingRemoved.Topic},
		retryNum,
	)

	return

}

type consumer struct {
	s app.UserService
}

func (c *consumer) HandleEventAddFollowing(body []byte, h map[string]string) (err error) {
	msg := domain.MsgFollowing{}

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

	err = c.s.AddFollower(&v)

	if err != nil {
		_, ok := err.(repository.ErrorDuplicateCreating)
		if ok {
			err = nil
		}
	}

	return

}

func (c *consumer) HandleEventRemoveFollowing(body []byte, h map[string]string) (err error) {
	msg := domain.MsgFollowing{}

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

	v := domain.FollowerInfo{
		User:     user,
		Follower: follower,
	}

	return c.s.RemoveFollower(&v)

}

type TopicConfig struct {
	FollowingAdd    string `json:"following_add"       required:"true"`
	FollowingRemove string `json:"following_remove"    required:"true"`
}
