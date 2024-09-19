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

	c.Specs = make([]domain.CloudSpec, 0, len(doc.Specs))
	for i := range doc.Specs {
		cloudSpec := domain.CloudSpec{}

		if cloudSpec.Desc, err = domain.NewCloudSpecDesc(doc.Specs[i].Desc); err != nil {
			return
		}

		if cloudSpec.CardsNum, err = domain.NewCloudSpecCardsNum(doc.Specs[i].CardsNum); err != nil {
			return
		}

		c.Specs = append(c.Specs, cloudSpec)
	}

	c.Images = make([]domain.CloudImage, 0, len(doc.Images))
	for i := range doc.Images {
		cloudImage := domain.CloudImage{}

		if cloudImage.Alias, err = domain.NewCloudImageAlias(doc.Images[i].Alias); err != nil {
			return
		}

		if cloudImage.Image, err = domain.NewICloudImage(doc.Images[i].Image); err != nil {
			return
		}

		c.Images = append(c.Images, cloudImage)
	}

	if c.Feature, err = domain.NewCloudFeature(doc.Feature); err != nil {
		return
	}

	if c.Processor, err = domain.NewCloudProcessor(doc.Processor); err != nil {
		return
	}

	if c.SingleLimited, err = domain.NewCloudLimited(doc.SingleLimited); err != nil {
		return
	}
	if c.MultiLimited, err = domain.NewCloudLimited(doc.MultiLimited); err != nil {
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
	p.Image = table.Image

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

	if p.CardsNum, err = domain.NewCloudSpecCardsNum(table.CardsNum); err != nil {
		return
	}

	return
}

func (table *TPod) toTPod(p *domain.PodInfo) {
	*table = TPod{
		CloudId: p.CloudId,
		Image:   p.Image,
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

	if p.Status != nil {
		table.Status = p.Status.PodStatus()
	}

	if p.CardsNum != nil {
		table.CardsNum = p.CardsNum.CloudSpecCardsNum()
	}
}
