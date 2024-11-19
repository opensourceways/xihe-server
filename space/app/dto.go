package app

import (
	sdk "github.com/opensourceways/xihe-sdk/space"

	"github.com/opensourceways/xihe-server/app"
	spacedomain "github.com/opensourceways/xihe-server/space/domain"
	"github.com/opensourceways/xihe-server/utils"
)

// Project
type GlobalProjectsDTO struct {
	Total    int                `json:"total"`
	Projects []GlobalProjectDTO `json:"projects"`
}

type GlobalProjectDTO struct {
	ProjectSummaryDTO
	AvatarId string `json:"avatar_id"`
}

type ProjectsDTO struct {
	Total    int                 `json:"total"`
	Projects []ProjectSummaryDTO `json:"projects"`
}

type ProjectSummaryDTO struct {
	Id            string   `json:"id"`
	Owner         string   `json:"owner"`
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	Title         string   `json:"title"`
	Level         string   `json:"level"`
	CoverId       string   `json:"cover_id"`
	Tags          []string `json:"tags"`
	UpdatedAt     string   `json:"updated_at"`
	LikeCount     int      `json:"like_count"`
	ForkCount     int      `json:"fork_count"`
	DownloadCount int      `json:"download_count"`
	IsNpu         bool     `json:"is_npu"`
}

type ProjectDTO struct {
	Id            string   `json:"id"`
	Owner         string   `json:"owner"`
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	Title         string   `json:"title"`
	Type          string   `json:"type"`
	CoverId       string   `json:"cover_id"`
	Protocol      string   `json:"protocol"`
	Training      string   `json:"training"`
	RepoType      string   `json:"repo_type"`
	RepoId        string   `json:"repo_id"`
	Tags          []string `json:"tags"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
	LikeCount     int      `json:"like_count"`
	ForkCount     int      `json:"fork_count"`
	DownloadCount int      `json:"download_count"`
	CommitId      string   `json:"commit_id"`
	BaseImage     string   `json:"base_image"`
	Hardware      string   `json:"hardware"`
	IsNpu         bool     `json:"is_npu"`
}

type ProjectDetailDTO struct {
	ProjectDTO

	RelatedModels   []app.ResourceDTO `json:"related_models"`
	RelatedDatasets []app.ResourceDTO `json:"related_datasets"`
}

// CmdToNotifyUpdateCode is to update no application file and commitId
type CmdToNotifyUpdateCode struct {
	CommitId          string
	NoApplicationFile bool
}

func (s projectService) toSpaceMetaDTO(v spacedomain.Project) sdk.SpaceMetaDTO {
	return sdk.SpaceMetaDTO{
		Id:           v.RepoId,
		SDK:          v.Type.ProjType(),
		Name:         v.Name.ResourceName(),
		Owner:        v.Owner.Account(),
		Hardware:     v.Hardware.Hardware(),
		BaseImage:    v.BaseImage.BaseImage(),
		Visibility:   v.RepoType.RepoType(),
		Disable:      false,
		HardwareType: v.Hardware.Hardware(),
		CommitId:     v.CommitId,
	}
}

func (s projectService) toProjectDTO(p *spacedomain.Project, dto *ProjectDTO) {
	*dto = ProjectDTO{
		Id:            p.Id,
		Owner:         p.Owner.Account(),
		Name:          p.Name.ResourceName(),
		Type:          p.Type.ProjType(),
		CoverId:       p.CoverId.CoverId(),
		Protocol:      p.Protocol.ProtocolName(),
		Training:      p.Training.TrainingPlatform(),
		RepoType:      p.RepoType.RepoType(),
		RepoId:        p.RepoId,
		Tags:          p.Tags,
		CreatedAt:     utils.ToDate(p.CreatedAt),
		UpdatedAt:     utils.ToDate(p.UpdatedAt),
		LikeCount:     p.LikeCount,
		ForkCount:     p.ForkCount,
		DownloadCount: p.DownloadCount,
		BaseImage:     p.BaseImage.BaseImage(),
		Hardware:      p.Hardware.Hardware(),
		IsNpu:         p.Hardware.IsNpu(),
	}

	if p.CommitId != "" {
		dto.CommitId = p.CommitId
	}

	if p.Desc != nil {
		dto.Desc = p.Desc.ResourceDesc()
	}

	if p.Title != nil {
		dto.Title = p.Title.ResourceTitle()
	}

}

func (s projectService) toProjectSummaryDTO(p *spacedomain.ProjectSummary, dto *ProjectSummaryDTO) {
	*dto = ProjectSummaryDTO{
		Id:            p.Id,
		Owner:         p.Owner.Account(),
		Name:          p.Name.ResourceName(),
		CoverId:       p.CoverId.CoverId(),
		Tags:          p.Tags,
		UpdatedAt:     utils.ToDate(p.UpdatedAt),
		LikeCount:     p.LikeCount,
		ForkCount:     p.ForkCount,
		DownloadCount: p.DownloadCount,
		IsNpu:         p.Hardware.IsNpu(),
	}

	if p.Desc != nil {
		dto.Desc = p.Desc.ResourceDesc()
	}

	if p.Title != nil {
		dto.Title = p.Title.ResourceTitle()
	}

	if p.Level != nil {
		dto.Level = p.Level.ResourceLevel()
	}
}
