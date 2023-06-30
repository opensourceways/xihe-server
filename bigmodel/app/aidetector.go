package app

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain"
	"github.com/opensourceways/xihe-server/bigmodel/infrastructure/bigmodels"
)

func (s bigModelService) AIDetector(cmd *AIDetectorCmd) (code string, ismachine bool, err error) {
	// operation log
	_ = s.sender.AddOperateLogForAccessBigModel(cmd.User, domain.BigmodelAIDetector)

	// audit
	if err = s.fm.CheckText(cmd.Text.AIDetectorText()); err != nil {
		code = ErrorBigModelSensitiveInfo

		return
	}

	// detector
	if ismachine, err = s.fm.AIDetector(domain.AIDetectorInput{
		Lang: cmd.Lang,
		Text: cmd.Text,
	}); err != nil {
		if bigmodels.IsErrorConcurrentRequest(err) {
			code = ErrorBigModelConcurrentRequest

			return
		}

		return
	}

	return
}
