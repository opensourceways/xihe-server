package repositoryimpl

import (
	"errors"
	"fmt"
	"strings"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/infrastructure/repositories"

	spacedomain "github.com/opensourceways/xihe-server/space/domain"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	"gorm.io/gorm/clause"
)

type projectAdapter struct {
	daoImpl
}

func (adapter *projectAdapter) Save(v *spacedomain.Project) (spacedomain.Project, error) {
	if v.Id != "" {
		err := errors.New("must be a new project")
		return spacedomain.Project{}, err
	}

	do := toProjectDO(v)
	err := adapter.db().Clauses(clause.Returning{}).Create(&do).Error
	if err != nil {
		return spacedomain.Project{}, err
	}

	doTags := toProjectTagsDO(v)
	if err := adapter.dbTag().CreateInBatches(doTags, 1000).Error; err != nil {
		return spacedomain.Project{}, err
	}
	return *v, nil
}

func (adapter *projectAdapter) Delete(index *domain.ResourceIndex) (err error) {
	return adapter.DeleteSingleRow(
		&projectDO{Id: index.Id, Owner: index.Owner.Account()},
	)
}

func (adapter *projectAdapter) GetByRepoId(id domain.Identity) (
	r spacedomain.Project, err error,
) {
	//filter
	do := projectDO{
		RepoId: id.Identity(),
	}

	// find project
	result := projectDO{}
	if err := adapter.daoImpl.GetProjectRecord(&do, &result); err != nil {
		return spacedomain.Project{}, err
	}

	// find tags
	var tagResults []projectTagsDO
	if err := adapter.daoImpl.dbTag().Where("project_id", id).Find(&tagResults).Error; err != nil {
		return spacedomain.Project{}, err
	}
	adapter.getProjectTags(&r, tagResults)

	// get datasets
	var datasetResults []datasetDO
	if err := adapter.daoImpl.dbDataset().Where("project_id", id).Find(&datasetResults).Error; err != nil {
		return spacedomain.Project{}, err
	}
	adapter.getDataset(&r, datasetResults)

	// get models
	var modelResults []modelDO
	if err := adapter.daoImpl.dbModel().Where("project_id", id).Find(&modelResults).Error; err != nil {
		return spacedomain.Project{}, err
	}
	adapter.getModel(&r, modelResults)

	return r, nil
}

func (adapter *projectAdapter) GetByName(owner domain.Account, name domain.ResourceName) (
	r spacedomain.Project, err error,
) {
	//filter
	do := projectDO{
		Owner: owner.Account(),
		Name:  name.ResourceName(),
	}

	// find project
	result := projectDO{}
	if err := adapter.daoImpl.GetProjectRecord(&do, &result); err != nil {
		return spacedomain.Project{}, err
	}

	id := result.RepoId
	if err = result.toProject(&r); err != nil {
		return spacedomain.Project{}, err
	}

	// find tags
	var tagResults []projectTagsDO
	if err := adapter.daoImpl.dbTag().Where("project_id", id).Find(&tagResults).Error; err != nil {
		return spacedomain.Project{}, err
	}
	adapter.getProjectTags(&r, tagResults)

	// get datasets
	var datasetResults []datasetDO
	if err := adapter.daoImpl.dbDataset().Where("project_id", id).Find(&datasetResults).Error; err != nil {
		return spacedomain.Project{}, err
	}
	adapter.getDataset(&r, datasetResults)

	// get models
	var modelResults []modelDO
	if err := adapter.daoImpl.dbModel().Where("project_id", id).Find(&modelResults).Error; err != nil {
		return spacedomain.Project{}, err
	}
	adapter.getModel(&r, modelResults)

	return r, nil
}

func (adapter *projectAdapter) getProjectTags(p *spacedomain.Project, tagResults []projectTagsDO) {
	p.Tags = make([]string, 0, len(tagResults))

	for _, tagDO := range tagResults {
		p.Tags = append(p.Tags, tagDO.TagName)
	}
}

func (adapter *projectAdapter) getDataset(p *spacedomain.Project, datasetResult []datasetDO) {
	if len(datasetResult) == 0 {
		return
	}

	relatedDatasets := make(domain.RelatedResources, len(datasetResult))

	for i, dataset := range datasetResult {
		relatedDatasets[i] = domain.ResourceIndex{
			Owner: domain.CreateAccount(dataset.Owner),
			Id:    dataset.DatasetId,
		}
	}

	p.RelatedDatasets = relatedDatasets

}

