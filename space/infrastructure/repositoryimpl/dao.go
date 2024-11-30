package repositoryimpl

import (
	"errors"

	"gorm.io/gorm"

	"github.com/opensourceways/xihe-server/common/domain/repository"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

var dbInstance *gorm.DB

type daoImpl struct {
	table        string
	tableTag     string
	tableDataset string
	tableModel   string
}

// Each operation must generate a new gorm.DB instance.
// If using the same gorm.DB instance by different operations, they will share the same error.
func (dao *daoImpl) db() *gorm.DB {
	if dbInstance == nil {
		return nil
	}

	return dbInstance.Table(dao.table)
}

func (dao *daoImpl) dbTag() *gorm.DB {
	if dbInstance == nil {
		return nil
	}

	return dbInstance.Table(dao.tableTag)
}

func (dao *daoImpl) dbDataset() *gorm.DB {
	if dbInstance == nil {
		return nil
	}

	return dbInstance.Table(dao.tableDataset)
}

func (dao *daoImpl) dbModel() *gorm.DB {
	if dbInstance == nil {
		return nil
	}

	return dbInstance.Table(dao.tableModel)
}

// GetRecord retrieves a single record from the database based on the provided filter
// and stores it in the result parameter.
func (dao *daoImpl) GetProjectRecord(filter, result interface{}) error {
	err := dao.db().Where(filter).First(result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

// DeleteByPrimaryKey deletes a single record from the database based on the primary key of the row parameter.
func (dao *daoImpl) DeleteSingleRow(row interface{}) error {
	err := dao.db().Delete(row).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return repository.NewErrorResourceNotExists(errors.New("not found"))
	}

	return err
}

func (dao *daoImpl) IncrementStatistic(filter *projectDO, fieldName string, increment int) error {
	result := dao.db().Model(&projectDO{}).
		Where(equalQuery(fieldOwner), filter.Owner).
		Where(equalQuery(fieldID), filter.Id).
		UpdateColumn(fieldName, gorm.Expr(fieldName+" + ?", increment))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return repository.NewErrorResourceNotExists(errors.New("project not found"))
	}

	return nil
}

func (dao *daoImpl) ListAndSortByUpdateTime(
	owner string, do *repositories.ResourceListDO,
) ([]ProjectSummaryDO, int, error) {
	var items []projectDO
	var count int64

	// 基础查询条件
	baseQuery := dao.db().Where(equalQuery(fieldOwner), owner)

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
		if err := dao.dbTag().Where(equalQuery(fieldProjectId), item.Id).Find(&tagResults).Error; err != nil {
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

func (dao *daoImpl) ListAndSortByFirstLetter(
	owner string, do *repositories.ResourceListDO,
) ([]ProjectSummaryDO, int, error) {
	var items []projectDO
	var count int64

	// 基础查询条件
	baseQuery := dao.db().Where(equalQuery(fieldOwner), owner)

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
		if err := dao.dbTag().Where(equalQuery(fieldProjectId), item.Id).Find(&tagResults).Error; err != nil {
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

func (dao *daoImpl) ListAndSortByDownloadCount(
	owner string, do *repositories.ResourceListDO,
) ([]ProjectSummaryDO, int, error) {
	var items []projectDO
	var count int64

	// 基础查询条件
	baseQuery := dao.db().Where(equalQuery(fieldOwner), owner)

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
		if err := dao.dbTag().Where(equalQuery(fieldProjectId), item.Id).Find(&tagResults).Error; err != nil {
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
