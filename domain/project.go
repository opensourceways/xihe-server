package domain

type Project struct {
	Id string

	Owner    Account
	Type     ProjType
	Protocol ProtocolName
	Training TrainingPlatform

	ProjectModifiableProperty

	RepoId string

	RelatedModels   RelatedResources
	RelatedDatasets RelatedResources

	CreatedAt int64
	UpdatedAt int64

	Version int

	// following fileds is not under the controlling of version
	LikeCount     int
	ForkCount     int
	DownloadCount int
}

func (p *Project) MaxRelatedResourceNum() int {
	return DomainConfig.MaxRelatedResourceNum
}

func (p *Project) IsPrivate() bool {
	return p.RepoType.RepoType() == RepoTypePrivate
}

func (p *Project) IsOnline() bool {
	return p.RepoType.RepoType() == RepoTypeOnline
}

func (p *Project) ResourceIndex() ResourceIndex {
	return ResourceIndex{
		Owner: p.Owner,
		Id:    p.Id,
	}
}

func (p *Project) ResourceObject() (ResourceObject, RepoType) {
	return ResourceObject{
		Type:          ResourceTypeProject,
		ResourceIndex: p.ResourceIndex(),
	}, p.RepoType
}

func (p *Project) RelatedResources() []ResourceObjects {
	r := make([]ResourceObjects, 0, 2)

	if len(p.RelatedModels) > 0 {
		r = append(r, ResourceObjects{
			Type:    ResourceTypeModel,
			Objects: p.RelatedModels,
		})
	}

	if len(p.RelatedDatasets) > 0 {
		r = append(r, ResourceObjects{
			Type:    ResourceTypeDataset,
			Objects: p.RelatedDatasets,
		})
	}

	return r
}

type ProjectModifiableProperty struct {
	Name     ResourceName
	Desc     ResourceDesc
	Title    ResourceTitle
	CoverId  CoverId
	RepoType RepoType
	Tags     []string
	TagKinds []string
	Level    ResourceLevel
}

type ProjectSummary struct {
	Id            string
	Owner         Account
	Name          ResourceName
	Desc          ResourceDesc
	Title         ResourceTitle
	Level         ResourceLevel
	CoverId       CoverId
	Tags          []string
	UpdatedAt     int64
	LikeCount     int
	ForkCount     int
	DownloadCount int
}
