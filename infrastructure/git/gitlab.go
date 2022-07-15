package git

import (
	"github.com/opensourceways/xihe-server/config"
	"github.com/xanzy/go-gitlab"
)

type GitLabClient struct {
	*gitlab.Client
}

func NewGitlabClient(cfg *config.Config) (userGitLabClient *GitLabClient, err error) {
	userGitLabClient = new(GitLabClient)
	userGitLabClient.Client, err = gitlab.NewBasicAuthClient(cfg.Gitlab.RootUser, cfg.Gitlab.RootPswd, gitlab.WithBaseURL(cfg.Gitlab.Host))
	return
}
