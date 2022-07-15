package git

import (
	"github.com/xanzy/go-gitlab"
)

type GitUserClient struct {
	*GitLabClient
}

func NewGitUserClient(gitlabclient *GitLabClient) (client *GitUserClient) {
	client = new(GitUserClient)
	client.GitLabClient = gitlabclient
	return
}

func (u *GitUserClient) CreateUser(username, name, email, pswd, bio string, isAdmin bool) (user *gitlab.User, err error) {
	var opts gitlab.CreateUserOptions
	opts.Admin = &isAdmin
	opts.Username = &username
	opts.Name = &name
	opts.Email = &email
	opts.Password = &pswd
	opts.Bio = &bio
	var skip bool = true
	opts.SkipConfirmation = &skip
	user, _, err = u.Client.Users.CreateUser(&opts)
	return
}

func (u *GitUserClient) UpdateUser(userid int, username, name, email, bio string) (user *gitlab.User, err error) {
	var opts gitlab.GetUsersOptions
	user, _, err = u.Client.Users.GetUser(userid, opts)
	if err != nil {
		return
	}
	var updateOpts gitlab.ModifyUserOptions
	updateOpts.Username = &username
	updateOpts.Name = &name
	updateOpts.Email = &email
	updateOpts.Bio = &bio
	user, _, err = u.Client.Users.ModifyUser(userid, &updateOpts)
	return
}
