package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type ModelUpdateCmd struct {
	Name     domain.ResourceName
	Desc     domain.ResourceDesc
	Title    domain.ResourceTitle
	RepoType domain.RepoType
}

func (cmd *ModelUpdateCmd) toModel(
	p *domain.ModelModifiableProperty, repo *platform.RepoOption,
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

	return
}

func (s modelService) Update(
	m *domain.Model, cmd *ModelUpdateCmd, pr platform.Repository,
) (dto ModelDTO, err error) {
	//sdk text audit
	title := cmd.Title.ResourceTitle()
	if title != "" {
		if err := s.audit.TextAudit(title, auditTitle); err != nil {
			return ModelDTO{}, err
		}
	}

	desc := cmd.Desc.ResourceDesc()
	if desc != "" {
		if err := s.audit.TextAudit(desc, auditProfile); err != nil {
			return ModelDTO{}, err
		}
	}

	opt := new(platform.RepoOption)
	if !cmd.toModel(&m.ModelModifiableProperty, opt) {
		s.toModelDTO(m, &dto)

		return
	}

	if opt.IsNotEmpty() {
		if err = pr.Update(m.RepoId, opt); err != nil {
			return
		}
	}

	info := repository.ModelPropertyUpdateInfo{
		ResourceToUpdate: s.toResourceToUpdate(m),
		Property:         m.ModelModifiableProperty,
	}
	if err = s.repo.UpdateProperty(&info); err != nil {
		return
	}

	s.toModelDTO(m, &dto)

	return
}

func (s modelService) SetTags(m *domain.Model, cmd *ResourceTagsUpdateCmd) error {
	tags, b := cmd.ToTags(m.ModelModifiableProperty.Tags)
	if !b {
		return nil
	}

	m.ModelModifiableProperty.Tags = tags
	m.ModelModifiableProperty.TagKinds = cmd.GenTagKinds(tags)

	info := repository.ModelPropertyUpdateInfo{
		ResourceToUpdate: s.toResourceToUpdate(m),
		Property:         m.ModelModifiableProperty,
	}

	return s.repo.UpdateProperty(&info)
}

func (s modelService) AddRelatedDataset(
	m *domain.Model, index *domain.ResourceIndex,
) error {
	if m.RelatedDatasets.Has(index) {
		return nil
	}

	if m.RelatedDatasets.Count()+1 > m.MaxRelatedResourceNum() {
		return ErrorExceedMaxRelatedResourceNum{
			errors.New("exceed max related reousrce num"),
		}
	}

	info := repository.RelatedResourceInfo{
		ResourceToUpdate: s.toResourceToUpdate(m),
		RelatedResource:  *index,
	}

	if err := s.repo.AddRelatedDataset(&info); err != nil {
		return err
	}

	_ = s.sender.AddRelatedResource(&message.RelatedResource{
		Promoter: &domain.ResourceObject{
			ResourceIndex: domain.ResourceIndex{
				Owner: m.Owner,
				Id:    m.Id,
			},
			Type: domain.ResourceTypeModel,
		},
		Resource: &domain.ResourceObject{
			ResourceIndex: *index,
			Type:          domain.ResourceTypeDataset,
		},
	})

	return nil
}

func (s modelService) RemoveRelatedDataset(
	m *domain.Model, index *domain.ResourceIndex,
) error {
	if !m.RelatedDatasets.Has(index) {
		return nil
	}

	info := repository.RelatedResourceInfo{
		ResourceToUpdate: s.toResourceToUpdate(m),
		RelatedResource:  *index,
	}

	if err := s.repo.RemoveRelatedDataset(&info); err != nil {
		return err
	}

	_ = s.sender.RemoveRelatedResource(&message.RelatedResource{
		Promoter: &domain.ResourceObject{
			ResourceIndex: domain.ResourceIndex{
				Owner: m.Owner,
				Id:    m.Id,
			},
			Type: domain.ResourceTypeModel,
		},
		Resource: &domain.ResourceObject{
			ResourceIndex: *index,
			Type:          domain.ResourceTypeDataset,
		},
	})

	return nil
}

func (s modelService) toResourceToUpdate(m *domain.Model) repository.ResourceToUpdate {
	return repository.ResourceToUpdate{
		Owner:     m.Owner,
		Id:        m.Id,
		Version:   m.Version,
		UpdatedAt: utils.Now(),
	}
}
