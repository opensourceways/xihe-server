package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
)

type ProjectForkCmd struct {
	Name      domain.ResourceName
	Desc      domain.ResourceDesc
	From      domain.Project
	Owner     domain.Account
	ValidTags []domain.DomainTags
}

func (cmd *ProjectForkCmd) toProject(r *domain.Project) {
	p := &cmd.From
	*r = domain.Project{
		Owner:     cmd.Owner,
		Type:      p.Type,
		Protocol:  p.Protocol,
		Training:  p.Training,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		ProjectModifiableProperty: domain.ProjectModifiableProperty{
			Name:     cmd.Name,
			Desc:     p.Desc,
			CoverId:  p.CoverId,
			RepoType: p.RepoType,
			Tags:     p.Tags,
		},
	}

	h := ResourceTagsUpdateCmd{
		All: cmd.ValidTags,
	}

	r.TagKinds = h.genTagKinds(p.Tags)

	if cmd.Desc != nil {
		r.Desc = cmd.Desc
	}
}

func (s projectService) Fork(cmd *ProjectForkCmd, pr platform.Repository) (dto ProjectDTO, err error) {
	pid, err := pr.Fork(cmd.From.RepoId, cmd.Name)
	if err != nil {
		return
	}

	v := new(domain.Project)
	cmd.toProject(v)
	v.RepoId = pid

	p, err := s.repo.Save(v)
	if err != nil {
		return
	}

	s.toProjectDTO(&p, &dto)

	// create activity
	r, repoType := p.ResourceObject()
	ua := genActivityForCreatingResource(r, repoType)
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
