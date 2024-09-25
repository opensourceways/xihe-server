package controller

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/promotion/app"
	promotiond "github.com/opensourceways/xihe-server/promotion/domain"
	userctl "github.com/opensourceways/xihe-server/user/controller"
	"github.com/opensourceways/xihe-server/user/domain"
)

type PromotionApplyReq struct {
	userctl.ApplyRequest
	Origin string `json:"origin"`
}

func (req PromotionApplyReq) ToCmd(promotionid string, user types.Account) (cmd app.UserRegistrationCmd, err error) {
	userRegistrationCmd, err := req.ApplyRequest.ToCmd(user)
	if err != nil {
		return
	}

	origin, err := promotiond.NewOrigin(req.Origin)
	if err != nil {
		return
	}

	return app.UserRegistrationCmd{
		PromotionId:      promotionid,
		UserRegistration: domain.UserRegInfo(userRegistrationCmd),
		Origin:           origin,
	}, nil
}
