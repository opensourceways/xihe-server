package app

import (
	"github.com/opensourceways/xihe-server/cloud/domain"
	"github.com/opensourceways/xihe-server/cloud/domain/cloud"
	"github.com/opensourceways/xihe-server/cloud/domain/repository"
	"github.com/sirupsen/logrus"
)

type CloudMessageService interface {
	CreatePodInstance(*domain.PodInfo) error
}

func NewCloudMessageService(
	repo repository.Pod,
	manager cloud.CloudPod,
	survivalTimeForPod int64,
) CloudMessageService {
	return &cloudMessageService{
		repo:               repo,
		manager:            manager,
		survivalTimeForPod: survivalTimeForPod,
	}
}

type cloudMessageService struct {
	repo               repository.Pod
	manager            cloud.CloudPod
	survivalTimeForPod int64
}

func (c *cloudMessageService) CreatePodInstance(p *domain.PodInfo) error {
	// create pod instance by SDK
	logrus.Debugf("send create pod info: %#v", p)
	err := c.manager.Create(
		&cloud.CloudPodCreateInfo{
			PodId:        p.Id,
			SurvivalTime: c.survivalTimeForPod,
			User:         p.Owner.Account(),
			CloudType:    p.GetCloudType(),
		},
	)

	if err != nil {
		return err
	}

	// update pod status in DB
	p.StatusSetCreating()

	return c.repo.UpdatePod(p)
}
