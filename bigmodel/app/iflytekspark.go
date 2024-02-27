package app

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain"
)

func (s bigModelService) IFlytekSpark(cmd *IFlytekSparkCmd) (code string, err error) {
	_ = s.sender.SendBigModelStarted(&domain.BigModelStartedEvent{
		Account:      cmd.User,
		BigModelType: domain.BigmodelIFlytekSpark,
	})

	input := &domain.IFlytekSparkInput{
		Text:              cmd.Text,
		Sampling:          cmd.Sampling,
		TopK:              cmd.TopK,
		Temperature:       cmd.Temperature,
		RepetitionPenalty: cmd.RepetitionPenalty,
	}

	if err = s.fm.IFlytekSpark(cmd.CH, input); err != nil {
		code = s.setCode(err)

		return
	}

	_ = s.sender.SendBigModelFinished(&domain.BigModelFinishedEvent{
		Account:      cmd.User,
		BigModelType: domain.BigmodelIFlytekSpark,
	})

	return
}
