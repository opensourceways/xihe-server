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
func (impl *messageAdapter) SendUserRegisterEvent(v *domain.UserRegisterEvent) error {
	cfg := &impl.cfg.UserRegister

	msg := common.MsgNormal{
		Type:      cfg.Name,
		User:      v.Account.Account(),
		Desc:      "Register",
		CreatedAt: utils.Now(),
	}

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

// Bind Email
func (impl *messageAdapter) SendBindEmailEvent(v *domain.UserBindEmailEvent) error {
	cfg := &impl.cfg.BindEmail

	msg := common.MsgNormal{
		Type:      cfg.Name,
		User:      v.Account.Account(),
		Desc:      fmt.Sprintf("Bind Email of %s", v.Email),
		CreatedAt: utils.Now(),
	}

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

// Set AvatarId
func (impl *messageAdapter) SendSetAvatarIdEvent(v *domain.UserSetAvatarIdEvent) error {
	cfg := &impl.cfg.SetAvatarId

	msg := common.MsgNormal{
		Type:      cfg.Name,
		User:      v.Account.Account(),
		Desc:      fmt.Sprintf("Set AvatarId of %s", v.AvatarId),
		CreatedAt: utils.Now(),
	}

	return impl.publisher.Publish(cfg.Topic, &msg, nil)
}

// Set Bio
func (impl *messageAdapter) SendSetBioEvent(v *domain.UserSetBioEvent) error {
	cfg := &impl.cfg.SetBio

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
	UserRegister common.TopicConfig `json:"user_register"`
	BindEmail    common.TopicConfig `json:"bind_email"`
	SetAvatarId  common.TopicConfig `json:"set-avatar-id"`
	SetBio       common.TopicConfig `json:"set_bio"`
}
