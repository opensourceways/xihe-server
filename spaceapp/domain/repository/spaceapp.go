package repository

import (
	"github.com/opensourceways/xihe-server/domain"
	spaceappdomain "github.com/opensourceways/xihe-server/spaceapp/domain"
)

type SpaceAppRepository interface {
	SaveWithoutAllBuildLog(*spaceappdomain.SpaceApp) error
	FindBySpaceId(domain.Identity) (spaceappdomain.SpaceApp, error)
	FindById(domain.Identity) (spaceappdomain.SpaceApp, error)
}
