package service

import (
	"github.com/opensourceways/xihe-server/cloud/domain"
	"github.com/opensourceways/xihe-server/cloud/domain/repository"
	commonrepo "github.com/opensourceways/xihe-server/common/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
)

type CloudService struct {
	podRepo repository.Pod
}

func NewCloudService(
	pod repository.Pod,
) CloudService {
	return CloudService{
		pod,
	}
}

func (r *CloudService) caculateRemain(
	c *domain.Cloud, p *repository.PodInfoList,
) (err error) {
	// caculate running and not expiry pod
	var count int
	for i := range p.PodInfos {
		if !p.PodInfos[i].IsExpiried() {
			count++
		}
	}
	remain := c.CloudConf.Limited.CloudLimited() - count
	if remain < 0 {
		remain = 0
	}

	if c.Remain, err = domain.NewCloudRemain(remain); err != nil {
		return
	}

	return
}

func (r *CloudService) ToCloud(c *domain.Cloud) (err error) {
	plist, err := r.podRepo.GetRunningPod(c.CloudConf.Id)
	if err != nil {
		return
	}

	return r.caculateRemain(c, &plist)
}

func (r *CloudService) HasHolding(user types.Account, c *domain.CloudConf) (bool, error) {
	p, err := r.podRepo.GetUserCloudIdLastPod(user, c.Id)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			return false, err
		}

		return false, err
	}

	if p.IsHoldingAndNotExpiried() {
		return true, nil
	}

	return false, nil
}
