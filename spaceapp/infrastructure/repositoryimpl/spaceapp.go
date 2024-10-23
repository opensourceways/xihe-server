package repositoryimpl

import (
	"context"

	"github.com/opensourceways/xihe-server/common/infrastructure/pgsql"
	"github.com/opensourceways/xihe-server/spaceapp/domain"
	"github.com/opensourceways/xihe-server/spaceapp/domain/repository"
)

const (
	tableSpaceApp = "spaceapp"

	fieldAllBuildLog = "all_build_log"
)

func NewSpaceAppRepository() repository.SpaceAppRepository {
	return spaceAppRepoImpl{dao: pgsql.NewDBTable(tableSpaceApp)}
}

type spaceAppRepoImpl struct {
	dao SpaceAppDAO
}

// Save saves space application into repository without all build log
func (impl spaceAppRepoImpl) SaveWithoutAllBuildLog(ctx context.Context, app *domain.SpaceApp) error {
	do := toSpaceAppDO(app)
	do.Version += 1

	return impl.dao.UpdateWithOmitingSpecificFields(
		&spaceappDO{Id: do.Id, Version: do.Version}, &do, fieldAllBuildLog,
	)
}

// FindBySpaceId finds a space application in the repository based on the space ID.
func (impl spaceAppRepoImpl) FindBySpaceId(
	ctx context.Context, id string) (domain.SpaceApp, error) {
	do := spaceappDO{SpaceId: id}

	// It must new a new DO, otherwise the sql statement will include duplicate conditions.
	result := spaceappDO{}

	if err := impl.dao.GetRecord(ctx, &do, &result); err != nil {
		return domain.SpaceApp{}, err
	}

	return result.toSpaceApp(), nil
}
