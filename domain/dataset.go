package domain

type Dataset struct {
	Id string

	Owner    Account
	Protocol ProtocolName

	DatasetModifiableProperty

	RepoId string

	CreatedAt int64
	UpdatedAt int64

	Version int

	// following fileds is not under the controlling of version
	LikeCount     int
	DownloadCount int

	RelatedModels   RelatedResources
	RelatedProjects RelatedResources
}

func (d *Dataset) IsPrivate() bool {
	return d.RepoType.RepoType() == RepoTypePrivate
}

func (d *Dataset) ResourceIndex() ResourceIndex {
	return ResourceIndex{
		Owner: d.Owner,
		Id:    d.Id,
	}
}

func (d *Dataset) ResourceObject() ResourceObject {
	return ResourceObject{
		Type:          ResourceTypeDataset,
		ResourceIndex: d.ResourceIndex(),
	}
}

func (d *Dataset) RelatedResources() []ResourceObjects {
	r := make([]ResourceObjects, 0, 2)

	if len(d.RelatedProjects) > 0 {
		r = append(r, ResourceObjects{
			Type:    ResourceTypeProject,
			Objects: d.RelatedProjects,
		})
	}

	if len(d.RelatedModels) > 0 {
		r = append(r, ResourceObjects{
			Type:    ResourceTypeModel,
			Objects: d.RelatedModels,
		})
	}

	return r
}

type DatasetModifiableProperty struct {
	Name     ResourceName
	Desc     ResourceDesc
	Title    ResourceTitle
	RepoType RepoType
	Tags     []string
	TagKinds []string
}

type DatasetSummary struct {
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