func (adapter *projectAdapter) getModel(p *spacedomain.Project, modelResult []modelDO) {
	if len(modelResult) == 0 {
		return
	}

	relatedModels := make(domain.RelatedResources, len(modelResult))

	for i, model := range modelResult {
		relatedModels[i] = domain.ResourceIndex{
			Owner: domain.CreateAccount(model.Owner),
			Id:    model.ModelId,
		}
	}

	p.RelatedModels = relatedModels

}

func (adapter *projectAdapter) FindUserProjects(opts []repository.UserResourceListOption) (
	[]spacedomain.ProjectSummary, error,
) {
	var projectSummaries []spacedomain.ProjectSummary

	for _, opt := range opts {
		var projects []projectDO
		query := adapter.db().Where("owner = ? AND id IN (?)", opt.Owner.Account(), opt.Ids).Order("updated_at DESC")

		err := query.Find(&projects).Error
		if err != nil {
			return nil, repositories.ConvertError(err)
		}

		for _, project := range projects {
			// ResourceName
			resourceName, err := domain.NewResourceName(project.Name)
			if err != nil {
				return nil, err
			}
			// ResourceDesc
			resourceDesc, err := domain.NewResourceDesc(project.Description)
			if err != nil {
				return nil, err
			}
			// ResourceTitle
			resourceTitle, err := domain.NewResourceTitle(project.Title)
			if err != nil {
				return nil, err
			}
			//projectCoverId
			resourceCoverId, err := domain.NewCoverId(project.CoverId)
			if err != nil {
				return nil, err
			}
			//projectHardware
			projectHardware, err := domain.NewHardware(project.Hardware, project.Type)
			if err != nil {
				return nil, err
			}
			// Tags
			var tagResults []projectTagsDO
			if err := adapter.dbTag().Where("project_id = ?", project.Id).Find(&tagResults).Error; err != nil {
				return nil, err
			}
			tags := make([]string, len(tagResults))
			for i, tag := range tagResults {
				tags[i] = tag.TagName
			}

			summary := spacedomain.ProjectSummary{
				Id:            project.Id,
				Owner:         domain.CreateAccount(project.Owner),
				Name:          resourceName,
				Desc:          resourceDesc,
				Title:         resourceTitle,
				Level:         domain.NewResourceLevelByNum(project.Level),
				CoverId:       resourceCoverId,
				UpdatedAt:     project.UpdatedAt,
				LikeCount:     project.LikeCount,
				ForkCount:     project.ForkCount,
				DownloadCount: project.DownloadCount,
				Hardware:      projectHardware,
				Tags:          tags,
			}

			projectSummaries = append(projectSummaries, summary)
		}
	}

	return projectSummaries, nil
}

func (adapter *projectAdapter) AddRelatedDataset(info *repository.RelatedResourceInfo) error {
	do := toDatasetDO(info)
	return adapter.dbDataset().Clauses(clause.Returning{}).Create(&do).Error
}

func (adapter *projectAdapter) RemoveRelatedDataset(info *repository.RelatedResourceInfo) error {
	do := toDatasetDO(info)

	if err := adapter.dbDataset().Where(&do).Delete(&do).Error; err != nil {
		return repositories.ConvertError(err)
	}

	return nil
}

func (adapter *projectAdapter) AddRelatedModel(info *repository.RelatedResourceInfo) error {
	do := toModelDO(info)
	return adapter.dbModel().Clauses(clause.Returning{}).Create(&do).Error
}

func (adapter *projectAdapter) RemoveRelatedModel(info *repository.RelatedResourceInfo) error {
	do := toModelDO(info)

	if err := adapter.dbModel().Where(&do).Delete(&do).Error; err != nil {
		return repositories.ConvertError(err)
	}

	return nil
}

func (adapter *projectAdapter) Get(owner domain.Account, identity string) (r spacedomain.Project, err error) {
	do := projectDO{Owner: owner.Account(), RepoId: identity}
	result := projectDO{}

	if err := adapter.daoImpl.GetProjectRecord(&do, &result); err != nil {
		return spacedomain.Project{}, err
	}

	if err = result.toProject(&r); err != nil {
		return spacedomain.Project{}, err
	}
	// find tags
	var tagResults []projectTagsDO
	if err := adapter.daoImpl.dbTag().Where("project_id", identity).Find(&tagResults).Error; err != nil {
		return spacedomain.Project{}, err
	}
	tags := make([]string, len(tagResults))
	for i, tag := range tagResults {
		tags[i] = tag.TagName
	}
	r.Tags = tags

	return

}

