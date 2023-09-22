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

	handleNameFollowingAdded   = "following_added"
	handleNameFollowingRemoved = "following_removed"
)

func Subscribe(s app.UserService, subscriber message.Subscriber, topics *TopicConfig) (err error) {
	c := &consumer{s}

	err = subscriber.SubscribeWithStrategyOfRetry(
		handleNameFollowingAdded,
		c.handleFollowingAddedEvent,
		[]string{topics.FollowingAdded},
		retryNum,
	)
	if err != nil {
		return err
	}

	err = subscriber.SubscribeWithStrategyOfRetry(
		handleNameFollowingRemoved,
		c.handleFollowingRemovedEvent,
		[]string{topics.FollowingRemoved},
		retryNum,
	)

	return
}

type msgFollowing struct {
	common.MsgNormal

	Follower string `json:"follower"`
}

func followerInfo(body []byte) (info domain.FollowerInfo, err error) {
	msg := msgFollowing{}
	if err = json.Unmarshal(body, &msg); err != nil {
		return
	}

	if info.User, err = domain.NewAccount(msg.MsgNormal.User); err != nil {
		return
	}

	info.Follower, err = domain.NewAccount(msg.Follower)

	return
}

type consumer struct {
	s app.UserService
}

func (c *consumer) handleFollowingAddedEvent(body []byte, h map[string]string) error {
	info, err := followerInfo(body)
	if err != nil {
		return err
	}

	err = c.s.AddFollower(&info)
	if err != nil && repository.IsErrorDuplicateCreating(err) {
		err = nil
	}

	return err
}

func (c *consumer) handleFollowingRemovedEvent(body []byte, h map[string]string) error {
	info, err := followerInfo(body)
	if err != nil {
		return err
	}

	return c.s.RemoveFollower(&info)
}

type TopicConfig struct {
	FollowingAdded   string `json:"following_added"    required:"true"`
	FollowingRemoved string `json:"following_removed"  required:"true"`
}
