package app

import "github.com/opensourceways/xihe-server/app"

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
	IsNPU         bool     `json:"is_npu"`
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
