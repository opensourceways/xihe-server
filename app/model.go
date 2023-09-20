package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type ModelCreateCmd struct {
	Owner    domain.Account
	Name     domain.ResourceName
	Desc     domain.ResourceDesc
	Title    domain.ResourceTitle
	RepoType domain.RepoType
	Protocol domain.ProtocolName
	Tags     []string
	TagKinds []string
	All      []domain.DomainTags
}

func (cmd *ModelCreateCmd) Validate() error {
	b := cmd.Owner != nil &&
		cmd.Name != nil &&
		cmd.RepoType != nil &&
		cmd.Protocol != nil

	if !b {
		return errors.New("invalid cmd of creating model")
	}

	return nil
}

func (cmd *ModelCreateCmd) genTagKinds(tags []string) []string {
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

func (cmd *ModelCreateCmd) toModel(v *domain.Model) {
	now := utils.Now()

	normTags := []string{cmd.Protocol.ProtocolName()}
	*v = domain.Model{
		Owner:     cmd.Owner,
		Protocol:  cmd.Protocol,
		CreatedAt: now,
		UpdatedAt: now,
		ModelModifiableProperty: domain.ModelModifiableProperty{
			Name:     cmd.Name,
			Desc:     cmd.Desc,
			Title:    cmd.Title,
			RepoType: cmd.RepoType,
			Tags:     append(normTags, cmd.Tags...),
			TagKinds: cmd.genTagKinds(cmd.Tags),
		},
	}
}

type ModelsDTO struct {
	Total  int               `json:"total"`
	Models []ModelSummaryDTO `json:"models"`
}

type ModelSummaryDTO struct {
	Id            string   `json:"id"`
	Owner         string   `json:"owner"`
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	Title         string   `json:"title"`
	Tags          []string `json:"tags"`
	UpdatedAt     string   `json:"updated_at"`
	LikeCount     int      `json:"like_count"`
	DownloadCount int      `json:"download_count"`
}

type ModelDTO struct {
	Id            string   `json:"id"`
	Owner         string   `json:"owner"`
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	Title         string   `json:"title"`
	Protocol      string   `json:"protocol"`
	RepoType      string   `json:"repo_type"`
	RepoId        string   `json:"repo_id"`
	Tags          []string `json:"tags"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
	LikeCount     int      `json:"like_count"`
	DownloadCount int      `json:"download_count"`
}

type ModelDetailDTO struct {
	ModelDTO

	RelatedDatasets []ResourceDTO `json:"related_datasets"`
	RelatedProjects []ResourceDTO `json:"related_projects"`
}

type ModelService interface {
	CanApplyResourceName(domain.Account, domain.ResourceName) bool
	Create(*ModelCreateCmd, platform.Repository) (ModelDTO, error)
	Delete(*domain.Model, platform.Repository) error
	Update(*domain.Model, *ModelUpdateCmd, platform.Repository) (ModelDTO, error)
	GetByName(domain.Account, domain.ResourceName, bool) (ModelDetailDTO, error)
	List(domain.Account, *ResourceListCmd) (ModelsDTO, error)
	ListGlobal(*GlobalResourceListCmd) (GlobalModelsDTO, error)

	AddRelatedDataset(*domain.Model, *domain.ResourceIndex) error
	RemoveRelatedDataset(*domain.Model, *domain.ResourceIndex) error

	SetTags(*domain.Model, *ResourceTagsUpdateCmd) error
}

func NewModelService(
	user userrepo.User,
	repo repository.Model,
	proj repository.Project,
	dataset repository.Dataset,
	activity repository.Activity,
	pr platform.Repository,
	sender message.Sender,
) ModelService {
	return modelService{
		repo:     repo,
		activity: activity,
		sender:   sender,
		rs: resourceService{
			user:    user,
			model:   repo,
			project: proj,
			dataset: dataset,
		},
	}
}

type modelService struct {
	repo repository.Model
	//pr       platform.Repository
	activity repository.Activity
	rs       resourceService
	sender   message.Sender
}

func (s modelService) CanApplyResourceName(owner domain.Account, name domain.ResourceName) bool {
	return s.rs.canApplyResourceName(owner, name)
}

func (s modelService) Create(cmd *ModelCreateCmd, pr platform.Repository) (dto ModelDTO, err error) {
	pid, err := pr.New(&platform.RepoOption{
		Name:     cmd.Name,
		RepoType: cmd.RepoType,
	})
	if err != nil {
		return
	}

	v := new(domain.Model)
	cmd.toModel(v)
	v.RepoId = pid

	m, err := s.repo.Save(v)
	if err != nil {
		return
	}

	s.toModelDTO(&m, &dto)

	// add activity
	r, repoType := m.ResourceObject()
	ua := genActivityForCreatingResource(r, repoType)
	_ = s.activity.Save(&ua)

	_ = s.sender.AddOperateLogForCreateResource(r, m.Name)

	_ = s.sender.DailyCreate(cmd.Owner, m.Name)

	return
}

func (s modelService) Delete(r *domain.Model, pr platform.Repository) (err error) {
	// step1: delete repo on gitlab
	if err = pr.Delete(r.RepoId); err != nil {
		return
	}

	obj, repoType := r.ResourceObject()

	// step2: message
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
	ua := genActivityForDeletingResource(&obj, repoType)

	// ignore the error
	_ = s.activity.Save(&ua)

	return
}

func (s modelService) GetByName(
	owner domain.Account, name domain.ResourceName,
	allowPrivacy bool,
) (dto ModelDetailDTO, err error) {
	v, err := s.repo.GetByName(owner, name)
	if err != nil {
		return
	}

	if !allowPrivacy && v.IsPrivate() {
		err = ErrorPrivateRepo{errors.New("private repo")}

		return
	}

	d, err := s.rs.listDatasets(v.RelatedDatasets)
	if err != nil {
		return
	}
	dto.RelatedDatasets = d

	d, err = s.rs.listProjects(v.RelatedProjects)
	if err != nil {
		return
	}
	dto.RelatedProjects = d

	s.toModelDTO(&v, &dto.ModelDTO)

	return
}

func (s modelService) List(owner domain.Account, cmd *ResourceListCmd) (
	dto ModelsDTO, err error,
) {
	option := cmd.toResourceListOption()

	var v repository.UserModelsInfo

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

	items := v.Models

	if err != nil || len(items) == 0 {
		return
	}

	dtos := make([]ModelSummaryDTO, len(items))
	for i := range items {
		s.toModelSummaryDTO(&items[i], &dtos[i])
	}

	dto.Total = v.Total
	dto.Models = dtos

	return
}

func (s modelService) toModelDTO(m *domain.Model, dto *ModelDTO) {
	*dto = ModelDTO{
		Id:            m.Id,
		Owner:         m.Owner.Account(),
		Name:          m.Name.ResourceName(),
		Protocol:      m.Protocol.ProtocolName(),
		RepoType:      m.RepoType.RepoType(),
		RepoId:        m.RepoId,
		Tags:          m.Tags,
		CreatedAt:     utils.ToDate(m.CreatedAt),
		UpdatedAt:     utils.ToDate(m.UpdatedAt),
		LikeCount:     m.LikeCount,
		DownloadCount: m.DownloadCount,
	}

	if m.Desc != nil {
		dto.Desc = m.Desc.ResourceDesc()
	}

	if m.Title != nil {
		dto.Title = m.Title.ResourceTitle()
	}

}

func (s modelService) toModelSummaryDTO(m *domain.ModelSummary, dto *ModelSummaryDTO) {
	*dto = ModelSummaryDTO{
		Id:            m.Id,
		Owner:         m.Owner.Account(),
		Name:          m.Name.ResourceName(),
		Tags:          m.Tags,
		UpdatedAt:     utils.ToDate(m.UpdatedAt),
		LikeCount:     m.LikeCount,
		DownloadCount: m.DownloadCount,
	}

	if m.Desc != nil {
		dto.Desc = m.Desc.ResourceDesc()
	}

	if m.Title != nil {
		dto.Title = m.Title.ResourceTitle()
	}

}
