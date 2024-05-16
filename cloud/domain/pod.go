package domain

import (
	types "github.com/opensourceways/xihe-server/common/domain"
	otypes "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

const (
	cloudIdCPU   = "cpu_001"
	cloudIdNPU   = "ascend_001"
	cloudTypeCPU = "cpu"
	cloudTypeNPU = "npu"
)

type Pod struct {
	Id      string
	CloudId string
	Owner   otypes.Account
}

type PodInfo struct {
	Pod

	Status    PodStatus
	Expiry    PodExpiry
	Error     PodError
	AccessURL AccessURL
	CreatedAt types.Time
}

func (r *Pod) IsOnwer(owner otypes.Account) bool {
	return r.Owner == owner
}

func (p *PodInfo) CanRelease() bool {
	return p.Status.IsRunning()
}

func (p *PodInfo) IsExpiried() bool {
	return utils.Now() > p.Expiry.PodExpiry()
}

func (p *PodInfo) IsFailedOrTerminated() bool {
	return p.Status.IsFailed() || p.Status.IsTerminated()
}

func (p *PodInfo) IsHoldingAndNotExpiried() bool {
	if p.IsExpiried() {
		return false
	}

	return p.Status.IsCreating() || p.Status.IsStarting() || p.Status.IsRunning()
}

func (p *PodInfo) CheckGoodAndSet() bool {
	if !p.Error.IsGood() {
		p.Status, _ = NewPodStatus(cloudPodStatusFailed)
		return false
	}

	return true
}

func (p *PodInfo) StatusSetCreating() {
	p.Status, _ = NewPodStatus(cloudPodStatusCreating)
}

func (p *PodInfo) StatusSetRunning() {
	p.Status, _ = NewPodStatus(cloudPodStatusRunning)
}

func (p *PodInfo) StatusSetFailed() {
	p.Status, _ = NewPodStatus(cloudPodStatusFailed)
}

func (p *PodInfo) SetStatus() {
	if p.AccessURL.AccessURL() != "" {
		p.StatusSetRunning()
	} else {
		p.StatusSetFailed()
	}
}

func (p *PodInfo) SetDefaultExpiry() (err error) {
	if p.Expiry, err = NewPodExpiry(utils.Now() + 2*60*60); err != nil { // TODO conifg
		return
	}

	return
}

func (p *PodInfo) SetStartingPodInfo(cid string, owner otypes.Account) (err error) {
	p.CloudId = cid
	p.Owner = owner

	if p.Status, err = NewPodStatus(cloudPodStatusStarting); err != nil {
		return
	}

	return
}

func (p *PodInfo) GetCloudType() string {
	if p.CloudId == cloudIdCPU {
		return cloudTypeCPU
	}

	return cloudTypeNPU
}
