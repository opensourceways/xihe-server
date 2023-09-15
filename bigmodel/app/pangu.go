package app

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

func (s bigModelService) PanGu(u types.Account, q string) (v string, code string, err error) {
	_ = s.sender.SendBigModelAccessLog(&domain.BigModelAccessLogEvent{
		Account:      u,
		BigModelType: domain.BigmodelPanGu,
	})

	if v, err = s.fm.PanGu(q); err != nil {
		code = s.setCode(err)
	}

	return
}
