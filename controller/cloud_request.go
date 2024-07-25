package controller

import (
	cloudapp "github.com/opensourceways/xihe-server/cloud/app"
	"github.com/opensourceways/xihe-server/domain"
)

type cloudSubscribeRequest struct {
	CloudId string `json:"cloud_id"`
}

func (req *cloudSubscribeRequest) toCmd(user domain.Account) cloudapp.SubscribeCloudCmd {
	return cloudapp.SubscribeCloudCmd{
		User:    user,
		CloudId: req.CloudId,
	}
}
