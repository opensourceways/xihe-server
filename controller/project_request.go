package controller

import (
	"github.com/opensourceways/xihe-server/domain"
	spaceapp "github.com/opensourceways/xihe-server/space/app"
	"k8s.io/apimachinery/pkg/util/sets"
)

type projectCreateRequest struct {
	Owner    string   `json:"owner" required:"true"`
	Name     string   `json:"name" required:"true"`
	Desc     string   `json:"desc"`
	Type     string   `json:"type" required:"true"`
	CoverId  string   `json:"cover_id" required:"true"`
	Protocol string   `json:"protocol" required:"true"`
	Training string   `json:"training" required:"true"`
	RepoType string   `json:"repo_type" required:"true"`
	Title    string   `json:"title"`
	Tags     []string `json:"tags"`
}

func (p *projectCreateRequest) toCmd(
	validTags []domain.DomainTags,
) (cmd spaceapp.ProjectCreateCmd, err error) {
	if cmd.Owner, err = domain.NewAccount(p.Owner); err != nil {
		return
	}

	if cmd.Name, err = domain.NewResourceName(p.Name); err != nil {
		return
	}

	if cmd.Type, err = domain.NewProjType(p.Type); err != nil {
		return
	}

	if cmd.Desc, err = domain.NewResourceDesc(p.Desc); err != nil {
		return
	}

	if cmd.CoverId, err = domain.NewCoverId(p.CoverId); err != nil {
		return
	}

	if cmd.Protocol, err = domain.NewProtocolName(p.Protocol); err != nil {
		return
	}

	if cmd.RepoType, err = domain.NewRepoType(p.RepoType); err != nil {
		return
	}

	if cmd.Training, err = domain.NewTrainingPlatform(p.Training); err != nil {
		return
	}

	if p.Title == "" {
		p.Title = p.Name
	}

	if cmd.Title, err = domain.NewResourceTitle(p.Title); err != nil {
		return
	}

	tags := sets.NewString()
	for i := range validTags {
		for _, item := range validTags[i].Items {
			tags.Insert(item.Items...)
		}
	}

	if len(p.Tags) > 0 && !tags.HasAll(p.Tags...) {
		return
	}
	cmd.Tags = p.Tags
	cmd.All = validTags

	err = cmd.Validate()

	return
}

type projectUpdateRequest struct {
	Name     *string `json:"name"`
	Desc     *string `json:"desc"`
	Title    *string `json:"title"`
	RepoType *string `json:"type"`
	CoverId  *string `json:"cover_id"`
}

func (p *projectUpdateRequest) toCmd() (cmd spaceapp.ProjectUpdateCmd, err error) {
	if p.Name != nil {
		if cmd.Name, err = domain.NewResourceName(*p.Name); err != nil {
			return
		}
	}

	if p.Desc != nil {
		if cmd.Desc, err = domain.NewResourceDesc(*p.Desc); err != nil {
			return
		}
	}

	if p.Title != nil {
		if cmd.Title, err = domain.NewResourceTitle(*p.Title); err != nil {
			return
		}
	}

	if p.RepoType != nil {
		if cmd.RepoType, err = domain.NewRepoType(*p.RepoType); err != nil {
			return
		}
	}

	if p.CoverId != nil {
		if cmd.CoverId, err = domain.NewCoverId(*p.CoverId); err != nil {
			return
		}
	}

	return
}

type projectDetail struct {
	Liked    bool   `json:"liked"`
	AvatarId string `json:"avatar_id"`

	*spaceapp.ProjectDetailDTO
}

type projectsInfo struct {
	Owner    string `json:"owner"`
	AvatarId string `json:"avatar_id"`

	*spaceapp.ProjectsDTO
}

type projectForkRequest struct {
	Name string `json:"name" required:"true"`
	Desc string `json:"desc"`
}

func (p *projectForkRequest) toCmd() (cmd spaceapp.ProjectForkCmd, err error) {
	if cmd.Name, err = domain.NewResourceName(p.Name); err != nil {
		return
	}

	cmd.Desc, err = domain.NewResourceDesc(p.Desc)

	return
}

type canApplyResourceNameResp struct {
	CanApply bool `json:"can_apply"`
}
