package repositoryimpl

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	spacedomain "github.com/opensourceways/xihe-server/space/domain"
)

var (
	projectTableName = ""
	tagsTableName    = ""
	datasetTableName = ""
	modelTableName   = ""
)

func (do *projectDO) TableName() string {
	return projectTableName
}

type projectDO struct {
	Id            string
	Owner         string
	Name          string
	FL            byte
	Description   string
	Title         string
	Type          string
	Level         int
	CoverId       string
	Protocol      string
	Training      string
	RepoType      string
	RepoId        string
	CreatedAt     int64
	UpdatedAt     int64
	Version       int
	LikeCount     int
	ForkCount     int
	DownloadCount int

	CommitId           string
	NoApplicationFile  bool
	CompPowerAllocated bool
	Exception          string

	Hardware  string
	BaseImage string
}

type projectTagsDO struct {
	Id        int
	ProjectId string
	TagName   string
}

type datasetDO struct {
	DatasetId string
	ProjectId string
	Owner     string
}

type modelDO struct {
	ModelId   string
	ProjectId string
	Owner     string
}

func toProjectDO(p *spacedomain.Project) projectDO {
	do := projectDO{
		Id:                 p.RepoId,
		Owner:              p.Owner.Account(),
		Name:               p.Name.ResourceName(),
		FL:                 p.Name.FirstLetterOfName(),
		Type:               p.Type.ProjType(),
		CoverId:            p.CoverId.CoverId(),
		RepoType:           p.RepoType.RepoType(),
		Protocol:           p.Protocol.ProtocolName(),
		Training:           p.Training.TrainingPlatform(),
		RepoId:             p.RepoId,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
		Version:            p.Version,
		Hardware:           p.Hardware.Hardware(),
		BaseImage:          p.BaseImage.BaseImage(),
		CommitId:           p.CommitId,
		NoApplicationFile:  p.NoApplicationFile,
		CompPowerAllocated: p.CompPowerAllocated,
	}

	if p.Desc != nil {
		do.Description = p.Desc.ResourceDesc()
	}

	if p.Title != nil {
		do.Title = p.Title.ResourceTitle()
	}
	return do
}

func toProjectTagsDO(p *spacedomain.Project) []projectTagsDO {
	var tags []projectTagsDO

	for _, v := range p.Tags {
		tags = append(tags, projectTagsDO{
			ProjectId: p.RepoId,
			TagName:   v,
		})
	}

	return tags
}

func (do *projectDO) toProject(r *spacedomain.Project) (err error) {
	r.Id = do.Id

	if r.Owner, err = domain.NewAccount(do.Owner); err != nil {
		return
	}

	if r.Name, err = domain.NewResourceName(do.Name); err != nil {
		return
	}

	if r.Desc, err = domain.NewResourceDesc(do.Description); err != nil {
		return
	}

	if r.Title, err = domain.NewResourceTitle(do.Title); err != nil {
		return
	}

	if r.Type, err = domain.NewProjType(do.Type); err != nil {
		return
	}

	if r.CoverId, err = domain.NewCoverId(do.CoverId); err != nil {
		return
	}

	if r.RepoType, err = domain.NewRepoType(do.RepoType); err != nil {
		return
	}

	if r.Protocol, err = domain.NewProtocolName(do.Protocol); err != nil {
		return
	}

	if r.Training, err = domain.NewTrainingPlatform(do.Training); err != nil {
		return
	}

	if do.Hardware != "" {
		if r.Hardware, err = domain.NewHardware(do.Hardware, do.Type); err != nil {
			return
		}

	}
	if do.BaseImage != "" {
		if r.BaseImage, err = domain.NewBaseImage(do.BaseImage, do.Hardware); err != nil {
			return
		}
	}

	r.Level = domain.NewResourceLevelByNum(do.Level)
	r.RepoId = do.RepoId
	r.Version = do.Version
	r.CreatedAt = do.CreatedAt
	r.UpdatedAt = do.UpdatedAt
	r.LikeCount = do.LikeCount
	r.ForkCount = do.ForkCount
	r.DownloadCount = do.DownloadCount
	r.CommitId = do.CommitId
	r.NoApplicationFile = do.NoApplicationFile
	r.Exception = domain.CreateException(do.Exception)
	r.CompPowerAllocated = do.CompPowerAllocated

	return
}

func toDatasetDO(r *repository.RelatedResourceInfo) datasetDO {
	do := datasetDO{
		ProjectId: r.ResourceToUpdate.Id,
		DatasetId: r.RelatedResource.Id,
		Owner:     r.RelatedResource.Owner.Account(),
	}
	return do
}

func toModelDO(r *repository.RelatedResourceInfo) modelDO {
	do := modelDO{
		ProjectId: r.ResourceToUpdate.Id,
		ModelId:   r.RelatedResource.Id,
		Owner:     r.ResourceToUpdate.Owner.Account(),
	}
	return do
}