package controller

import (
	cloudapp "github.com/opensourceways/xihe-server/cloud/app"
	cloudtypes "github.com/opensourceways/xihe-server/cloud/domain"
	"github.com/opensourceways/xihe-server/domain"
)

type cloudSubscribeRequest struct {
	CloudId string `json:"cloud_id"`
	Image   string `json:"image" binding:"required"`
}

func (req *cloudSubscribeRequest) toCmd(user domain.Account) cloudapp.SubscribeCloudCmd {
	cmd := cloudapp.SubscribeCloudCmd{
		User:    user,
		CloudId: req.CloudId,
	}

	cmd.ImageAlias, _ = cloudtypes.NewCloudImageAlias(req.Image)

	return cmd
}
