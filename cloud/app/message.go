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
	survivalTimeForPodCPU int64,
	survivalTimeForPodAscend int64,
) CloudMessageService {
	return &cloudMessageService{
		repo:                     repo,
		manager:                  manager,
		survivalTimeForPodCPU:    survivalTimeForPodCPU,
		survivalTimeForPodAscend: survivalTimeForPodAscend,
	}
}

type cloudMessageService struct {
	repo                     repository.Pod
	manager                  cloud.CloudPod
	survivalTimeForPodCPU    int64
	survivalTimeForPodAscend int64
}

func (c *cloudMessageService) CreatePodInstance(p *domain.PodInfo) error {
	// create pod instance by SDK
	logrus.Infof("send create pod info: %#v", p)
	survivalTime := p.Expiry.PodExpiry()
	if p.IsCpu() {
		survivalTime = c.survivalTimeForPodCPU
	} else if p.IsAscend() {
		survivalTime = c.survivalTimeForPodAscend
	}

	expire, err := domain.NewPodExpiry(survivalTime)
	if err != nil {
		return err
	}

	err = c.manager.Create(
		&cloud.CloudPodCreateInfo{
			PodId:        p.Id,
			SurvivalTime: survivalTime,
			User:         p.Owner.Account(),
			CloudType:    p.GetCloudType(),
		},
	)

	if err != nil {
		return err
	}

	// update pod status in DB
	p.StatusSetCreating()
	p.Expiry = expire

	return c.repo.UpdatePod(p)
}
