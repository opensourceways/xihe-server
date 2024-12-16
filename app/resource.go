package app

import (
	"errors"
	"fmt"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	spacedomain "github.com/opensourceways/xihe-server/space/domain"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	userdomain "github.com/opensourceways/xihe-server/user/domain"
	userrepo "github.com/opensourceways/xihe-server/user/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type ResourceDTO struct {
	Owner struct {
		Name     string `json:"name"`
		AvatarId string `json:"avatar_id"`
	} `json:"owner"`

	Id            string   `json:"id"`
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	Desc          string   `json:"description"`
	Title         string   `json:"title"`
	CoverId       string   `json:"cover_id"`
	UpdateAt      string   `json:"update_at"`
	Tags          []string `json:"tags"`
	ResourceLevel string   `json:"level"`

	LikeCount     int `json:"like_count"`
	ForkCount     int `json:"fork_count"`
	DownloadCount int `json:"download_count"`
}

func (r *ResourceDTO) identity() string {
	return fmt.Sprintf("%s_%s_%s", r.Owner.Name, r.Type, r.Id)
}

type ResourceListCmd struct {
	repository.ResourceListOption

	SortType domain.SortType
}

func (cmd *ResourceListCmd) ToResourceListOption() repository.ResourceListOption {
	return cmd.ResourceListOption
}

type ResourceService struct {
	User      userrepo.User
	Model     repository.Model
	ProjectPg spacerepo.ProjectPg
	Dataset   repository.Dataset
}

func (s ResourceService) list(resources []*domain.ResourceObject) (
	dtos []ResourceDTO, err error,
) {
	users, projects, datasets, models := s.toOptions(resources)

	return s.listResources(users, projects, datasets, models, len(resources))
}

func (s ResourceService) listProjects(resources []domain.ResourceIndex) (
	dtos []ResourceDTO, err error,
) {
	if len(resources) == 0 {
		return
	}

	users, options := s.singleResourceOptions(resources)

	return s.listResources(users, options, nil, nil, len(resources))
}

func (s ResourceService) ListModels(resources []domain.ResourceIndex) (
	dtos []ResourceDTO, err error,
) {
	if len(resources) == 0 {
		return
	}

	users, options := s.singleResourceOptions(resources)

	return s.listResources(users, nil, nil, options, len(resources))
}

func (s ResourceService) ListDatasets(resources []domain.ResourceIndex) (
	dtos []ResourceDTO, err error,
) {
	if len(resources) == 0 {
		return
	}

	users, options := s.singleResourceOptions(resources)

	return s.listResources(users, nil, options, nil, len(resources))
}

func (s ResourceService) singleResourceOptions(resources []domain.ResourceIndex) (
	users []userdomain.Account,
	options []repository.UserResourceListOption,
) {
	ul := make(map[string]userdomain.Account)
	ro := make(map[string][]string)

	for i := range resources {
		item := resources[i]

		account := item.Owner.Account()
		if _, ok := ul[account]; !ok {
			ul[account] = item.Owner
		}

		s.store(&item, ro)
	}

	options = s.toOptionList(ro, ul)
	users = s.userMapToList(ul)

	return
}

func (s ResourceService) toOptions(resources []*domain.ResourceObject) (
	users []userdomain.Account,
	projects []repository.UserResourceListOption,
	datasets []repository.UserResourceListOption,
	models []repository.UserResourceListOption,
) {
	ul := make(map[string]userdomain.Account)
	po := make(map[string][]string)
	do := make(map[string][]string)
	mo := make(map[string][]string)

	for i := range resources {
		item := resources[i]

		account := item.Owner.Account()
		if _, ok := ul[account]; !ok {
			ul[account] = item.Owner
		}

		switch item.Type.ResourceType() {
		case domain.ResourceTypeProject.ResourceType():
			s.store(&item.ResourceIndex, po)

		case domain.ResourceTypeModel.ResourceType():
			s.store(&item.ResourceIndex, mo)

		case domain.ResourceTypeDataset.ResourceType():
			s.store(&item.ResourceIndex, do)
		}
	}

	projects = s.toOptionList(po, ul)
	datasets = s.toOptionList(do, ul)
	models = s.toOptionList(mo, ul)
	users = s.userMapToList(ul)

	return
}

