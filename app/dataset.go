package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type DatasetCreateCmd struct {
	Owner    domain.Account
	Name     domain.ResourceName
	Desc     domain.ResourceDesc
	RepoType domain.RepoType
	Protocol domain.ProtocolName
}

func (cmd *DatasetCreateCmd) Validate() error {
	b := cmd.Owner != nil &&
		cmd.Name != nil &&
		cmd.RepoType != nil &&
		cmd.Protocol != nil

	if !b {
		return errors.New("invalid cmd of creating dataset")
	}

	return nil
}

func (cmd *DatasetCreateCmd) toDataset(v *domain.Dataset) {
	now := utils.Now()

	*v = domain.Dataset{
		Owner:     cmd.Owner,
		Protocol:  cmd.Protocol,
		CreatedAt: now,
		UpdatedAt: now,
		DatasetModifiableProperty: domain.DatasetModifiableProperty{
			Name:     cmd.Name,
			Desc:     cmd.Desc,
			RepoType: cmd.RepoType,
			Tags:     []string{cmd.Protocol.ProtocolName()},
			TagKinds: []string{},
		},
	}
}

type DatasetsDTO struct {
	Total    int                 `json:"total"`
	Datasets []DatasetSummaryDTO `json:"models"`
}

type DatasetSummaryDTO struct {
	Id            string   `json:"id"`
	Owner         string   `json:"owner"`
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	Tags          []string `json:"tags"`
	UpdatedAt     string   `json:"updated_at"`
	LikeCount     int      `json:"like_count"`
	DownloadCount int      `json:"download_count"`
}

