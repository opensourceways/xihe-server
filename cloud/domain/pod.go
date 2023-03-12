package domain

import (
	types "github.com/opensourceways/xihe-server/common/domain"
	otypes "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
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
