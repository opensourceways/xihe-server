package controller

import (
	cloudapp "github.com/opensourceways/xihe-server/cloud/app"
	cloudtypes "github.com/opensourceways/xihe-server/cloud/domain"
	"github.com/opensourceways/xihe-server/domain"
)

type cloudSubscribeRequest struct {
	CloudId  string `json:"cloud_id"`
	Image    string `json:"image" binding:"required"`
	CardsNum int    `json:"cards_num" binding:"required,min=1"`
}

func (req *cloudSubscribeRequest) toCmd(user domain.Account) cloudapp.SubscribeCloudCmd {
	cmd := cloudapp.SubscribeCloudCmd{
		User:    user,
		CloudId: req.CloudId,
	}

	cmd.ImageAlias, _ = cloudtypes.NewCloudImageAlias(req.Image)
	cmd.CardsNum, _ = cloudtypes.NewCloudSpecCardsNum(req.CardsNum)

	return cmd
}
