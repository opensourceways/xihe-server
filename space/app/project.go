package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	spacedomain "github.com/opensourceways/xihe-server/space/domain"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	userdomain "github.com/opensourceways/xihe-server/user/domain"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type ProjectCreateCmd struct {
	Owner    domain.Account
	Name     domain.ResourceName
	Desc     domain.ResourceDesc
	Title    domain.ResourceTitle
	Type     domain.ProjType
	CoverId  domain.CoverId
	RepoType domain.RepoType
	Protocol domain.ProtocolName
	Training domain.TrainingPlatform
	Tags     []string
	TagKinds []string
	All      []domain.DomainTags
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

func (cmd *ProjectCreateCmd) genTagKinds(tags []string) []string {
	if len(tags) == 0 {
		return nil
	}

	r := make([]string, 0, len(cmd.All))

	for i := range cmd.All {
		if v := cmd.All[i].GetKindsOfTags(tags); len(v) > 0 {
			r = append(r, v...)
		}
	}

	return r
}

func (cmd *ProjectCreateCmd) toProject(r *spacedomain.Project) {
	now := utils.Now()
	normTags := []string{cmd.Type.ProjType(),
		cmd.Protocol.ProtocolName(),
		cmd.Training.TrainingPlatform()}

	*r = spacedomain.Project{
		Owner:     cmd.Owner,
		Type:      cmd.Type,
		Protocol:  cmd.Protocol,
		Training:  cmd.Training,
		CreatedAt: now,
		UpdatedAt: now,
		ProjectModifiableProperty: spacedomain.ProjectModifiableProperty{
			Name:     cmd.Name,
			Desc:     cmd.Desc,
			Title:    cmd.Title,
			CoverId:  cmd.CoverId,
			RepoType: cmd.RepoType,
			Tags:     append(cmd.Tags, normTags...),
			TagKinds: cmd.genTagKinds(cmd.Tags),
		},
	}
}

type ProjectService interface {
	CanApplyResourceName(domain.Account, domain.ResourceName) bool
	Create(*ProjectCreateCmd, platform.Repository) (ProjectDTO, error)
	Delete(*spacedomain.Project, platform.Repository) error
	GetByName(domain.Account, domain.ResourceName, bool) (ProjectDetailDTO, error)
	List(domain.Account, *app.ResourceListCmd) (ProjectsDTO, error)
	ListGlobal(*app.GlobalResourceListCmd) (GlobalProjectsDTO, error)
	Update(*spacedomain.Project, *ProjectUpdateCmd, platform.Repository) (ProjectDTO, error)
	Fork(*ProjectForkCmd, platform.Repository) (ProjectDTO, error)

	AddRelatedModel(*spacedomain.Project, *domain.ResourceIndex) error
	RemoveRelatedModel(*spacedomain.Project, *domain.ResourceIndex) error

	AddRelatedDataset(*spacedomain.Project, *domain.ResourceIndex) error
	RemoveRelatedDataset(*spacedomain.Project, *domain.ResourceIndex) error

	SetTags(*spacedomain.Project, *app.ResourceTagsUpdateCmd) error
}

func NewProjectService(
	user userrepo.User,
	repo spacerepo.Project,
	model repository.Model,
	dataset repository.Dataset,
	activity repository.Activity,
	pr platform.Repository,
	sender message.ResourceProducer,
) ProjectService {
	return projectService{
		repo:     repo,
		activity: activity,
		sender:   sender,
		rs: app.ResourceService{
			User:    user,
			Model:   model,
			Project: repo,
			Dataset: dataset,
		},
	}
}

type projectService struct {
	repo spacerepo.Project
	//pr       platform.Repository
	activity repository.Activity
	sender   message.ResourceProducer
	rs       app.ResourceService
}

func (s projectService) CanApplyResourceName(owner domain.Account, name domain.ResourceName) bool {
	return s.rs.CanApplyResourceName(owner, name)
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
	v := new(spacedomain.Project)
	cmd.toProject(v)
	v.RepoId = pid

	p, err := s.repo.Save(v)
	if err != nil {
		return
	}

	s.toProjectDTO(&p, &dto)

	// add activity
	r, repoType := p.ResourceObject()
	ua := app.GenActivityForCreatingResource(r, repoType)
	_ = s.activity.Save(&ua)

	_ = s.sender.AddOperateLogForCreateResource(r, p.Name)

	_ = s.sender.CreateProject(message.ProjectCreatedEvent{
		Account:     r.Owner,
		ProjectName: dto.Name,
	})

	return
}

func (s projectService) Delete(r *spacedomain.Project, pr platform.Repository) (err error) {
	// step1: delete repo on gitlab
	if err = pr.Delete(r.RepoId); err != nil {
		return
	}

	obj, repoType := r.ResourceObject()

	// step2:
	if resources := r.RelatedResources(); len(resources) > 0 {
		err = s.sender.RemoveRelatedResources(&message.RelatedResources{
			Promoter:  obj,
			Resources: resources,
		})
		if err != nil {
			return
		}
	}

	// step3: delete
	if err = s.repo.Delete(&obj.ResourceIndex); err != nil {
		return
	}

	// add activity
	ua := app.GenActivityForDeletingResource(&obj, repoType)

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
		err = PrivateRepoError{
			errors.New("private repo"),
		}

		return
	}

	m, err := s.rs.ListModels(v.RelatedModels)
	if err != nil {
		return
	}
	dto.RelatedModels = m

	d, err := s.rs.ListDatasets(v.RelatedDatasets)
	if err != nil {
		return
	}
	dto.RelatedDatasets = d

	s.toProjectDTO(&v, &dto.ProjectDTO)

	return
}

func (s projectService) ListGlobal(cmd *app.GlobalResourceListCmd) (
	dto GlobalProjectsDTO, err error,
) {
	option := cmd.ToResourceListOption()

	var v spacerepo.UserProjectsInfo

	if cmd.SortType == nil {
		v, err = s.repo.ListGlobalAndSortByUpdateTime(&option)
	} else {
		switch cmd.SortType.SortType() {
		case domain.SortTypeUpdateTime:
			v, err = s.repo.ListGlobalAndSortByUpdateTime(&option)

		case domain.SortTypeFirstLetter:
			v, err = s.repo.ListGlobalAndSortByFirstLetter(&option)

		case domain.SortTypeDownloadCount:
			v, err = s.repo.ListGlobalAndSortByDownloadCount(&option)
		}
	}

	items := v.Projects

	if err != nil || len(items) == 0 {
		return
	}

	// find avatars
	users := make([]userdomain.Account, len(items))
	for i := range items {
		users[i] = items[i].Owner
	}

	avatars, err := s.rs.FindUserAvater(users)
	if err != nil {
		return
	}

	// gen result
	dtos := make([]GlobalProjectDTO, len(items))
	for i := range items {
		s.toProjectSummaryDTO(&items[i], &dtos[i].ProjectSummaryDTO)
		dtos[i].AvatarId = avatars[i]
	}

	dto.Total = v.Total
	dto.Projects = dtos

	return
}

func (s projectService) List(owner domain.Account, cmd *app.ResourceListCmd) (
	dto ProjectsDTO, err error,
) {
	option := cmd.ToResourceListOption()

	var v spacerepo.UserProjectsInfo

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