package repository

import (
	"context"

	"github.com/opensourceways/xihe-server/spaceapp/domain"
)

type SpaceAppRepository interface {
	SaveWithoutAllBuildLog(context.Context, *domain.SpaceApp) error
	FindBySpaceId(context.Context, string) (domain.SpaceApp, error)
}
