package repositoryimpl

import (
	"errors"

	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/spaceapp/domain"
	"github.com/opensourceways/xihe-server/spaceapp/domain/repository"
)

const (
	tableSpaceApp = "space_app"

	fieldAllBuildLog = "all_build_log"
	fieldSpaceId     = "space_id"
)

func NewSpaceAppRepository() (repository.SpaceAppRepository, error) {
	if err := pgsql.AutoMigrate(&spaceappDO{}); err != nil {
		return nil, err
	}

	return spaceAppRepoImpl{dao: pgsql.NewDBTable(tableSpaceApp)}, nil
}

type spaceAppRepoImpl struct {
	dao SpaceAppDAO
}

// Add adds a space application to the repository.
func (adapter spaceAppRepoImpl) Add(m *domain.SpaceApp) error {
	do := toSpaceAppDO(m)

	err := adapter.dao.DB().Create(&do).Error

	if err != nil && adapter.dao.IsRecordExists(err) {
		return errors.New("space app exists")

	}

	return err
}

// Save saves space application into repository without all build log
func (impl spaceAppRepoImpl) Save(app *domain.SpaceApp) error {
	do := toSpaceAppDO(app)
	do.Version += 1

	return impl.dao.UpdateWithOmittingSpecificFields(
		&spaceappDO{Id: do.Id, Version: do.Version}, &do, fieldAllBuildLog,
	)
}

// FindBySpaceId finds a space application in the repository based on the space ID.
func (impl spaceAppRepoImpl) FindBySpaceId(id types.Identity) (domain.SpaceApp, error) {
	do := spaceappDO{SpaceId: id.Integer()}

	// It must new a new DO, otherwise the sql statement will include duplicate conditions.
	result := spaceappDO{}

	if err := impl.dao.GetRecord(&do, &result); err != nil {
		return domain.SpaceApp{}, err
	}

	return result.toSpaceApp(), nil
}

func (impl spaceAppRepoImpl) FindById(id types.Identity) (domain.SpaceApp, error) {
	do := spaceappDO{Id: id.Integer()}

	if err := impl.dao.GetByPrimaryKey(&do); err != nil {
		return domain.SpaceApp{}, err
	}

	return do.toSpaceApp(), nil
}

// SaveWithBuildLog saves a space application and build log in the repository.
func (impl spaceAppRepoImpl) SaveWithBuildLog(m *domain.SpaceApp, log *domain.SpaceAppBuildLog) error {
	do := toSpaceAppDO(m)
	do.Version += 1
	do.AllBuildLog = log.Logs

	return impl.dao.Update(&spaceappDO{Id: do.Id, Version: m.Version}, &do)
}

func (adapter spaceAppRepoImpl) Remove(spaceId types.Identity) error {
	return adapter.dao.DB().Where(
		adapter.dao.EqualQuery(fieldSpaceId), spaceId.Identity(),
	).Delete(
		spaceappDO{},
	).Error
}
