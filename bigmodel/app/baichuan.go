package app

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain"
)

func (s bigModelService) BaiChuan(cmd *BaiChuanCmd) (code string, dto BaiChuanDTO, err error) {
	_ = s.sender.SendBigModelStarted(&domain.BigModelStartedEvent{
		Account:      cmd.User,
		BigModelType: domain.BigmodelBaiChuan,
	})

	input := &domain.BaiChuanInput{
		Text:              cmd.Text,
		TopK:              cmd.TopK,
		TopP:              cmd.TopP,
		Temperature:       cmd.Temperature,
		RepetitionPenalty: cmd.RepetitionPenalty,
	}

	if code, dto.Text, err = s.fm.BaiChuan(input); err != nil {
		return
	}

	return
}
