package app

import (
	"context"
	"errors"

	"github.com/opensourceways/xihe-server/common/domain/allerror"
	commonrepo "github.com/opensourceways/xihe-server/common/domain/repository"
	"github.com/opensourceways/xihe-server/domain"
	spacedomain "github.com/opensourceways/xihe-server/space/domain"
	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	spaceappdomain "github.com/opensourceways/xihe-server/spaceapp/domain"
	spacemesage "github.com/opensourceways/xihe-server/spaceapp/domain/message"
	"github.com/opensourceways/xihe-server/spaceapp/domain/repository"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

// SpaceappAppService is the interface for the space app service.
type SpaceappAppService interface {
	GetByName(context.Context, domain.Account, *spacedomain.SpaceIndex) (SpaceAppDTO, error)
	GetBuildLog(context.Context, domain.Account, *spacedomain.SpaceIndex) (string, error)
	GetBuildLogs(context.Context, domain.Account, *spacedomain.SpaceIndex) (BuildLogsDTO, error)
	GetRequestDataStream(*spaceappdomain.SeverSentStream) error
	GetSpaceLog(context.Context, domain.Account, *spacedomain.SpaceIndex) (string, error)
	CheckPermissionRead(context.Context, domain.Account, *spacedomain.SpaceIndex) error
}

// NewSpaceappAppService creates a new instance of the space app service.
func NewSpaceappAppService(
	repo repository.SpaceAppRepository,
	spaceRepo spacerepo.Project,
	sse spaceappdomain.SeverSentEvent,
	spacesender spacemesage.SpaceAppMessageProducer,
) *spaceappAppService {
	return &spaceappAppService{
		repo:        repo,
		spaceRepo:   spaceRepo,
		sse:         sse,
		spacesender: spacesender,
	}
}

// spaceappAppService
type spaceappAppService struct {
	repo        repository.SpaceAppRepository
	spaceRepo   spacerepo.Project
	sse         spaceappdomain.SeverSentEvent
	spacesender spacemesage.SpaceAppMessageProducer
}

// GetByName retrieves the space app by name.
func (s *spaceappAppService) GetByName(
	ctx context.Context, user domain.Account, index *spacedomain.SpaceIndex,
) (SpaceAppDTO, error) {
	var dto SpaceAppDTO

	space, err := s.spaceRepo.GetByName(index.Owner, index.Name)
	if err != nil {
		logrus.WithField("space_index", index).Errorf("fail to get space, err: %s", err.Error())
		return dto, err
	}

	if space.IsPrivate() && space.Owner.Account() != index.Owner.Account() {
		return dto, commonrepo.NewNotAccessedError(xerrors.New("no permission to access the space"))
	}

	spaceId, err := domain.NewIdentity(space.RepoId)
	if err != nil {
		return dto, err
	}

	app, err := s.repo.FindBySpaceId(spaceId)
	if err == nil {
		return toSpaceAppDTO(&app), nil
	}

	if space.Exception.Exception() != "" {
		return toSpaceDTO(&space), nil
	}

	// FIXME:
	// if space.Hardware.IsNpu() && !space.CompPowerAllocated {
	// 	return toSpaceNoCompQuotaDTO(&space), nil
	// }

	if commonrepo.IsErrorResourceNotExists(err) {
		err = allerror.NewNotFound(allerror.ErrorCodeSpaceAppNotFound, "space app not found", err)
	} else {
		err = xerrors.Errorf("find space app by id failed, err: %w", err)
	}
	return dto, err
}

// GetBuildLogs
func (s *spaceappAppService) GetBuildLogs(ctx context.Context, user domain.Account, index *spacedomain.SpaceIndex) (
	dto BuildLogsDTO, err error,
) {
	app, err := s.getPrivateReadSpaceApp(user, index)
	if err != nil {
		err = xerrors.Errorf("failed to get space app, err:%w", err)
		return
	}

	spaceApp, err := s.repo.FindById(app.Id)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeSpaceAppNotFound, "space app not found", err)
		} else {
			err = xerrors.Errorf("find space app by id failed, err: %w", err)
		}

		return
	}

	dto.Logs = spaceApp.AppLogURL.URL()

	return
}

// GetBuildLog for get build log
func (s *spaceappAppService) GetBuildLog(
	ctx context.Context, user domain.Account, index *spacedomain.SpaceIndex,
) (string, error) {
	app, err := s.getPrivateReadSpaceApp(user, index)
	if err != nil {
		return "", xerrors.Errorf("failed to get space app:%w", err)
	}

	if app.BuildLogURL == nil {
		return "", xerrors.New("space app is not building")
	}
	if app.BuildLogURL.URL() == "" {
		return "", xerrors.New("space app is not building")
	}

	return app.BuildLogURL.URL(), nil
}

func (s *spaceappAppService) getPrivateReadSpaceApp(
	user domain.Account, index *spacedomain.SpaceIndex,
) (spaceappdomain.SpaceApp, error) {
	var spaceApp spaceappdomain.SpaceApp

	space, err := s.spaceRepo.GetByName(user, index.Name)
	if err != nil {
		return spaceApp, err
	}

	if space.Owner.Account() != index.Owner.Account() {
		return spaceApp, commonrepo.NewErrorResourceNotExists(xerrors.New("not found"))
	}

	spaceId, err := domain.NewIdentity(space.RepoId)
	if err != nil {
		return spaceApp, err
	}

	app, err := s.repo.FindBySpaceId(spaceId)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			err = allerror.NewNotFound(allerror.ErrorCodeSpaceAppNotFound, "space app not found", err)
		} else {
			err = xerrors.Errorf("find space app by id failed, err: %w", err)
		}
	}

	return app, err
}

// GetRequestDataStream
func (s *spaceappAppService) GetRequestDataStream(cmd *spaceappdomain.SeverSentStream) error {
	return s.sse.Request(cmd)
}

// GetSpaceLog for get serving log
func (s *spaceappAppService) GetSpaceLog(
	ctx context.Context, user domain.Account, index *spacedomain.SpaceIndex,
) (string, error) {
	app, err := s.getPrivateReadSpaceApp(user, index)
	if err != nil {
		return "", xerrors.Errorf("failed to get space app:%w", err)
	}

	if app.AppLogURL == nil {
		return "", xerrors.New("space app is not serving")
	}
	if app.AppLogURL.URL() == "" {
		return "", xerrors.New("space app is not serving")
	}

	return app.AppLogURL.URL(), nil
}

func toSpaceAppDTO(app *spaceappdomain.SpaceApp) SpaceAppDTO {
	dto := SpaceAppDTO{
		Id:     app.Id.Identity(),
		Status: app.Status.AppStatus(),
		Reason: app.GetFailedReason(),
	}

	if app.AppURL != nil {
		dto.AppURL = app.AppURL.AppURL()
	}

	if app.AppLogURL != nil {
		dto.AppLogURL = app.AppLogURL.URL()
	}

	if app.BuildLogURL != nil {
		dto.BuildLogURL = app.BuildLogURL.URL()
	}

	return dto
}

// CheckPermissionRead  check user permission for read space app.
func (s *spaceappAppService) CheckPermissionRead(
	ctx context.Context, user domain.Account, index *spacedomain.SpaceIndex) error {
	space, err := s.spaceRepo.GetByName(index.Owner, index.Name)
	if err != nil {
		return err
	}

	if space.IsPrivate() && user != index.Owner {
		return errors.New("space app read permission denied")
	}

	return nil
}
