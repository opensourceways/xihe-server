package app

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain"
)

func (s bigModelService) GLM2(cmd *GLM2Cmd) (code string, err error) {
	_ = s.sender.SendBigModelStarted(&domain.BigModelStartedEvent{
		Account:      cmd.User,
		BigModelType: domain.BigmodelGLM2,
	})

	input := &domain.GLM2Input{
		Text:              cmd.Text,
		Sampling:          cmd.Sampling,
		History:           cmd.History,
		TopK:              cmd.TopK,
		TopP:              cmd.TopP,
		Temperature:       cmd.Temperature,
		RepetitionPenalty: cmd.RepetitionPenalty,
	}

	if err = s.fm.GLM2(cmd.CH, input); err != nil {
		code = s.setCode(err)
		
		return
	}

	_ = s.sender.SendBigModelFinished(&domain.BigModelFinishedEvent{
		Account:      cmd.User,
		BigModelType: domain.BigmodelGLM2,
	})

	return
}
