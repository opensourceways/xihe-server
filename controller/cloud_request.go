package controller

import (
	cloudapp "github.com/opensourceways/xihe-server/cloud/app"
	"github.com/opensourceways/xihe-server/domain"
	userapp "github.com/opensourceways/xihe-server/user/app"
	userdomain "github.com/opensourceways/xihe-server/user/domain"
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

func toWhiteListCmd(user domain.Account) (cmd userapp.UserWhiteListCmd, err error) {
	cmd.Account = user

	v, err := userdomain.NewWhiteListType("cloud")
	if err != nil {
		return
	}
	cmd.Type = v

	return
}
