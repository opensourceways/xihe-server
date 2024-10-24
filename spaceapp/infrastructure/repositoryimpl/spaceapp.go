package repositoryimpl

import (
	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	"github.com/opensourceways/xihe-server/domain"
	spaceappdomain "github.com/opensourceways/xihe-server/spaceapp/domain"
	"github.com/opensourceways/xihe-server/spaceapp/domain/repository"
)

const (
	tableSpaceApp = "space_app"

	fieldAllBuildLog = "all_build_log"
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

// Save saves space application into repository without all build log
func (impl spaceAppRepoImpl) SaveWithoutAllBuildLog(app *spaceappdomain.SpaceApp) error {
	do := toSpaceAppDO(app)
	do.Version += 1

	return impl.dao.UpdateWithOmittingSpecificFields(
		&spaceappDO{Id: do.Id, Version: do.Version}, &do, fieldAllBuildLog,
	)
}

// FindBySpaceId finds a space application in the repository based on the space ID.
func (impl spaceAppRepoImpl) FindBySpaceId(id domain.Identity) (spaceappdomain.SpaceApp, error) {
	do := spaceappDO{SpaceId: id.Integer()}

	// It must new a new DO, otherwise the sql statement will include duplicate conditions.
	result := spaceappDO{}

	if err := impl.dao.GetRecord(&do, &result); err != nil {
		return spaceappdomain.SpaceApp{}, err
	}

	return result.toSpaceApp(), nil
}

func (impl spaceAppRepoImpl) FindById(id domain.Identity) (spaceappdomain.SpaceApp, error) {
	do := spaceappDO{Id: id.Integer()}

	if err := impl.dao.GetByPrimaryKey(&do); err != nil {
		return spaceappdomain.SpaceApp{}, err
	}

	return do.toSpaceApp(), nil
}
