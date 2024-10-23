package app

import (
	"context"

	commonrepo "github.com/opensourceways/xihe-server/common/domain/repository"
	"github.com/opensourceways/xihe-server/domain"
	spacedomain "github.com/opensourceways/xihe-server/space/domain"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	"github.com/opensourceways/xihe-server/spaceapp/domain/repository"
	"golang.org/x/xerrors"
)

// SpaceappAppService is the interface for the space app service.
type SpaceappAppService interface {
	GetByName(context.Context, domain.Account, *spacedomain.SpaceIndex) (SpaceAppDTO, error)
	// GetBuildLog(context.Context, domain.Account, *spacedomain.SpaceIndex) (string, error)
	// GetBuildLogs(context.Context, domain.Account, *spacedomain.SpaceIndex) (BuildLogsDTO, error)
	// GetSpaceLog(context.Context, domain.Account, *spacedomain.SpaceIndex) (string, error)
}

// NewSpaceappAppService creates a new instance of the space app service.
func NewSpaceappAppService(
	repo repository.SpaceAppRepository,
	spaceRepo spacerepo.Project,
) *spaceappAppService {
	return &spaceappAppService{
		repo:      repo,
		spaceRepo: spaceRepo,
	}
}

// spaceappAppService
type spaceappAppService struct {
	repo      repository.SpaceAppRepository
	spaceRepo spacerepo.Project
}

// GetByName retrieves the space app by name.
func (s *spaceappAppService) GetByName(
	ctx context.Context, user domain.Account, index *spacedomain.SpaceIndex,
) (SpaceAppDTO, error) {
	var dto SpaceAppDTO

	space, err := s.spaceRepo.GetByName(user, index.Name)
	if err != nil {
		return dto, err
	}

	if space.IsPrivate() && space.Owner.Account() != index.Owner.Account() {
		return dto, commonrepo.NewErrorResourceNotExists(xerrors.New("not found"))
	}

	app, err := s.repo.FindBySpaceId(ctx, space.Id)
	if err == nil {
		return toSpaceAppDTO(&app), nil
	}

	if space.Hardware.IsNpu() && !space.CompPowerAllocated {
		return toSpaceNoCompQuotaDTO(&space), nil
	}

	if space.GitTemplate != "" {
		return toSpaceGitTemplateDTO(&space), nil
	}

	if commonrepo.IsErrorResourceNotExists(err) {
		err = newSpaceAppNotFound(xerrors.Errorf("not found, err: %w", err))
	} else {
		err = xerrors.Errorf("find space app by id failed, err: %w", err)
	}
	return dto, err
}
