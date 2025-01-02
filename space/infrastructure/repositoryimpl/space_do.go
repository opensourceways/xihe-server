package repositoryimpl

import (
	"strconv"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
	spacedomain "github.com/opensourceways/xihe-server/space/domain"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
)

const (
	fieldOwner         = "owner"
	fieldName          = "name"
	fieldID            = "id"
	fieldLikeCount     = "like_count"
	fieldProjectId     = "project_id"
	fieldForkCount     = "fork_count"
	fieldDownload      = "download_count"
	fieldRepoType      = "repo_type"
	fieldLevel         = "level"
	fieldKind          = "kind"
	fieldTagName       = "tag_name"
	tableTagCategories = "tag_categories"
	tableProjectTags   = "project_tags"
	tableProjects      = "projects"
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
	Id            int64  `gorm:"column:id;primaryKey"`
	Owner         string `gorm:"column:owner"`
	Name          string `gorm:"column:name"`
	FL            byte   `gorm:"column:fl"`
	Description   string `gorm:"column:description"`
	Title         string `gorm:"column:title"`
	Type          string `gorm:"column:type"`
	Level         int    `gorm:"column:level"`
	CoverId       string `gorm:"column:cover_id"`
	Protocol      string `gorm:"column:protocol"`
	Training      string `gorm:"column:training"`
	RepoType      string `gorm:"column:repo_type"`
	RepoId        int64  `gorm:"column:repo_id"`
	CreatedAt     int64  `gorm:"column:created_at"`
	UpdatedAt     int64  `gorm:"column:updated_at"`
	Version       int    `gorm:"column:version"`
	LikeCount     int    `gorm:"column:like_count"`
	ForkCount     int    `gorm:"column:fork_count"`
	DownloadCount int    `gorm:"column:download_count"`

	CommitId           string `gorm:"column:commit_id"`
	NoApplicationFile  bool   `gorm:"column:no_application_file"`
	CompPowerAllocated bool   `gorm:"column:comp_power_allocated"`
	Exception          string `gorm:"column:exception"`

	Hardware  string `gorm:"column:hardware"`
	BaseImage string `gorm:"column:base_image"`
}

type projectTagsDO struct {
	Id        int    `gorm:"column:id;primaryKey"`
	ProjectId int64  `gorm:"column:project_id"`
	TagName   string `gorm:"column:tag_name"`
}

type datasetDO struct {
	DatasetId string `gorm:"column:dataset_id;primaryKey"`
	ProjectId int64  `gorm:"column:project_id"`
	Owner     string `gorm:"column:owner"`
}

type modelDO struct {
	ModelId   string `gorm:"column:model_id;primaryKey"`
	ProjectId int64  `gorm:"column:project_id"`
	Owner     string `gorm:"column:owner"`
}

func toProjectDO(p *spacedomain.Project) projectDO {
	idInt64, err := strconv.ParseInt(p.RepoId, 10, 64)
	if err != nil {
		return projectDO{}
	}
	do := projectDO{
		Id:                 idInt64,
		Owner:              p.Owner.Account(),
		Name:               p.Name.ResourceName(),
		FL:                 p.Name.FirstLetterOfName(),
		Type:               p.Type.ProjType(),
		CoverId:            p.CoverId.CoverId(),
		RepoType:           p.RepoType.RepoType(),
		Protocol:           p.Protocol.ProtocolName(),
		Training:           p.Training.TrainingPlatform(),
		RepoId:             idInt64,
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
		idInt64, err := strconv.ParseInt(p.RepoId, 10, 64)
		if err != nil {
			return nil
		}
		tags = append(tags, projectTagsDO{
			ProjectId: idInt64,
			TagName:   v,
		})
	}

	return tags
}

func (do *projectDO) toProject(r *spacedomain.Project) (err error) {
	r.Id = strconv.Itoa(int(do.Id))

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

	if r.Hardware, err = domain.NewHardware(do.Hardware, do.Type); err != nil {
		return
	}

	if r.BaseImage, err = domain.NewBaseImage(do.BaseImage, do.Hardware); err != nil {
		return
	}

	r.Level = domain.NewResourceLevelByNum(do.Level)
	r.RepoId = strconv.Itoa(int(do.RepoId))
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
	projectIdInt64, err := strconv.ParseInt(r.ResourceToUpdate.Id, 10, 64)
	if err != nil {
		return datasetDO{}
	}

	do := datasetDO{
		ProjectId: projectIdInt64,
		DatasetId: r.RelatedResource.Id,
		Owner:     r.RelatedResource.Owner.Account(),
	}
	return do
}

func toModelDO(r *repository.RelatedResourceInfo) modelDO {
	projectIdInt64, err := strconv.ParseInt(r.ResourceToUpdate.Id, 10, 64)
	if err != nil {
		return modelDO{}
	}

	do := modelDO{
		ProjectId: projectIdInt64,
		ModelId:   r.RelatedResource.Id,
		Owner:     r.RelatedResource.Owner.Account(),
	}
	return do
}

func toProjectSummaryDO(item projectDO, tags []string) ProjectSummaryDO {
	return ProjectSummaryDO{
		Id:            strconv.FormatInt(int64(item.Id), 10),
		Owner:         item.Owner,
		Name:          item.Name,
		Desc:          item.Description,
		Title:         item.Title,
		Level:         item.Level,
		CoverId:       item.CoverId,
		UpdatedAt:     item.UpdatedAt,
		LikeCount:     item.LikeCount,
		ForkCount:     item.ForkCount,
		DownloadCount: item.DownloadCount,
		Hardware:      item.Hardware,
		Type:          item.Type,
		Tags:          tags,
	}
}

func toProjectResourceSummaryDO(project projectDO, tags []string) ProjectResourceSummaryDO {
	return ProjectResourceSummaryDO{
		ResourceSummaryDO: repositories.ResourceSummaryDO{
			Owner:    project.Owner,
			Name:     project.Name,
			Id:       strconv.Itoa(int(project.Id)),
			RepoId:   strconv.Itoa(int(project.RepoId)),
			RepoType: project.RepoType,
		},
		Tags: tags,
	}
}
func toProjectDOFromUpdateInfo(info spacerepo.ProjectPropertyUpdateInfo) projectDO {
	p := &info.Property

	idInt64, err := strconv.ParseInt(info.Id, 10, 64)
	if err != nil {
		return projectDO{}
	}

	return projectDO{
		Id:          idInt64,
		Owner:       info.Owner.Account(),
		Version:     info.Version,
		UpdatedAt:   info.UpdatedAt,
		Name:        p.Name.ResourceName(),
		FL:          p.Name.FirstLetterOfName(),
		Description: p.Desc.ResourceDesc(),
		Title:       p.Title.ResourceTitle(),
		Level: func() int {
			if p.Level != nil {
				return p.Level.Int()
			}
			return 0
		}(),
		CoverId:           p.CoverId.CoverId(),
		RepoType:          p.RepoType.RepoType(),
		CommitId:          p.CommitId,
		NoApplicationFile: p.NoApplicationFile,
		Exception:         p.Exception.Exception(),
	}
}