func (s ResourceService) listResources(
	users []userdomain.Account,
	projects []repository.UserResourceListOption,
	datasets []repository.UserResourceListOption,
	models []repository.UserResourceListOption,
	total int,
) (
	dtos []ResourceDTO, err error,
) {
	allUsers, err := s.User.FindUsersInfo(users)
	if err != nil {
		return
	}

	userInfos := make(map[string]*userdomain.UserInfo)
	for i := range allUsers {
		item := &allUsers[i]
		userInfos[item.Account.Account()] = item
	}

	dtos = make([]ResourceDTO, total)
	total = 0
	r := dtos

	if len(projects) > 0 {
		all, err := s.ProjectPg.FindUserProjects(projects)
		if err != nil {
			return nil, err
		}

		n := len(all)
		if n > 0 {
			if len(r) < n {
				return nil, errors.New("unmatched size")
			}
			s.projectToResourceDTO(userInfos, all, r)
			r = r[n:]
			total += n
		}
	}

	if len(models) > 0 {
		all, err := s.Model.FindUserModels(models)
		if err != nil {
			return nil, err
		}

		n := len(all)
		if n > 0 {
			if len(r) < n {
				return nil, errors.New("unmatched size")
			}
			s.modelToResourceDTO(userInfos, all, r)
			r = r[n:]
			total += n
		}
	}

	if n := len(datasets); n > 0 {
		all, err := s.Dataset.FindUserDatasets(datasets)
		if err != nil {
			return nil, err
		}

		n := len(all)
		if n > 0 {
			if len(r) < n {
				return nil, errors.New("unmatched size")
			}
			s.datasetToResourceDTO(userInfos, all, r)
			total += n
		}
	}

	dtos = dtos[:total]

	return
}

func (s ResourceService) projectToResourceDTO(
	userInfos map[string]*userdomain.UserInfo,
	projects []spacedomain.ProjectSummary, dtos []ResourceDTO,
) {
	for i := range projects {
		p := &projects[i]

		v := ResourceDTO{
			Id:            p.Id,
			Name:          p.Name.ResourceName(),
			Type:          domain.ResourceTypeProject.ResourceType(),
			CoverId:       p.CoverId.CoverId(),
			Tags:          p.Tags,
			UpdateAt:      utils.ToDate(p.UpdatedAt),
			LikeCount:     p.LikeCount,
			ForkCount:     p.ForkCount,
			DownloadCount: p.DownloadCount,
		}

		if p.Desc != nil {
			v.Desc = p.Desc.ResourceDesc()
		}

		if p.Title != nil {
			v.Title = p.Title.ResourceTitle()
		}

		if p.Level != nil {
			v.ResourceLevel = p.Level.ResourceLevel()
		}

		if u, ok := userInfos[p.Owner.Account()]; ok {
			v.Owner.Name = u.Account.Account()
			if u.AvatarId != nil {
				v.Owner.AvatarId = u.AvatarId.AvatarId()
			}
		}

		dtos[i] = v
	}
}

func (s ResourceService) modelToResourceDTO(
	userInfos map[string]*userdomain.UserInfo,
	data []domain.ModelSummary, dtos []ResourceDTO,
) {
	for i := range data {
		d := &data[i]

		v := ResourceDTO{
			Id:            d.Id,
			Name:          d.Name.ResourceName(),
			Tags:          d.Tags,
			Type:          domain.ResourceTypeModel.ResourceType(),
			UpdateAt:      utils.ToDate(d.UpdatedAt),
			LikeCount:     d.LikeCount,
			DownloadCount: d.DownloadCount,
		}

		if d.Desc != nil {
			v.Desc = d.Desc.ResourceDesc()
		}

		if d.Title != nil {
			v.Title = d.Title.ResourceTitle()
		}

		if u, ok := userInfos[d.Owner.Account()]; ok {
			v.Owner.Name = u.Account.Account()
			if u.AvatarId != nil {
				v.Owner.AvatarId = u.AvatarId.AvatarId()
			}
		}

		dtos[i] = v
	}
}

