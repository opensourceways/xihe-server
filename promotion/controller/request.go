package controller

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/promotion/app"
	userctl "github.com/opensourceways/xihe-server/user/controller"
	"github.com/opensourceways/xihe-server/user/domain"
)

type PromotionApplyReq struct {
	userctl.ApplyRequest
}

func (req PromotionApplyReq) ToCmd(promotionid string, user types.Account) (cmd app.UserRegistrationCmd, err error) {
	userRegisrationCmd, err := req.ApplyRequest.ToCmd(user)
	if err != nil {
		return
	}

	return app.UserRegistrationCmd{
		PromotionId:      promotionid,
		UserRegistration: domain.UserRegInfo(userRegisrationCmd),
	}, nil
}
