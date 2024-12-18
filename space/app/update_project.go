package app

import (
	"errors"

	"golang.org/x/xerrors"

	"github.com/opensourceways/xihe-audit-sync-sdk/audit"
	auditapi "github.com/opensourceways/xihe-audit-sync-sdk/audit/api"
	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/common/domain/allerror"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	spacedomain "github.com/opensourceways/xihe-server/space/domain"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

// only admin user can update level of project
type ProjectUpdateCmd struct {
	Name     domain.ResourceName
	Desc     domain.ResourceDesc
	Title    domain.ResourceTitle
	RepoType domain.RepoType
	CoverId  domain.CoverId
}

func (cmd *ProjectUpdateCmd) toProject(
	p *spacedomain.ProjectModifiableProperty, repo *platform.RepoOption,
) (b bool) {
	f := func() {
		if !b {
			b = true
		}
	}

	if cmd.Name != nil && p.Name.ResourceName() != cmd.Name.ResourceName() {
		p.Name = cmd.Name
		repo.Name = cmd.Name
		f()
	}

	if cmd.Desc != nil && !domain.IsSameDomainValue(cmd.Desc, p.Desc) {
		p.Desc = cmd.Desc
		f()
	}

	if cmd.Title != nil && !domain.IsSameDomainValue(cmd.Title, p.Title) {
		p.Title = cmd.Title
		f()
	}

	if cmd.RepoType != nil && p.RepoType.RepoType() != cmd.RepoType.RepoType() {
		p.RepoType = cmd.RepoType
		repo.RepoType = cmd.RepoType
		f()
	}

	if cmd.CoverId != nil && p.CoverId.CoverId() != cmd.CoverId.CoverId() {
		p.CoverId = cmd.CoverId
		f()
	}

	return
}

// the step1 must be done before step2.
// For example, it can't set the project's name to the one existing.
// gitlab will help to avoid this case.
func (s projectService) Update(
	p *spacedomain.Project, cmd *ProjectUpdateCmd, pr platform.Repository,
) (dto ProjectDTO, err error) {
	//sdk text audit
	var resp audit.ModerationDTO

	title := cmd.Title.ResourceTitle()
	if title != "" {
		resp, _, err = auditapi.Text(title, "title")
		if err != nil {
			return ProjectDTO{}, allerror.New(
				allerror.ErrorCodeFailToModerate,
				resp.Result, err)
		} else if resp.Result != "pass" {
			e := xerrors.Errorf("moderate unpass")
			return ProjectDTO{}, allerror.New(
				allerror.ErrorCodeModerateUnpass,
				resp.Result, e)
		}
	}
	desc := cmd.Desc.ResourceDesc()
	if desc != "" {
		resp, _, err = auditapi.Text(desc, "profile")
		if err != nil {
			return ProjectDTO{}, allerror.New(
				allerror.ErrorCodeFailToModerate,
				resp.Result, err)
		} else if resp.Result != "pass" {
			e := xerrors.Errorf("moderate unpass")
			return ProjectDTO{}, allerror.New(
				allerror.ErrorCodeModerateUnpass,
				resp.Result, e)
		}
	}

	opt := new(platform.RepoOption)
	if !cmd.toProject(&p.ProjectModifiableProperty, opt) {
		s.toProjectDTO(p, &dto)

		return
	}

	// step1
	if opt.IsNotEmpty() {
		if err = pr.Update(p.RepoId, opt); err != nil {
			return
		}
	}

	// step2
	info := spacerepo.ProjectPropertyUpdateInfo{
		ResourceToUpdate: s.toResourceToUpdate(p),
		Property:         p.ProjectModifiableProperty,
	}
	if err = s.repo.UpdateProperty(&info); err != nil {
		return
	}

	s.toProjectDTO(p, &dto)

	return
}

func (s projectService) SetTags(p *spacedomain.Project, cmd *app.ResourceTagsUpdateCmd) error {
	tags, b := cmd.ToTags(p.ProjectModifiableProperty.Tags)
	if !b {
		return nil
	}

	p.ProjectModifiableProperty.Tags = tags
	p.ProjectModifiableProperty.TagKinds = cmd.GenTagKinds(tags)

	info := spacerepo.ProjectPropertyUpdateInfo{
		ResourceToUpdate: s.toResourceToUpdate(p),
		Property:         p.ProjectModifiableProperty,
	}

	return s.repo.UpdateProperty(&info)
}

func (s projectService) AddRelatedModel(
	p *spacedomain.Project, index *domain.ResourceIndex,
) error {
	return s.addRelatedResource(
		p, p.RelatedModels, index, domain.ResourceTypeModel,
		s.repo.AddRelatedModel,
	)
}

func (s projectService) AddRelatedDataset(
	p *spacedomain.Project, index *domain.ResourceIndex,
) error {
	return s.addRelatedResource(
		p, p.RelatedDatasets, index, domain.ResourceTypeDataset,
		s.repo.AddRelatedDataset,
	)
}

func (s projectService) addRelatedResource(
	p *spacedomain.Project, v domain.RelatedResources,
	index *domain.ResourceIndex, t domain.ResourceType,
	f func(*repository.RelatedResourceInfo) error,
) error {
	if v.Has(index) {
		return nil
	}

	if v.Count()+1 > p.MaxRelatedResourceNum() {
		return ExceedMaxRelatedResourceNumError{
			errors.New("exceed max related reousrce num"),
		}
	}

	info := repository.RelatedResourceInfo{
		ResourceToUpdate: s.toResourceToUpdate(p),
		RelatedResource:  *index,
	}

	if err := f(&info); err != nil {
		return err
	}

	_ = s.sender.AddRelatedResource(&message.RelatedResource{
		Promoter: &domain.ResourceObject{
			ResourceIndex: domain.ResourceIndex{
				Owner: p.Owner,
				Id:    p.Id,
			},
			Type: domain.ResourceTypeProject,
		},
		Resource: &domain.ResourceObject{
			ResourceIndex: *index,
			Type:          t,
		},
	})

	return nil
}

func (s projectService) RemoveRelatedModel(
	p *spacedomain.Project, index *domain.ResourceIndex,
) error {
	return s.removeRelatedResource(
		p, p.RelatedModels, index, domain.ResourceTypeModel,
		s.repo.RemoveRelatedModel,
	)
}

func (s projectService) RemoveRelatedDataset(
	p *spacedomain.Project, index *domain.ResourceIndex,
) error {
	return s.removeRelatedResource(
		p, p.RelatedDatasets, index, domain.ResourceTypeDataset,
		s.repo.RemoveRelatedDataset,
	)
}

func (s projectService) removeRelatedResource(
	p *spacedomain.Project, v domain.RelatedResources,
	index *domain.ResourceIndex, t domain.ResourceType,
	f func(*repository.RelatedResourceInfo) error,
) error {
	if !v.Has(index) {
		return nil
	}

	info := repository.RelatedResourceInfo{
		ResourceToUpdate: s.toResourceToUpdate(p),
		RelatedResource:  *index,
	}

	if err := f(&info); err != nil {
		return err
	}

	_ = s.sender.RemoveRelatedResource(&message.RelatedResource{
		Promoter: &domain.ResourceObject{
			ResourceIndex: domain.ResourceIndex{
				Owner: p.Owner,
				Id:    p.Id,
			},
			Type: domain.ResourceTypeProject,
		},
		Resource: &domain.ResourceObject{
			ResourceIndex: *index,
			Type:          t,
		},
	})

	return nil
}

func (s projectService) toResourceToUpdate(p *spacedomain.Project) repository.ResourceToUpdate {
	return repository.ResourceToUpdate{
		Owner:     p.Owner,
		Id:        p.Id,
		Version:   p.Version,
		UpdatedAt: utils.Now(),
	}
}
