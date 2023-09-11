package app

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain"
)

func (s bigModelService) BaiChuan(cmd *BaiChuanCmd) (code string, dto BaiChuanDTO, err error) {
	_ = s.sender.AddOperateLogForAccessBigModel(cmd.User, domain.BigmodelBaiChuan)

	_ = s.sender.ExperienceBigmodel(cmd.User, domain.BigmodelBaiChuan)

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
