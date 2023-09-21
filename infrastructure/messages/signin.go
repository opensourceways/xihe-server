package messages

import (
	common "github.com/opensourceways/xihe-server/common/domain/message"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

func NewSignInMessageAdapter(cfg *SignInConfig, p common.Publisher) *signinMessageAdapter {
	return &signinMessageAdapter{cfg: *cfg, publisher: p}
}

type signinMessageAdapter struct {
	cfg       SignInConfig
	publisher common.Publisher
}

func (impl *signinMessageAdapter) SendUserSignedIn(v *domain.UserSignedInEvent) error {
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

type SignInConfig struct {
	UserSignedIn common.TopicConfig `json:"user_signedin" required:"true"`
}
