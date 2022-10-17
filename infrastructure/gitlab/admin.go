package gitlab

import (
	"strconv"
	"strings"

	sdk "github.com/xanzy/go-gitlab"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
)

var (
	admin     *administrator
	obsHelper *obsService

	endpoint      string
	defaultBranch string
)

func NewUserSerivce() platform.User {
	return admin
}

func Init(cfg *Config) error {
	v, err := sdk.NewClient(cfg.RootToken, sdk.WithBaseURL(cfg.Endpoint))
	if err != nil {
		return err
	}

	admin = &administrator{v}
	endpoint = strings.TrimSuffix(cfg.Endpoint, "/")
	defaultBranch = cfg.DefaultBranch

	return nil
}

type administrator struct {
	cli *sdk.Client
}

func (m *administrator) New(u platform.UserOption) (r domain.PlatformUser, err error) {
	name := u.Name.Account()
	email := u.Email.Email()
	pass := u.Password.Password()
	b := true

	v, _, err := m.cli.Users.CreateUser(&sdk.CreateUserOptions{
		Name:             &name,
		Email:            &email,
		Username:         &name,
		Password:         &pass,
		SkipConfirmation: &b,
	})

	if err != nil {
		return
	}

	r.Id = strconv.Itoa(v.ID)
	r.NamespaceId = strconv.Itoa(v.NamespaceID)

	return
}

func (m *administrator) NewToken(u domain.PlatformUser) (string, error) {
	uid, err := strconv.Atoi(u.Id)
	if err != nil {
		return "", err
	}

	name := "___"
	scope := []string{"api"}

	v, _, err := m.cli.Users.CreatePersonalAccessToken(
		uid, &sdk.CreatePersonalAccessTokenOptions{
			Name:   &name,
			Scopes: &scope,
		},
	)

	if err != nil {
		return "", err
	}

	return v.Token, nil
}