type DatasetDTO struct {
	Id            string   `json:"id"`
	Owner         string   `json:"owner"`
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	Protocol      string   `json:"protocol"`
	RepoType      string   `json:"repo_type"`
	RepoId        string   `json:"repo_id"`
	Tags          []string `json:"tags"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
	LikeCount     int      `json:"like_count"`
	DownloadCount int      `json:"download_count"`
}

type DatasetDetailDTO struct {
	DatasetDTO

	RelatedModels   []ResourceDTO `json:"related_models"`
	RelatedProjects []ResourceDTO `json:"related_projects"`
}

type DatasetService interface {
	CanApplyResourceName(domain.Account, domain.ResourceName) bool
	Create(*DatasetCreateCmd, platform.Repository) (DatasetDTO, error)
	Delete(*domain.Dataset, platform.Repository) error
	Update(*domain.Dataset, *DatasetUpdateCmd, platform.Repository) (DatasetDTO, error)
	GetByName(domain.Account, domain.ResourceName, bool) (DatasetDetailDTO, error)
	List(domain.Account, *ResourceListCmd) (DatasetsDTO, error)
	ListGlobal(*GlobalResourceListCmd) (GlobalDatasetsDTO, error)

	AddLike(*domain.ResourceIndex) error
	RemoveLike(*domain.ResourceIndex) error

	AddRelatedProject(*domain.ReverselyRelatedResourceInfo) error
	RemoveRelatedProject(*domain.ReverselyRelatedResourceInfo) error

	AddRelatedModel(*domain.ReverselyRelatedResourceInfo) error
	RemoveRelatedModel(*domain.ReverselyRelatedResourceInfo) error

	SetTags(*domain.Dataset, *ResourceTagsUpdateCmd) error
}

func NewDatasetService(
	user repository.User,
	repo repository.Dataset,
	proj repository.Project,
	model repository.Model,
	activity repository.Activity,
	pr platform.Repository,
	sender message.Sender,
) DatasetService {
	return datasetService{
		repo:     repo,
		activity: activity,
		sender:   sender,
		rs: resourceService{
			user:    user,
			model:   model,
			project: proj,
			dataset: repo,
		},
	}
}

type datasetService struct {
	repo repository.Dataset
	//pr       platform.Repository
	activity repository.Activity
	sender   message.Sender
	rs       resourceService
}

func (s datasetService) CanApplyResourceName(owner domain.Account, name domain.ResourceName) bool {
	return s.rs.canApplyResourceName(owner, name)
}

func (s datasetService) Create(cmd *DatasetCreateCmd, pr platform.Repository) (dto DatasetDTO, err error) {
	pid, err := pr.New(&platform.RepoOption{
		Name:     cmd.Name,
		RepoType: cmd.RepoType,
	})
	if err != nil {
		return
	}

	v := new(domain.Dataset)
	cmd.toDataset(v)
	v.RepoId = pid

	d, err := s.repo.Save(v)
	if err != nil {
		return
	}

	s.toDatasetDTO(&d, &dto)

	// add activity
	ua := genActivityForCreatingResource(
		d.Owner, domain.ResourceTypeDataset, d.Id,
	)
	// ignore the error
	_ = s.activity.Save(&ua)

	return
}

func (s datasetService) Delete(r *domain.Dataset, pr platform.Repository) (err error) {
	// step1: delete repo on gitlab
	if err = pr.Delete(r.RepoId); err != nil {
		return
	}

	obj := r.ResourceObject()

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
	ua := genActivityForDeletingResource(&obj)

	// ignore the error
	_ = s.activity.Save(&ua)

	return
}

func (s datasetService) GetByName(
	owner domain.Account, name domain.ResourceName,
	allowPrivacy bool,
) (dto DatasetDetailDTO, err error) {
	v, err := s.repo.GetByName(owner, name)
	if err != nil {
		return
	}

	if !allowPrivacy && v.IsPrivate() {
		err = ErrorPrivateRepo{errors.New("private repo")}

		return
	}

	d, err := s.rs.listModels(v.RelatedModels)
	if err != nil {
		return
	}
	dto.RelatedModels = d

	d, err = s.rs.listProjects(v.RelatedProjects)
	if err != nil {
		return
	}
	dto.RelatedProjects = d

	s.toDatasetDTO(&v, &dto.DatasetDTO)

	return
}

func (s datasetService) List(owner domain.Account, cmd *ResourceListCmd) (
	dto DatasetsDTO, err error,
) {
	option := cmd.toResourceListOption()

	var v repository.UserDatasetsInfo

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

	items := v.Datasets

	if err != nil || len(items) == 0 {
		return
	}

	dtos := make([]DatasetSummaryDTO, len(items))
	for i := range items {
		s.toDatasetSummaryDTO(&items[i], &dtos[i])
	}

	dto.Total = v.Total
	dto.Datasets = dtos

	return
}

func (s datasetService) toDatasetDTO(d *domain.Dataset, dto *DatasetDTO) {
	*dto = DatasetDTO{
		Id:            d.Id,
		Owner:         d.Owner.Account(),
		Name:          d.Name.ResourceName(),
		Protocol:      d.Protocol.ProtocolName(),
		RepoType:      d.RepoType.RepoType(),
		RepoId:        d.RepoId,
		Tags:          d.Tags,
		CreatedAt:     utils.ToDate(d.CreatedAt),
		UpdatedAt:     utils.ToDate(d.UpdatedAt),
		LikeCount:     d.LikeCount,
		DownloadCount: d.DownloadCount,
	}

	if d.Desc != nil {
		dto.Desc = d.Desc.ResourceDesc()
	}
}

func (s datasetService) toDatasetSummaryDTO(d *domain.DatasetSummary, dto *DatasetSummaryDTO) {
	*dto = DatasetSummaryDTO{
		Id:            d.Id,
		Owner:         d.Owner.Account(),
		Name:          d.Name.ResourceName(),
		Tags:          d.Tags,
		UpdatedAt:     utils.ToDate(d.UpdatedAt),
		LikeCount:     d.LikeCount,
		DownloadCount: d.DownloadCount,
	}

	if d.Desc != nil {
		dto.Desc = d.Desc.ResourceDesc()
	}

}
