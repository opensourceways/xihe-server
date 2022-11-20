package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type ProjectCreateCmd struct {
	Owner    domain.Account
	Name     domain.ResourceName
	Desc     domain.ResourceDesc
	Type     domain.ProjType
	CoverId  domain.CoverId
	RepoType domain.RepoType
	Protocol domain.ProtocolName
	Training domain.TrainingPlatform
}

func (cmd *ProjectCreateCmd) Validate() error {
	b := cmd.Owner != nil &&
		cmd.Name != nil &&
		cmd.Type != nil &&
		cmd.CoverId != nil &&
		cmd.RepoType != nil &&
		cmd.Protocol != nil &&
		cmd.Training != nil

	if !b {
		return errors.New("invalid cmd of creating project")
	}

	return nil
}

func (cmd *ProjectCreateCmd) toProject(r *domain.Project) {
	now := utils.Now()

	*r = domain.Project{
		Owner:     cmd.Owner,
		Type:      cmd.Type,
		Protocol:  cmd.Protocol,
		Training:  cmd.Training,
		CreatedAt: now,
		UpdatedAt: now,
		ProjectModifiableProperty: domain.ProjectModifiableProperty{
			Name:     cmd.Name,
			Desc:     cmd.Desc,
			CoverId:  cmd.CoverId,
			RepoType: cmd.RepoType,
			Tags: []string{
				cmd.Type.ProjType(),
				cmd.Protocol.ProtocolName(),
				cmd.Training.TrainingPlatform(),
			},
			TagKinds: []string{},
		},
	}
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
}

type ProjectDetailDTO struct {
	ProjectDTO

	RelatedModels   []ResourceDTO `json:"related_models"`
	RelatedDatasets []ResourceDTO `json:"related_datasets"`
}

type ProjectService interface {
	CanApplyResourceName(domain.Account, domain.ResourceName) bool
	Create(*ProjectCreateCmd, platform.Repository) (ProjectDTO, error)
	Delete(*domain.ResourceSummary, platform.Repository) error
	GetByName(domain.Account, domain.ResourceName, bool) (ProjectDetailDTO, error)
	List(domain.Account, *ResourceListCmd) (ProjectsDTO, error)
	ListGlobal(*GlobalResourceListCmd) (GlobalProjectsDTO, error)
	Update(*domain.Project, *ProjectUpdateCmd, platform.Repository) (ProjectDTO, error)
	Fork(*ProjectForkCmd, platform.Repository) (ProjectDTO, error)

	IncreaseFork(*domain.ResourceIndex) error

	AddLike(*domain.ResourceIndex) error
	RemoveLike(*domain.ResourceIndex) error

	AddRelatedModel(*domain.Project, *domain.ResourceIndex) error
	RemoveRelatedModel(*domain.Project, *domain.ResourceIndex) error

	AddRelatedDataset(*domain.Project, *domain.ResourceIndex) error
	RemoveRelatedDataset(*domain.Project, *domain.ResourceIndex) error

	SetTags(*domain.Project, *ResourceTagsUpdateCmd) error
}

func NewProjectService(
	user repository.User,
	repo repository.Project,
	model repository.Model,
	dataset repository.Dataset,
	activity repository.Activity,
	pr platform.Repository,
	sender message.Sender,
) ProjectService {
	return projectService{
		repo:     repo,
		activity: activity,
		sender:   sender,
		rs: resourceService{
			user:    user,
			model:   model,
			project: repo,
			dataset: dataset,
		},
	}
}

type projectService struct {
	repo repository.Project
	//pr       platform.Repository
	activity repository.Activity
	sender   message.Sender
	rs       resourceService
}

func (s projectService) CanApplyResourceName(owner domain.Account, name domain.ResourceName) bool {
	return s.rs.canApplyResourceName(owner, name)
}

func (s projectService) Create(cmd *ProjectCreateCmd, pr platform.Repository) (dto ProjectDTO, err error) {
	// step1: create repo on gitlab
	pid, err := pr.New(&platform.RepoOption{
		Name:     cmd.Name,
		RepoType: cmd.RepoType,
	})
	if err != nil {
		return
	}

	// step2: save
	v := new(domain.Project)
	cmd.toProject(v)
	v.RepoId = pid

	p, err := s.repo.Save(v)
	if err != nil {
		return
	}

	s.toProjectDTO(&p, &dto)

	// add activity
	ua := genActivityForCreatingResource(
		p.Owner, domain.ResourceTypeProject, p.Id,
	)
	// ignore the error
	_ = s.activity.Save(&ua)

	return
}

func (s projectService) Delete(r *domain.ResourceSummary, pr platform.Repository) (err error) {
	// step1: delete repo on gitlab
	if err = pr.Delete(r.RepoId); err != nil {
		return
	}

	// step2: save
	index := r.ResourceIndex()
	if err = s.repo.Delete(&index); err != nil {
		return
	}

	// add activity
	ua := genActivityForDeletingResource(r, domain.ResourceTypeProject)

	// ignore the error
	_ = s.activity.Save(&ua)

	return
}

func (s projectService) GetByName(
	owner domain.Account, name domain.ResourceName,
	allowPrivacy bool,
) (dto ProjectDetailDTO, err error) {
	v, err := s.repo.GetByName(owner, name)
	if err != nil {
		return
	}

	if !allowPrivacy && v.IsPrivate() {
		err = ErrorPrivateRepo{errors.New("private repo")}

		return
	}

	m, err := s.rs.listModels(v.RelatedModels)
	if err != nil {
		return
	}
	dto.RelatedModels = m

	d, err := s.rs.listDatasets(v.RelatedDatasets)
	if err != nil {
		return
	}
	dto.RelatedDatasets = d

	s.toProjectDTO(&v, &dto.ProjectDTO)

	return
}

type ResourceListCmd struct {
	repository.ResourceListOption

	SortType domain.SortType
}

func (cmd *ResourceListCmd) toResourceListOption() repository.ResourceListOption {
	return cmd.ResourceListOption
}

func (s projectService) List(owner domain.Account, cmd *ResourceListCmd) (
	dto ProjectsDTO, err error,
) {
	option := cmd.toResourceListOption()

	var v repository.UserProjectsInfo

	if cmd.SortType == nil {
		v, err = s.repo.ListAndSortByUpdateTime(owner, &option)
	} else {
		switch cmd.SortType.SortType() {
		case domain.SortTypeUpdateTime:
			v, err = s.repo.ListAndSortByUpdateTime(owner, &option)

		case domain.SortTypeFirstLetter:
			v, err = s.repo.ListAndSortByFirstLetter(owner, &option)

		case domain.SortTypeDownloadCount:
			v, err = s.repo.ListAndSortByDownloadCount(owner, &option)
		}
	}

	items := v.Projects

	if err != nil || len(items) == 0 {
		return
	}

	dtos := make([]ProjectSummaryDTO, len(items))
	for i := range items {
		s.toProjectSummaryDTO(&items[i], &dtos[i])
	}

	dto.Total = v.Total
	dto.Projects = dtos

	return
}

func (s projectService) toProjectDTO(p *domain.Project, dto *ProjectDTO) {
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
	}

	if p.Desc != nil {
		dto.Desc = p.Desc.ResourceDesc()
	}

}

func (s projectService) toProjectSummaryDTO(p *domain.ProjectSummary, dto *ProjectSummaryDTO) {
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
	}

	if p.Desc != nil {
		dto.Desc = p.Desc.ResourceDesc()
	}

	if p.Level != nil {
		dto.Level = p.Level.ResourceLevel()
	}
}
