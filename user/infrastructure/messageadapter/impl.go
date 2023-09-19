package messageadapter

import (
	"fmt"

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

// Config
type Config struct {
	UserSignedUp common.TopicConfig `json:"user-signed-up"`
	BioSet       common.TopicConfig `json:"bio_set"`
	AvatarSet    common.TopicConfig `json:"avatar_set"`
}
