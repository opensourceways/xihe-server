package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/cloud/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type SubscribeCloudCmd struct {
	User    types.Account
	CloudId string
}

type PodInfoCmd domain.PodInfo

type RelasePodCmd struct {
	User  types.Account
	PodId string
}

type UpdatePodInternalCmd struct {
	PodId     string
	PodError  domain.PodError
	AccessURL domain.AccessURL
}

type CloudConfDTO struct {
	Spec      string `json:"spec"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	Feature   string `json:"feature"`
	Processor string `json:"processor"`
}

type CloudDTO struct {
	CloudConfDTO

	IsFree bool `json:"is_free"`
}

type PodInfoDTO struct {
	Id         string `json:"id"`
	CloudId    string `json:"cloud_id"`
	Owner      string `json:"owner"`
	Status     string `json:"status"`
	ExpiryDate string `json:"expiry_date"`
	Error      string `json:"error"`
	AccessURL  string `json:"access_url"`
	CreatedAt  string `json:"created_at"`
}

func (cmd *SubscribeCloudCmd) Validate() error {
	b := cmd.User.Account() != "" &&
		cmd.CloudId != ""

	if !b {
		return errors.New("invalid cmd")
	}

	return nil
}

func (cmd *PodInfoCmd) Validate() error {
	b := cmd.Owner.Account() != "" &&
		cmd.CloudId != ""

	if !b {
		return errors.New("invalid cmd")
	}

	return nil
}

func (r *UpdatePodInternalCmd) toPodInfo(p *domain.PodInfo) (err error) {
	p.Id = r.PodId
	p.Error = r.PodError
	p.AccessURL = r.AccessURL

	return
}

func (r *CloudConfDTO) toCloudConfDTO(c *domain.CloudConf) {
	*r = CloudConfDTO{
		Spec:      c.Spec.CloudSpec(),
		Name:      c.Name.CloudName(),
		Image:     c.Image.CloudImage(),
		Feature:   c.Feature.CloudFeature(),
		Processor: c.Processor.CloudProcessor(),
	}
}

func (r *CloudDTO) toCloudDTO(c *domain.Cloud, isFree bool) {
	r.CloudConfDTO.toCloudConfDTO(&c.CloudConf)

	r.IsFree = isFree
}

func (r *PodInfoDTO) toPodInfoDTO(p *domain.PodInfo) {
	r.Id = p.Id
	r.CloudId = p.CloudId

	if p.Owner != nil {
		r.Owner = p.Owner.Account()
	}

	if p.Status != nil {
		r.Status = p.Status.PodStatus()
	}

	if p.Expiry != nil {
		r.ExpiryDate = p.Expiry.PodExpiryDate()
	}

	if p.Error != nil {
		r.Error = p.Error.PodError()
	}

	if p.AccessURL != nil {
		r.AccessURL = p.AccessURL.AccessURL()
	}

	if p.CreatedAt != nil {
		r.CreatedAt = p.CreatedAt.TimeDate()
	}
}
