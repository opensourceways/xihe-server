package repository

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/spaceapp/domain"
)

type SpaceAppRepository interface {
	Save(*domain.SpaceApp) error
	FindBySpaceId(types.Identity) (domain.SpaceApp, error)
	FindById(types.Identity) (domain.SpaceApp, error)
	SaveWithBuildLog(*domain.SpaceApp, *domain.SpaceAppBuildLog) error
}
