package domain

type Model struct {
	Id string

	Owner    Account
	Protocol ProtocolName

	ModelModifiableProperty

	RepoId string

	RelatedDatasets RelatedResources

	CreatedAt int64
	UpdatedAt int64

	Version int

	// following fileds is not under the controlling of version
	LikeCount       int
	DownloadCount   int
	RelatedProjects RelatedResources
}

func (m *Model) MaxRelatedResourceNum() int {
	return DomainConfig.MaxRelatedResourceNum
}

func (m *Model) IsPrivate() bool {
	return m.RepoType.RepoType() == RepoTypePrivate
}

func (m *Model) ResourceIndex() ResourceIndex {
	return ResourceIndex{
		Owner: m.Owner,
		Id:    m.Id,
	}
}

func (m *Model) ResourceObject() (ResourceObject, RepoType) {
	return ResourceObject{
		Type:          ResourceTypeModel,
		ResourceIndex: m.ResourceIndex(),
	}, m.RepoType
}

func (m *Model) RelatedResources() []ResourceObjects {
	r := make([]ResourceObjects, 0, 2)

	if len(m.RelatedProjects) > 0 {
		r = append(r, ResourceObjects{
			Type:    ResourceTypeProject,
			Objects: m.RelatedProjects,
		})
	}

	if len(m.RelatedDatasets) > 0 {
		r = append(r, ResourceObjects{
			Type:    ResourceTypeDataset,
			Objects: m.RelatedDatasets,
		})
	}

	return r
}

type ModelModifiableProperty struct {
	Name     ResourceName
	Desc     ResourceDesc
	Title    ResourceTitle
	RepoType RepoType
	Tags     []string
	TagKinds []string
}

type ModelSummary struct {
	Id            string
	Owner         Account
	Name          ResourceName
	Desc          ResourceDesc
	Title         ResourceTitle
	Tags          []string
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int
}
