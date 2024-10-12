package domain

import "github.com/opensourceways/xihe-server/domain"

type Project struct {
	Id string

	Owner    domain.Account
	Type     domain.ProjType
	Protocol domain.ProtocolName
	Training domain.TrainingPlatform

	ProjectModifiableProperty

	RepoId string

	RelatedModels   domain.RelatedResources
	RelatedDatasets domain.RelatedResources

	CreatedAt int64
	UpdatedAt int64

	Version int

	// following fields is not under the controlling of version
	LikeCount     int
	ForkCount     int
	DownloadCount int
}

func (p *Project) MaxRelatedResourceNum() int {
	return domain.DomainConfig.MaxRelatedResourceNum
}

func (p *Project) IsPrivate() bool {
	return p.RepoType.RepoType() == domain.RepoTypePrivate
}

func (p *Project) IsOnline() bool {
	return p.RepoType.RepoType() == domain.RepoTypeOnline
}

func (p *Project) ResourceIndex() domain.ResourceIndex {
	return domain.ResourceIndex{
		Owner: p.Owner,
		Id:    p.Id,
	}
}

func (p *Project) ResourceObject() (domain.ResourceObject, domain.RepoType) {
	return domain.ResourceObject{
		Type:          domain.ResourceTypeProject,
		ResourceIndex: p.ResourceIndex(),
	}, p.RepoType
}

func (p *Project) RelatedResources() []domain.ResourceObjects {
	r := make([]domain.ResourceObjects, 0, 2)

	if len(p.RelatedModels) > 0 {
		r = append(r, domain.ResourceObjects{
			Type:    domain.ResourceTypeModel,
			Objects: p.RelatedModels,
		})
	}

	if len(p.RelatedDatasets) > 0 {
		r = append(r, domain.ResourceObjects{
			Type:    domain.ResourceTypeDataset,
			Objects: p.RelatedDatasets,
		})
	}

	return r
}

type ProjectModifiableProperty struct {
	Name     domain.ResourceName
	Desc     domain.ResourceDesc
	Title    domain.ResourceTitle
	CoverId  domain.CoverId
	RepoType domain.RepoType
	Tags     []string
	TagKinds []string
	Level    domain.ResourceLevel
}

type ProjectSummary struct {
	Id            string
	Owner         domain.Account
	Name          domain.ResourceName
	Desc          domain.ResourceDesc
	Title         domain.ResourceTitle
	Level         domain.ResourceLevel
	CoverId       domain.CoverId
	Tags          []string
	UpdatedAt     int64
	LikeCount     int
	ForkCount     int
	DownloadCount int
}
