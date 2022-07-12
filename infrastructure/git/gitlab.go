package git

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/xanzy/go-gitlab"
)

type UserGitlabClient struct {
	*gitlab.Client
}

func NewUserGitlabClient(username, pswd, host string) (userGitLabClient *UserGitlabClient, err error) {
	userGitLabClient = new(UserGitlabClient)
	userGitLabClient.Client, err = gitlab.NewBasicAuthClient(username, pswd, gitlab.WithBaseURL(host))
	return
}
func (u *UserGitlabClient) CreateProject(name, desc, visibility string, mergeRequestsEnabled, snippetsEnabled bool) error {
	var newGitlabProject domain.GitlabProject
	createOpts, err := newGitlabProject.MakeCreateOpt(name, desc, visibility, mergeRequestsEnabled, snippetsEnabled)
	if err != nil {
		return err
	}
	_, _, err = u.Client.Projects.CreateProject(createOpts)

	return err
}
