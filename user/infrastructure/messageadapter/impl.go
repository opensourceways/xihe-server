package messageadapter

import (
	"fmt"

	common "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/utils"
	"github.com/opensourceways/xihe-server/user/domain/message"
)

func MessageAdapter(cfg *Config, p common.Publisher) *messageAdapter {
	return &messageAdapter{cfg: *cfg, publisher: p}
}

type messageAdapter struct {
	cfg       Config
	publisher common.Publisher
}

// Register
func (impl *messageAdapter) SendUserSignedUpEvent(v *domain.UserSignedUpEvent) error {
	cfg := &impl.cfg.UserSignedUp

	msg := common.MsgNormal{
		Type:      cfg.Name,
		User:      v.Account.Account(),
		Desc:      "Register",
		CreatedAt: utils.Now(),
	}

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

// Set AvatarId
func (impl *messageAdapter) SendUserAvatarSetEvent(v *domain.UserAvatarSetEvent) error {
	cfg := &impl.cfg.AvatarSet

	msg := common.MsgNormal{
		Type:      cfg.Name,
		User:      v.Account.Account(),
		Desc:      fmt.Sprintf("Set AvatarId of %s", v.AvatarId),
		CreatedAt: utils.Now(),
	}

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

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

// Add following
func (impl *messageAdapter) SendFollowingAddedEvent(v *domain.FollowerInfo) error {
	cfg := &impl.cfg.FollowingAdded

	msg := message.msgFolloweing{
		Follower: v.Follower.Account(),
		com: common.MsgNormal{
			Type:      cfg.Name,
			User:      v.User.Account(),
			Desc:      fmt.Sprintf("Add Following of %s", v.Follower.Account()),
			CreatedAt: utils.Now(),
		},
	}

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

// Remove following
func (impl *messageAdapter) SendFollowingRemovedEvent(v *domain.FollowerInfo) error {
	cfg := &impl.cfg.FollowingRemoved

	msg := message.msgFolloweing{
		Follower: v.Follower.Account(),
		com: common.MsgNormal{
			Type:      cfg.Name,
			User:      v.User.Account(),
			Desc:      fmt.Sprintf("Remove Following of %s", v.Follower.Account()),
			CreatedAt: utils.Now(),
		},
	}

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

// Config
type Config struct {
	UserSignedUp common.TopicConfig `json:"sign_up"           required:"true"`
	BioSet       common.TopicConfig `json:"set_bio"           required:"true"`
	AvatarSet    common.TopicConfig `json:"set_avatar"        required:"true"`

	FollowingAdded   common.TopicConfig `json:"following_added"   required:"true"`
	FollowingRemoved common.TopicConfig `json:"following_removed" required:"true"`
}
