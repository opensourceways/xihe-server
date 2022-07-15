package domain

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

type GitlabProject struct {
	Name                 string
	Description          string
	MergeRequestsEnabled bool
	SnippetsEnabled      bool
	Visibility           bool
}

func (p *GitlabProject) MakeCreateOpt(name, desc, visibility string, mergeRequestsEnabled, snippetsEnabled bool) (createProjectOptions *gitlab.CreateProjectOptions, err error) {
	createProjectOptions = new(gitlab.CreateProjectOptions)
	switch visibility {
	case string(gitlab.PrivateVisibility):
		temp := gitlab.VisibilityValue(gitlab.PrivateVisibility)
		createProjectOptions.Visibility = &temp
	case string(gitlab.InternalVisibility):
		temp := gitlab.VisibilityValue(gitlab.InternalVisibility)
		createProjectOptions.Visibility = &temp
	case string(gitlab.PublicVisibility):
		temp := gitlab.VisibilityValue(gitlab.PublicVisibility)
		createProjectOptions.Visibility = &temp
	default:
		return nil, fmt.Errorf("visibility must be set as one of ( %s,%s,%s)", gitlab.PrivateVisibility, gitlab.InternalVisibility, gitlab.PublicVisibility)
	}
	if len(name) == 0 {
		return nil, fmt.Errorf("project Name must be fill")
	}
	if len(desc) == 0 {
		return nil, fmt.Errorf("project Description must be fill")
	}
	createProjectOptions.Name = &name
	createProjectOptions.Description = &desc
	createProjectOptions.MergeRequestsEnabled = &mergeRequestsEnabled
	createProjectOptions.SnippetsEnabled = &snippetsEnabled
	return
}
