package app

import (
	"github.com/opensourceways/xihe-server/cloud/domain"
	"github.com/opensourceways/xihe-server/cloud/domain/cloud"
	"github.com/opensourceways/xihe-server/cloud/domain/message"
	"github.com/opensourceways/xihe-server/cloud/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
	"github.com/sirupsen/logrus"
)

type CloudMessageService interface {
	CreatePodInstance(*domain.PodInfo) error
	ReleasePodInstance(Id, cloudType string) error
}

func NewCloudMessageService(
	repo repository.Pod,
	manager cloud.CloudPod,
	survivalTimeForPodCPU int64,
	survivalTimeForPodAscend int64,
	cloudRecordEventPublisher message.CloudRecordEventPublisher,
) CloudMessageService {
	return &cloudMessageService{
		repo:                      repo,
		manager:                   manager,
		survivalTimeForPodCPU:     survivalTimeForPodCPU,
		survivalTimeForPodAscend:  survivalTimeForPodAscend,
		cloudRecordEventPublisher: cloudRecordEventPublisher,
	}
}

type cloudMessageService struct {
	repo                      repository.Pod
	manager                   cloud.CloudPod
	survivalTimeForPodCPU     int64
	survivalTimeForPodAscend  int64
	cloudRecordEventPublisher message.CloudRecordEventPublisher
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

	expire, err := domain.NewPodExpiry(utils.Now() + survivalTime)
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

	if err = c.repo.UpdatePod(p); err != nil {
		return err
	}

	return c.cloudRecordEventPublisher.Publish(&message.CloudRecordEvent{
		Owner:  p.Owner,
		ClouId: p.CloudId,
	})
}

func (c *cloudMessageService) ReleasePodInstance(Id, cloudType string) error {
	logrus.Infof("release pod id: %s, type: %s", Id, cloudType)

	return c.manager.Release(Id, cloudType)
}
