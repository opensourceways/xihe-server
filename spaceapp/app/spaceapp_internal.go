package app

import (
	"context"

	spacerepo "github.com/opensourceways/xihe-server/space/domain/repository"
	"github.com/opensourceways/xihe-server/spaceapp/domain"
	"github.com/opensourceways/xihe-server/spaceapp/domain/repository"
	"github.com/sirupsen/logrus"
)

// SpaceappInternalAppService is an interface that defines the methods for creating and managing a SpaceApp.
type SpaceappInternalAppService interface {
	NotifyIsServing(ctx context.Context, cmd *CmdToNotifyServiceIsStarted) error
}

// NewSpaceappInternalAppService creates a new instance of spaceappInternalAppService
// with the provided message and repository.
func NewSpaceappInternalAppService(
	repo repository.SpaceAppRepository,
) *spaceappInternalAppService {
	return &spaceappInternalAppService{
		repo: repo,
	}
}

// spaceappInternalAppService
type spaceappInternalAppService struct {
	spaceRepo spacerepo.Project
	repo      repository.SpaceAppRepository
}

// NotifyIsServing notifies that a service of a SpaceApp has serving.
func (s *spaceappInternalAppService) NotifyIsServing(ctx context.Context, cmd *CmdToNotifyServiceIsStarted) error {
	v, err := s.getSpaceApp(ctx, cmd.SpaceAppIndex)
	if err != nil {
		return err
	}

	if err := v.StartServing(cmd.AppURL, cmd.LogURL); err != nil {
		logrus.Errorf("spaceId:%s set space app serving failed, err:%s", cmd.SpaceId.Identity(), err)
		return err
	}

	if err := s.repo.SaveWithoutAllBuildLog(&v); err != nil {
		logrus.Errorf("spaceId:%s save db failed", cmd.SpaceId.Identity())
		return err
	}
	logrus.Infof("spaceId:%s notify serving successful", cmd.SpaceId.Identity())

	return nil
}

func (s *spaceappInternalAppService) getSpaceApp(ctx context.Context, cmd CmdToCreateApp) (domain.SpaceApp, error) {
	// space, err := s.spaceRepo.Get(nil, cmd.SpaceId.Identity())
	// if err != nil {
	// 	if commonrepo.IsErrorResourceNotExists(err) {
	// 		err = newSpaceNotFound(xerrors.Errorf("space not found, err:%w", err))
	// 	} else {
	// 		err = xerrors.Errorf("failed to get space, err:%w", err)
	// 	}
	// 	logrus.Errorf("spaceId:%s get space failed, err:%s", cmd.SpaceId.Identity(), err)
	// 	return domain.SpaceApp{}, err
	// }

	// if space.CommitId != cmd.CommitId {
	// 	err = allerror.New(allerror.ErrorCodeSpaceCommitConflict, "commit conflict",
	// 		xerrors.Errorf("spaceId:%s commit conflict", space.Id.Identity()))
	// 	logrus.Errorf("spaceId:%s latest commitId:%s, old commitId:%s, err:%s",
	// 		cmd.SpaceId.Identity(), space.CommitId, cmd.CommitId, err)
	// 	return domain.SpaceApp{}, err
	// }

	// v, err := s.repo.FindBySpaceId(ctx, space.Id)
	// if err != nil {
	// 	if commonrepo.IsErrorResourceNotExists(err) {
	// 		err = newSpaceAppNotFound(err)
	// 	}
	// 	logrus.Errorf("spaceId:%s get space app failed, err:%s", space.Id.Identity(), err)
	// 	return domain.SpaceApp{}, err
	// }
	// return v, nil
	return domain.SpaceApp{}, nil
}
