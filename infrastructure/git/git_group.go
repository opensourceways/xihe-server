package git

import (
	"github.com/xanzy/go-gitlab"
)

type GitGroup struct {
	*gitlab.Client
}

func (u *GitGroup) CreateGroup(name, desc, visibility string, mergeRequestsEnabled, snippetsEnabled bool) error {
	var opts gitlab.CreateGroupOptions
	opts.Name = &name
	opts.Description = &desc
	temp := gitlab.VisibilityValue(visibility)
	opts.Visibility = &temp
	_, _, err := u.Client.Groups.CreateGroup(&opts)

	return err
}
