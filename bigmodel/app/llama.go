package app

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain"
)

func (s bigModelService) LLAMA2(cmd *LLAMA2Cmd) (code string, err error) {
	_ = s.sender.SendBigModelStarted(&domain.BigModelStartedEvent{
		Account:      cmd.User,
		BigModelType: domain.BigmodelLLAMA2,
	})

	input := &domain.LLAMA2Input{
		Text:              cmd.Text,
		Sampling:          cmd.Sampling,
		History:           cmd.History,
		TopK:              cmd.TopK,
		TopP:              cmd.TopP,
		Temperature:       cmd.Temperature,
		RepetitionPenalty: cmd.RepetitionPenalty,
	}

	if err = s.fm.LLAMA2(cmd.CH, input); err != nil {
		code = s.setCode(err)
		
		return
	}

	_ = s.sender.SendBigModelFinished(&domain.BigModelFinishedEvent{
		Account:      cmd.User,
		BigModelType: domain.BigmodelLLAMA2,
	})

	return
}
