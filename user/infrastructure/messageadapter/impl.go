package messageadapter

import (
	"fmt"

	"github.com/sirupsen/logrus"

	common "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/utils"
)

func MessageAdapter(cfg *Config, p common.Publisher) *messageAdapter {
	return &messageAdapter{cfg: *cfg, publisher: p}
}

type messageAdapter struct {
	cfg       Config
	publisher common.Publisher
}

type msgFollowing struct {
	common.MsgNormal

	Follower string `json:"follower"`
}

// Sign Up
func (impl *messageAdapter) SendUserSignedUpEvent(v *domain.UserSignedUpEvent) error {
	cfg := &impl.cfg.UserSignedUp

	msg := common.MsgNormal{
		Type:      cfg.Name,
		User:      v.Account.Account(),
		Desc:      "Sign Up",
		CreatedAt: utils.Now(),
	}

	logrus.Debugf("Send sign up msg: %+v", msg)

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

// Set Avatar
func (impl *messageAdapter) SendUserAvatarSetEvent(v *domain.UserAvatarSetEvent) error {
	cfg := &impl.cfg.AvatarSet

	msg := common.MsgNormal{
		Type:      cfg.Name,
		User:      v.Account.Account(),
		Desc:      fmt.Sprintf("Set Avatar of %s", v.AvatarId),
		CreatedAt: utils.Now(),
	}

	logrus.Debugf("Send set avatar  msg: %+v", msg)

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

// Set Bio
func (impl *messageAdapter) SendUserBioSetEvent(v *domain.UserBioSetEvent) error {
	cfg := &impl.cfg.BioSet

	msg := common.MsgNormal{
		Type:      cfg.Name,
		User:      v.Account.Account(),
		Desc:      fmt.Sprintf("Set Bio of %s", v.Bio),
		CreatedAt: utils.Now(),
	}

	logrus.Debugf("Send set biomsg: %+v", msg)

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

// Add following
func (impl *messageAdapter) SendFollowingAddedEvent(v *domain.FollowerInfo) error {
	cfg := &impl.cfg.FollowingAdded

	msg := msgFollowing{
		Follower: v.Follower.Account(),
		MsgNormal: common.MsgNormal{
			Type:      cfg.Name,
			User:      v.User.Account(),
			Desc:      fmt.Sprintf("Add Following of %s", v.Follower.Account()),
			CreatedAt: utils.Now(),
		},
	}

	logrus.Debugf("Send add follower msg: %+v", msg)

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

// Remove following
func (impl *messageAdapter) SendFollowingRemovedEvent(v *domain.FollowerInfo) error {
	cfg := &impl.cfg.FollowingRemoved

	msg := msgFollowing{
		Follower: v.Follower.Account(),
		MsgNormal: common.MsgNormal{
			Type:      cfg.Name,
			User:      v.User.Account(),
			Desc:      fmt.Sprintf("Remove Following of %s", v.Follower.Account()),
			CreatedAt: utils.Now(),
		},
	}

	logrus.Debugf("Send remove followermsg: %+v", msg)

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

// operate log
func (impl *messageAdapter) AddOperateLogForNewUser(u domain.Account) error {
	a := ""
	if u != nil {
		a = u.Account()
	}

	msg := common.MsgOperateLog{
		When: utils.Now(),
		User: a,
		Type: "user",
	}

	logrus.Debugf("Send new user oprate msg: %+v", msg)

	return impl.publisher.Publish(impl.cfg.OperateLog, &msg, nil)
}

// Config
type Config struct {
	BioSet           common.TopicConfig `json:"bio_set"           required:"true"`
	AvatarSet        common.TopicConfig `json:"avatar_set"        required:"true"`
	OperateLog       string             `json:"operate_log"       required:"true"`
	UserSignedUp     common.TopicConfig `json:"user_signedup"     required:"true"`
	FollowingAdded   common.TopicConfig `json:"following_added"   required:"true"`
	FollowingRemoved common.TopicConfig `json:"following_removed" required:"true"`
}
