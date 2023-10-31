package repositoryimpl

import (
	"github.com/opensourceways/xihe-server/cloud/domain"
	types "github.com/opensourceways/xihe-server/common/domain"
	otypes "github.com/opensourceways/xihe-server/domain"
)

const (
	fieldId      = "id"
	fieldCloudId = "cloud_id"
	fieldStatus  = "status"
	fieldOwner   = "owner"
)

func (doc *DCloudConf) toCloudConf(c *domain.CloudConf) (err error) {
	c.Id = doc.Id

	if c.Name, err = domain.NewCloudName(doc.Name); err != nil {
		return
	}

	if c.Spec, err = domain.NewCloudSpec(doc.Spec); err != nil {
		return
	}

	if c.Image, err = domain.NewCloudImage(doc.Image); err != nil {
		return
	}

	if c.Feature, err = domain.NewCloudFeature(doc.Feature); err != nil {
		return
	}

	if c.Processor, err = domain.NewCloudProcessor(doc.Processor); err != nil {
		return
	}

	if c.Limited, err = domain.NewCloudLimited(doc.Limited); err != nil {
		return
	}

	if c.Credit, err = domain.NewCredit(doc.Credit); err != nil {
		return
	}

	return
}

func (table *TPod) toPodInfo(p *domain.PodInfo) (err error) {
	p.Id = table.Id
	p.CloudId = table.CloudId

	if p.Owner, err = otypes.NewAccount(table.Owner); err != nil {
		return
	}

	if p.Status, err = domain.NewPodStatus(table.Status); err != nil {
		return
	}

	if p.Expiry, err = domain.NewPodExpiry(table.Expiry); err != nil {
		return
	}

	if p.Error, err = domain.NewPodError(table.Error); err != nil {
		return
	}

	if p.AccessURL, err = domain.NewAccessURL(table.AccessURL); err != nil {
		return
	}

	if p.CreatedAt, err = types.NewTime(table.CreatedAt); err != nil {
		return
	}

	return
}

func (table *TPod) toTPod(p *domain.PodInfo) {
	*table = TPod{
		CloudId: p.CloudId,
		Status:  p.Status.PodStatus(),
	}

	if p.Id != "" {
		table.Id = p.Id
	}

	if p.Owner != nil {
		table.Owner = p.Owner.Account()
	}

	if p.Expiry != nil {
		table.Expiry = p.Expiry.PodExpiry()
	}

	if p.Error != nil {
		table.Error = p.Error.PodError()
	}

	if p.AccessURL != nil {
		table.AccessURL = p.AccessURL.AccessURL()
	}

	if p.CreatedAt != nil {
		table.CreatedAt = p.CreatedAt.Time()
	}
}
