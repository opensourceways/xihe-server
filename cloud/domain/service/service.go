package service

import (
	"github.com/opensourceways/xihe-server/cloud/domain"
	"github.com/opensourceways/xihe-server/cloud/domain/message"
	"github.com/opensourceways/xihe-server/cloud/domain/repository"
	commonrepo "github.com/opensourceways/xihe-server/common/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
)

type CloudService struct {
	podRepo repository.Pod
	sender  message.CloudMessageProducer
}

func NewCloudService(
	pod repository.Pod,
	sender message.CloudMessageProducer,
) CloudService {
	return CloudService{
		pod,
		sender,
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

func (r *CloudService) SubscribeCloud(
	c *domain.CloudConf, u types.Account, imageAlias string,
) (err error) {
	var image domain.ICloudImage
	if image, err = c.GetImage(imageAlias); err != nil {
		return
	}

	// save into repo
	p := new(domain.PodInfo)
	if err := p.SetStartingPodInfo(c.Id, u); err != nil {
		return err
	}

	var pid string
	if pid, err = r.podRepo.AddStartingPod(p); err != nil {
		return
	}

	// send msg to call pod instance api
	msg := new(message.MsgCloudConf)
	msg.ToMsgCloudConf(c, u, pid, image)

	return r.sender.SubscribeCloud(msg)
}

func (r *CloudService) CheckUserCanSubsribe(user types.Account, cid string) (
	p domain.PodInfo, ok bool, err error,
) {
	p, err = r.podRepo.GetUserCloudIdLastPod(user, cid)
	if err != nil {
		if commonrepo.IsErrorResourceNotExists(err) {
			return p, true, nil
		}

		return
	}

	if p.IsExpiried() || p.IsFailedOrTerminated() {
		return p, true, err
	}

	return p, false, err
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
