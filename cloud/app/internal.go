package app

import (
	"github.com/opensourceways/xihe-server/cloud/domain"
	"github.com/opensourceways/xihe-server/cloud/domain/repository"
)

type CloudInternalService interface {
	UpdateInfo(*UpdatePodInternalCmd) error
}

func NewCloudInternalService(
	repo repository.Pod,
) CloudInternalService {
	return &cloudInternalService{
		repo: repo,
	}
}

type cloudInternalService struct {
	repo repository.Pod
}

func (s *cloudInternalService) UpdateInfo(cmd *UpdatePodInternalCmd) error {
	p := new(domain.PodInfo)
	cmd.toPodInfo(p)

	p.SetStatus()

	err := s.repo.UpdatePod(p)
	return err
}
