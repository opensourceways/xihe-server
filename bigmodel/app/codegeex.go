package app

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain"
	"github.com/opensourceways/xihe-server/bigmodel/domain/bigmodel"
	types "github.com/opensourceways/xihe-server/domain"
)

func (s bigModelService) CodeGeex(user types.Account, cmd *CodeGeexCmd) (
	dto CodeGeexDTO, code string, err error,
) {
	_ = s.sender.SendBigModelStarted(&domain.BigModelStartedEvent{
		Account:      user,
		BigModelType: domain.BigmodelCodeGeex,
	})

	if dto, err = s.fm.CodeGeex((*bigmodel.CodeGeexReq)(cmd)); err != nil {
		code = s.setCode(err)

		return
	}

	_ = s.sender.SendBigModelFinished(&domain.BigModelFinishedEvent{
		Account:      user,
		BigModelType: domain.BigmodelCodeGeex,
	})

	return
}
