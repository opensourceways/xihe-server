package app

import (
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
	spacedomain "github.com/opensourceways/xihe-server/space/domain"
)

type ProjectForkCmd struct {
	Name      domain.ResourceName
	Desc      domain.ResourceDesc
	From      spacedomain.Project
	Owner     domain.Account
	ValidTags []domain.DomainTags
}

func (cmd *ProjectForkCmd) toProject(r *spacedomain.Project) {
	p := &cmd.From
	*r = spacedomain.Project{
		Owner:     cmd.Owner,
		Type:      p.Type,
		Protocol:  p.Protocol,
		Training:  p.Training,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		ProjectModifiableProperty: spacedomain.ProjectModifiableProperty{
			Name:     cmd.Name,
			Desc:     p.Desc,
			CoverId:  p.CoverId,
			RepoType: p.RepoType,
			Tags:     p.Tags,
		},
	}

	h := app.ResourceTagsUpdateCmd{
		All: cmd.ValidTags,
	}

	r.TagKinds = h.GenTagKinds(p.Tags)

	if cmd.Desc != nil {
		r.Desc = cmd.Desc
	}
}

func (s projectService) Fork(cmd *ProjectForkCmd, pr platform.Repository) (dto ProjectDTO, err error) {
	pid, err := pr.Fork(cmd.From.RepoId, cmd.Name)
	if err != nil {
		return
	}

	v := new(spacedomain.Project)
	cmd.toProject(v)
	v.RepoId = pid

	p, err := s.repo.Save(v)
	if err != nil {
		return
	}

	s.toProjectDTO(&p, &dto)

	// create activity
	r, repoType := p.ResourceObject()
	ua := app.GenActivityForCreatingResource(r, repoType)
	ua.Type = domain.ActivityTypeFork
	_ = s.activity.Save(&ua)

	// send event
	_ = s.sender.IncreaseFork(&domain.ResourceIndex{
		Owner: cmd.From.Owner,
		Id:    cmd.From.Id,
	})

	_ = s.sender.AddOperateLogForCreateResource(r, p.Name)

	return
}