func (s ResourceService) datasetToResourceDTO(
	userInfos map[string]*userdomain.UserInfo,
	data []domain.DatasetSummary, dtos []ResourceDTO,
) {
	for i := range data {
		d := &data[i]

		v := ResourceDTO{
			Id:            d.Id,
			Name:          d.Name.ResourceName(),
			Tags:          d.Tags,
			Type:          domain.ResourceTypeDataset.ResourceType(),
			UpdateAt:      utils.ToDate(d.UpdatedAt),
			LikeCount:     d.LikeCount,
			DownloadCount: d.DownloadCount,
		}

		if d.Desc != nil {
			v.Desc = d.Desc.ResourceDesc()
		}

		if d.Title != nil {
			v.Title = d.Title.ResourceTitle()
		}

		if u, ok := userInfos[d.Owner.Account()]; ok {
			v.Owner.Name = u.Account.Account()
			if u.AvatarId != nil {
				v.Owner.AvatarId = u.AvatarId.AvatarId()
			}
		}

		dtos[i] = v
	}
}

func (s ResourceService) store(v *domain.ResourceIndex, m map[string][]string) {
	a := v.Owner.Account()

	if p, ok := m[a]; !ok {
		m[a] = []string{v.Id}
	} else {
		m[a] = append(p, v.Id)
	}
}

func (s ResourceService) toOptionList(
	m map[string][]string, users map[string]userdomain.Account,
) []repository.UserResourceListOption {

	if len(m) == 0 {
		return nil
	}

	r := make([]repository.UserResourceListOption, len(m))

	i := 0
	for k, v := range m {
		r[i] = repository.UserResourceListOption{
			Owner: users[k],
			Ids:   v,
		}

		i++
	}

	return r
}

func (s ResourceService) userMapToList(m map[string]userdomain.Account) []userdomain.Account {
	if len(m) == 0 {
		return nil
	}

	r := make([]userdomain.Account, len(m))

	i := 0
	for _, u := range m {
		r[i] = u
		i++
	}

	return r
}

func (s ResourceService) FindUserAvater(users []userdomain.Account) ([]string, error) {
	allUsers, err := s.User.FindUsersInfo(users)
	if err != nil {
		return nil, err
	}

	userInfos := make(map[string]string)
	for i := range allUsers {
		if item := &allUsers[i]; item.AvatarId != nil {
			userInfos[item.Account.Account()] = item.AvatarId.AvatarId()
		}
	}

	r := make([]string, len(users))
	for i := range users {
		r[i] = userInfos[users[i].Account()]
	}

	return r, nil
}

func (s ResourceService) CanApplyResourceName(owner domain.Account, name domain.ResourceName) bool {
	if _, err := s.ProjectPg.GetSummaryByName(owner, name); err == nil {
		return false
	}

	if _, err := s.Model.GetSummaryByName(owner, name); err == nil {
		return false
	}

	if _, err := s.Dataset.GetSummaryByName(owner, name); err == nil {
		return false
	}

	return true
}

func (s ResourceService) IsPrivate(owner domain.Account, resourceType domain.ResourceType, id string,
) (isprivate bool, ok bool) {
	switch resourceType.ResourceType() {
	case domain.ResourceProject:
		if p, err := s.ProjectPg.Get(owner, id); err == nil {
			return p.IsPrivate(), true
		}
	case domain.ResourceModel:
		if m, err := s.Model.Get(owner, id); err == nil {
			return m.IsPrivate(), true
		}
	case domain.ResourceDataset:
		if d, err := s.Dataset.Get(owner, id); err == nil {
			return d.IsPrivate(), true
		}
	default:
		return false, false
	}

	return false, false
}
