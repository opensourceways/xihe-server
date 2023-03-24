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

type PodInfoCmd SubscribeCloudCmd

type GetCloudConfCmd struct {
	IsVisitor bool
	User      types.Account
}

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
	Id        string `json:"id"`
	Spec      string `json:"spec"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	Feature   string `json:"feature"`
	Processor string `json:"processor"`
	Credit    int64  `json:"credit"`
}

type CloudDTO struct {
	CloudConfDTO

	IsIdle     bool `json:"is_idle"`
	HasHolding bool `json:"has_holding"`
}

type PodInfoDTO struct {
	Id        string `json:"id"`
	CloudId   string `json:"cloud_id"`
	Owner     string `json:"owner"`
	Status    string `json:"status"`
	Expiry    int64  `json:"expiry"`
	Error     string `json:"error"`
	AccessURL string `json:"access_url"`
	CreatedAt int64  `json:"created_at"`
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
	b := cmd.User.Account() != "" &&
		cmd.CloudId != ""

	if !b {
		return errors.New("invalid cmd")
	}

	return nil
}

func (cmd *GetCloudConfCmd) ToCmd(user types.Account, visitor bool) {
	*cmd = GetCloudConfCmd{
		IsVisitor: visitor,
		User:      user,
	}
}

func (r *UpdatePodInternalCmd) toPodInfo(p *domain.PodInfo) (err error) {
	p.Id = r.PodId
	p.Error = r.PodError
	p.AccessURL = r.AccessURL

	return
}

func (r *CloudConfDTO) toCloudConfDTO(c *domain.CloudConf) {
	*r = CloudConfDTO{
		Id:        c.Id,
		Spec:      c.Spec.CloudSpec(),
		Name:      c.Name.CloudName(),
		Image:     c.Image.CloudImage(),
		Feature:   c.Feature.CloudFeature(),
		Processor: c.Processor.CloudProcessor(),
		Credit:    c.Credit.Credit(),
	}
}

func (r *CloudDTO) toCloudDTO(c *domain.Cloud, isIdle bool, hasHolding bool) {
	r.CloudConfDTO.toCloudConfDTO(&c.CloudConf)

	r.IsIdle = isIdle
	r.HasHolding = hasHolding
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
		r.Expiry = p.Expiry.PodExpiry()
	}

	if p.Error != nil {
		r.Error = p.Error.PodError()
	}

	if p.AccessURL != nil {
		r.AccessURL = p.AccessURL.AccessURL()
	}

	if p.CreatedAt != nil {
		r.CreatedAt = p.CreatedAt.Time()
	}
}
