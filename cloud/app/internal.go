package app

import (
	"time"

	"github.com/opensourceways/xihe-server/cloud/domain"
	"github.com/opensourceways/xihe-server/cloud/domain/repository"
)

type CloudInternalService interface {
	UpdateInfo(*UpdatePodInternalCmd) error
	Release(*ReleaseInternalCmd) error
}

func NewCloudInternalService(
	repo repository.Pod,
	terminationWait int64,
) CloudInternalService {
	return &cloudInternalService{
		repo:            repo,
		terminationWait: terminationWait,
	}
}

type cloudInternalService struct {
	repo            repository.Pod
	terminationWait int64
}

func (s *cloudInternalService) UpdateInfo(cmd *UpdatePodInternalCmd) error {
	p := new(domain.PodInfo)
	if err := cmd.toPodInfo(p); err != nil {
		return err
	}

	p.SetStatus()

	err := s.repo.UpdatePod(p)

	return err
}

func (s *cloudInternalService) Release(cmd *ReleaseInternalCmd) error {
	p := new(domain.PodInfo)
	p.Id = cmd.PodId

	// wait a moment because pod doesn't terminate immediately
	expiry, err := domain.NewPodExpiry(time.Now().Add(time.Duration(s.terminationWait) * time.Second).Unix())
	if err != nil {
		return err
	}

	p.Expiry = expiry
	p.StatusSetTerminated()

	return s.repo.UpdatePod(p)
}
