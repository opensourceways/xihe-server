package app

import (
	"github.com/opensourceways/xihe-server/common/domain/allerror"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type DatasetUpdateCmd struct {
	Name     domain.ResourceName
	Desc     domain.ResourceDesc
	Title    domain.ResourceTitle
	RepoType domain.RepoType
}

func (cmd *DatasetUpdateCmd) toDataset(
	p *domain.DatasetModifiableProperty, repo *platform.RepoOption,
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

func (s datasetService) Update(
	d *domain.Dataset, cmd *DatasetUpdateCmd, pr platform.Repository,
) (dto DatasetDTO, err error) {
	//sdk text audit
	title := cmd.Title.ResourceTitle()
	if title != "" {
		if err := s.audit.TextAudit(title, audiTitle); err != nil {
			return DatasetDTO{}, allerror.New(
				allerror.ErrorCodeFailToModerate,
				"", err)
		}
	}

	desc := cmd.Desc.ResourceDesc()
	if desc != "" {
		if err := s.audit.TextAudit(desc, audiProfile); err != nil {
			return DatasetDTO{}, allerror.New(
				allerror.ErrorCodeFailToModerate,
				"", err)
		}
	}
	opt := new(platform.RepoOption)
	if !cmd.toDataset(&d.DatasetModifiableProperty, opt) {
		s.toDatasetDTO(d, &dto)

		return
	}

	if opt.IsNotEmpty() {
		if err = pr.Update(d.RepoId, opt); err != nil {
			return
		}
	}

	info := repository.DatasetPropertyUpdateInfo{
		ResourceToUpdate: s.toResourceToUpdate(d),
		Property:         d.DatasetModifiableProperty,
	}
	if err = s.repo.UpdateProperty(&info); err != nil {
		return
	}

	s.toDatasetDTO(d, &dto)

	return
}

func (s datasetService) SetTags(d *domain.Dataset, cmd *ResourceTagsUpdateCmd) error {
	tags, b := cmd.ToTags(d.DatasetModifiableProperty.Tags)
	if !b {
		return nil
	}

	d.DatasetModifiableProperty.Tags = tags
	d.DatasetModifiableProperty.TagKinds = cmd.GenTagKinds(tags)

	info := repository.DatasetPropertyUpdateInfo{
		ResourceToUpdate: s.toResourceToUpdate(d),
		Property:         d.DatasetModifiableProperty,
	}

	return s.repo.UpdateProperty(&info)
}

func (s datasetService) toResourceToUpdate(d *domain.Dataset) repository.ResourceToUpdate {
	return repository.ResourceToUpdate{
		Owner:     d.Owner,
		Id:        d.Id,
		Version:   d.Version,
		UpdatedAt: utils.Now(),
	}
}