func (adapter *projectAdapter) ListGlobalAndSortByUpdateTime(
	option *repository.GlobalResourceListOption,
) (spacerepo.UserProjectsInfo, error) {
	return adapter.listGlobal(
		option, adapter.listGlobalAndSortByUpdateTime,
	)
}

func (adapter *projectAdapter) listGlobalAndSortByUpdateTime(
	do *repositories.GlobalResourceListDO) ([]ProjectSummaryDO, int, error) {
	var items []projectDO
	var count int64

	// 基础查询条件
	baseQuery := adapter.db()

	// 计算总数
	if err := baseQuery.Model(&projectDO{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 构建分页查询
	query := baseQuery.Order("updated_at DESC")
	if do.PageNum > 0 && do.CountPerPage > 0 {
		query = query.Limit(int(do.CountPerPage)).Offset((int(do.PageNum) - 1) * int(do.CountPerPage))
	}

	// 执行分页查询
	err := query.Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为 ProjectSummaryDO
	var projectSummaries []ProjectSummaryDO
	for _, item := range items {
		summary := ProjectSummaryDO{
			Id:            item.Id,
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
		}
		// 查询标签
		var tagResults []projectTagsDO
		if err := adapter.dbTag().Where("project_id = ?", item.Id).Find(&tagResults).Error; err != nil {
			return nil, 0, err
		}
		tags := make([]string, len(tagResults))
		for i, tag := range tagResults {
			tags[i] = tag.TagName
		}
		summary.Tags = tags

		projectSummaries = append(projectSummaries, summary)
	}

	return projectSummaries, int(count), nil
}
func (adapter *projectAdapter) ListGlobalAndSortByDownloadCount(
	option *repository.GlobalResourceListOption,
) (spacerepo.UserProjectsInfo, error) {
	return adapter.listGlobal(
		option, adapter.listGlobalAndSortByDownloadCount,
	)
}

func (adapter *projectAdapter) listGlobalAndSortByDownloadCount(
	do *repositories.GlobalResourceListDO) ([]ProjectSummaryDO, int, error) {
	var items []projectDO
	var count int64

	// 基础查询条件
	baseQuery := adapter.db()

	// 计算总数
	if err := baseQuery.Model(&projectDO{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 构建分页查询
	query := baseQuery.Order("download_count ASC")
	if do.PageNum > 0 && do.CountPerPage > 0 {
		query = query.Limit(int(do.CountPerPage)).Offset((int(do.PageNum) - 1) * int(do.CountPerPage))
	}

	// 执行分页查询
	err := query.Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为 ProjectSummaryDO
	var projectSummaries []ProjectSummaryDO
	for _, item := range items {
		summary := ProjectSummaryDO{
			Id:            item.Id,
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
		}
		// 查询标签
		var tagResults []projectTagsDO
		if err := adapter.dbTag().Where("project_id = ?", item.Id).Find(&tagResults).Error; err != nil {
			return nil, 0, err
		}
		tags := make([]string, len(tagResults))
		for i, tag := range tagResults {
			tags[i] = tag.TagName
		}
		summary.Tags = tags

		projectSummaries = append(projectSummaries, summary)
	}

	return projectSummaries, int(count), nil
}
func (adapter *projectAdapter) ListGlobalAndSortByFirstLetter(
	option *repository.GlobalResourceListOption,
) (spacerepo.UserProjectsInfo, error) {
	return adapter.listGlobal(
		option, adapter.listGlobalAndSortByFirstLetter,
	)
}

func (adapter *projectAdapter) listGlobalAndSortByFirstLetter(
	do *repositories.GlobalResourceListDO) ([]ProjectSummaryDO, int, error) {
	var items []projectDO
	var count int64

	// 基础查询条件
	baseQuery := adapter.db()

	// 计算总数
	if err := baseQuery.Model(&projectDO{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// 构建分页查询
	query := baseQuery.Order("LOWER(name) COLLATE \"C\" ASC")
	if do.PageNum > 0 && do.CountPerPage > 0 {
		query = query.Limit(int(do.CountPerPage)).Offset((int(do.PageNum) - 1) * int(do.CountPerPage))
	}

	// 执行分页查询
	err := query.Find(&items).Error
	if err != nil {
		return nil, 0, err
	}

	// 转换为 ProjectSummaryDO
	var projectSummaries []ProjectSummaryDO
	for _, item := range items {
		summary := ProjectSummaryDO{
			Id:            item.Id,
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
		}
		// 查询标签
		var tagResults []projectTagsDO
		if err := adapter.dbTag().Where("project_id = ?", item.Id).Find(&tagResults).Error; err != nil {
			return nil, 0, err
		}
		tags := make([]string, len(tagResults))
		for i, tag := range tagResults {
			tags[i] = tag.TagName
		}
		summary.Tags = tags

		projectSummaries = append(projectSummaries, summary)
	}

	return projectSummaries, int(count), nil
}

func (adapter *projectAdapter) listGlobal(
	option *repository.GlobalResourceListOption,
	f func(*repositories.GlobalResourceListDO) ([]ProjectSummaryDO, int, error),
) (
	info spacerepo.UserProjectsInfo, err error,
) {
	return adapter.doList(func() ([]ProjectSummaryDO, int, error) {
		do := repositories.ToGlobalResourceListDO(option)

		return f(&do)
	})
}

func (adapter *projectAdapter) GetSummary(owner domain.Account, projectId string) (
	r spacerepo.ProjectSummary, err error,
) {
	v, err := adapter.getSummary(owner.Account(), projectId)
	if err != nil {
		err = repositories.ConvertError(err)

		return
	}

	if r.ResourceSummary, err = v.ToProject(); err == nil {
		r.Tags = v.Tags
	}

	return
}

func (adapter *projectAdapter) getSummary(owner string, projectId string) (
	do ProjectResourceSummaryDO, err error,
) {
	//filter
	filter := projectDO{
		Owner:  owner,
		RepoId: projectId,
	}

	// find project
	project := projectDO{}
	if err := adapter.daoImpl.GetProjectRecord(&filter, &project); err != nil {
		return ProjectResourceSummaryDO{}, err
	}

	// find tags
	var tagResults []projectTagsDO
	if err := adapter.daoImpl.dbTag().Where("project_id", projectId).Find(&tagResults).Error; err != nil {
		return ProjectResourceSummaryDO{}, err
	}
	// Convert tags to a slice of strings
	tags := make([]string, len(tagResults))
	for i, tag := range tagResults {
		tags[i] = tag.TagName
	}

	// Convert projectDO to ProjectResourceSummaryDO
	do = ProjectResourceSummaryDO{
		ResourceSummaryDO: repositories.ResourceSummaryDO{
			Owner:    project.Owner,
			Name:     project.Name,
			Id:       project.Id,
			RepoId:   project.RepoId,
			RepoType: project.RepoType,
		},
		Tags: tags,
	}

	return do, nil
}

func (adapter *projectAdapter) GetSummaryByName(owner domain.Account, name domain.ResourceName) (
	domain.ResourceSummary, error,
) {
	v, err := adapter.getSummaryByName(owner.Account(), name.ResourceName())
	if err != nil {
		return domain.ResourceSummary{}, err
	}

	return v.ToProject()
}

func (adapter *projectAdapter) getSummaryByName(owner, name string) (
	do repositories.ResourceSummaryDO, err error,
) {
	//filter
	filter := projectDO{
		Owner: owner,
		Name:  name,
	}

	// find project
	project := projectDO{}
	if err := adapter.daoImpl.GetProjectRecord(&filter, &project); err != nil {
		return repositories.ResourceSummaryDO{}, err
	}

	// Convert projectDO to ProjectResourceSummaryDO
	do = repositories.ResourceSummaryDO{
		Owner:    project.Owner,
		Name:     project.Name,
		Id:       project.Id,
		RepoId:   project.RepoId,
		RepoType: project.RepoType,
	}

	return do, nil

}

func (adapter *projectAdapter) AddLike(p *domain.ResourceIndex) error {
	filter := projectDO{
		Owner: p.Owner.Account(),
		Id:    p.Id,
	}

	err := adapter.daoImpl.IncrementStatistic(&filter, fieldLikeCount, 1)
	if err != nil {
		return err
	}

	return nil
}

func (adapter *projectAdapter) RemoveLike(p *domain.ResourceIndex) error {
	filter := projectDO{
		Owner: p.Owner.Account(),
		Id:    p.Id,
	}

	if err := adapter.daoImpl.IncrementStatistic(&filter, fieldLikeCount, -1); err != nil {
		return err
	}

	return nil
}

func (adapter *projectAdapter) Search(option *repository.ResourceSearchOption) (
	r repository.ResourceSearchResult, err error) {
	var projectDOs []projectDO

	// 构建查询
	query := adapter.db().Model(&projectDO{})

	// 添加名称搜索条件
	if option.Name != "" {
		query = query.Where("name LIKE ?", "%"+strings.TrimSpace(option.Name)+"%")
	}

	// 添加RepoType搜索条件
	if len(option.RepoType) > 0 {
		repoTypes := make([]string, len(option.RepoType))
		for i, rt := range option.RepoType {
			repoTypes[i] = rt.RepoType()
		}
		query = query.Where("repo_type IN (?)", repoTypes)
	}

	// 执行查询
	if err = query.Limit(option.TopNum).Find(&projectDOs).Error; err != nil {
		return r, err
	}

	// 计算总数
	var total int64
	if err = query.Count(&total).Error; err != nil {
		return r, err
	}
	r.Total = int(total)

	// 转换结果
	r.Top = make([]domain.ResourceSummary, len(projectDOs))
	for i, do := range projectDOs {
		if r.Top[i].Owner, err = domain.NewAccount(do.Owner); err != nil {
			return r, err
		}

		if r.Top[i].Name, err = domain.NewResourceName(do.Name); err != nil {
			return r, err
		}

		if r.Top[i].RepoType, err = domain.NewRepoType(do.RepoType); err != nil {
			return r, err
		}

		r.Top[i].Id = do.Id
		r.Top[i].RepoId = do.RepoId
	}

	return r, nil
}

func (adapter *projectAdapter) UpdateProperty(info *spacerepo.ProjectPropertyUpdateInfo) error {
	p := &info.Property

	fmt.Printf("==================p.Level: %+v\n", p.Level)

	do := projectDO{
		Id:          info.Id,
		Owner:       info.Owner.Account(),
		Version:     info.Version,
		UpdatedAt:   info.UpdatedAt,
		Name:        p.Name.ResourceName(),
		FL:          p.Name.FirstLetterOfName(),
		Description: p.Desc.ResourceDesc(),
		Title:       p.Title.ResourceTitle(),
		Type:        p.RepoType.RepoType(),
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

	result := adapter.db().Model(&projectDO{}).Where("id = ?", do.Id).Updates(do)
	if result.Error != nil {
		return repositories.ConvertError(result.Error)
	}

	if result.RowsAffected == 0 {
		return repository.NewErrorResourceNotExists(errors.New("project not found"))
	}

	// 删除旧的标签关联
	if err := adapter.dbTag().Where("project_id = ?", do.Id).Delete(&projectTagsDO{}).Error; err != nil {
		return repositories.ConvertError(err)
	}

	// 创建新的标签关联
	var newTagsDOs []projectTagsDO
	for _, tagName := range p.Tags {
		newTagsDOs = append(newTagsDOs, projectTagsDO{
			ProjectId: do.Id,
			TagName:   tagName,
		})
	}

	if err := adapter.dbTag().Create(&newTagsDOs).Error; err != nil {
		return repositories.ConvertError(err)
	}
	fmt.Printf("======================newTagsDOs: %+v\n", newTagsDOs)
	return nil
}

func (adapter *projectAdapter) IncreaseFork(r *domain.ResourceIndex) error {
	filter := projectDO{
		Owner: r.Owner.Account(),
		Id:    r.Id,
	}

	err := adapter.daoImpl.IncrementStatistic(&filter, fieldForkCount, 1)
	if err != nil {
		return repositories.ConvertError(err)
	}

	return nil
}

func (adapter *projectAdapter) IncreaseDownload(r *domain.ResourceIndex) error {
	filter := projectDO{
		Owner: r.Owner.Account(),
		Id:    r.Id,
	}

	err := adapter.daoImpl.IncrementStatistic(&filter, fieldDownload, 1)
	if err != nil {
		return repositories.ConvertError(err)
	}

	return nil
}
