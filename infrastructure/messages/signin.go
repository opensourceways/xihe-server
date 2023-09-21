package messages

import (
	common "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

func NewUserSignedInMessageAdapter(cfg *UserSignedInConfig, p common.Publisher) *signInMessageAdapter {
	return &signInMessageAdapter{cfg: *cfg, publisher: p}
}

type signInMessageAdapter struct {
	cfg       UserSignedInConfig
	publisher common.Publisher
}

func (impl *signInMessageAdapter) SendUserSignedIn(v *domain.UserSignedInEvent) error {
	t := &impl.cfg.UserSignedIn

	return impl.publisher.Publish(
		t.Topic,
		&common.MsgNormal{
			Type:      t.Name,
			User:      v.Account.Account(),
			CreatedAt: utils.Now(),
		},
		nil,
	)
}

type UserSignedInConfig struct {
	UserSignedIn common.TopicConfig `json:"user_signedin" required:"true"`
}
