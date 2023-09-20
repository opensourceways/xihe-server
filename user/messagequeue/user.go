package messagequeue

import (
	"encoding/json"

	"github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/user/app"
	"github.com/opensourceways/xihe-server/user/domain"
)

const (
	retryNum = 3

	handleNameUserFollowingAdd    = "user_following_add"
	handleNameUserFollowingRemove = "user_following_remove"
)

func Subscribe(s app.UserService, subscriber message.Subscriber) (err error) {
	c := &consumer{s}

	topicmsg := TopicConfig{}

	err = subscriber.SubscribeWithStrategyOfRetry(
		handleNameUserFollowingAdd,
		c.HandleEventAddFollowing,
		[]string{topicmsg.FollowingAdd},
		retryNum,
	)

	if err != nil {
		return
	}

	err = subscriber.SubscribeWithStrategyOfRetry(
		handleNameUserFollowingRemove,
		c.HandleEventRemoveFollowing,
		[]string{topicmsg.FollowingRemove},
		retryNum,
	)

	return

}

type consumer struct {
	s app.UserService
}

func (c *consumer) HandleEventAddFollowing(body []byte, h map[string]string) error {
	msg := message.MsgNormal{}

	f := domain.FollowerInfo{}

	if err := json.Unmarshal(body, &msg); err != nil {
		return err
	}

	user, err := domain.NewAccount(msg.User)

	if err != nil {
		return nil
	}

	v := domain.FollowerInfo{
		User:     user,
		Follower: f.Follower,
	}

	err = c.s.AddFollower(&v)

	if err != nil {
		_, ok := err.(repository.ErrorDuplicateCreating)
		if ok {
			err = nil
		}
	}

	return c.s.AddFollower(&v)

}

func (c *consumer) HandleEventRemoveFollowing(body []byte, h map[string]string) error {
	msg := message.MsgNormal{}

	f := domain.FollowerInfo{}

	if err := json.Unmarshal(body, &msg); err != nil {
		return err
	}

	user, err := domain.NewAccount(msg.User)

	if err != nil {
		return nil
	}

	v := domain.FollowerInfo{
		User:     user,
		Follower: f.Follower,
	}

	return c.s.RemoveFollower(&v)

}

type TopicConfig struct {
	FollowingAdd    string `json:"following_add"       required:"true"`
	FollowingRemove string `json:"following_remove"    required:"true"`
}
