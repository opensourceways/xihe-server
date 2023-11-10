package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	"github.com/opensourceways/xihe-server/bigmodel/domain/bigmodel"
	"github.com/opensourceways/xihe-server/bigmodel/domain/message"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/sirupsen/logrus"
)

type AsyncBigModelService interface {
	WuKong(uint64, types.Account, *WuKongCmd) error
	GetIdleEndpoint(bid string) (int, error)
}

func NewAsyncBigModelService(
	fm bigmodel.BigModel,
	sender message.MessageProducer,
) AsyncBigModelService {
	return &asyncBigModelService{
		fm:     fm,
		sender: sender,
	}
}

type asyncBigModelService struct {
	fm     bigmodel.BigModel
	sender message.MessageProducer
}

func (s *asyncBigModelService) WuKong(tid uint64, user types.Account, cmd *WuKongCmd) (err error) {
	// 1. inference
	sendErr := s.sender.SendBigModelStarted(&domain.BigModelStartedEvent{
		Account:      user,
		BigModelType: domain.BigmodelWuKong,
	})
	if sendErr != nil {
		logrus.Warnf("send bigmodel started failed, err: %s", sendErr.Error())
	}

	if err = s.sender.SendWuKongAsyncTaskStart(&domain.WuKongAsyncTaskStartEvent{
		Account: user,
		TaskId:  tid,
	}); err != nil {
		return
	}

	links, err := s.fm.GenPicturesByWuKong(user, &cmd.WuKongPictureMeta, cmd.EsType)
	if err != nil {
		if !bigmodel.IsErrorSensitiveInfo(err) {
			err = errors.New("internal error")
		}

		if err = s.sender.SendWuKongInferenceError(&domain.WuKongInferenceErrorEvent{
			Account: user,
			TaskId:  tid,
			ErrMsg:  err.Error(),
		}); err != nil {
			return
		}

		return
	}

	// 3. send msg
	return s.sender.SendWuKongAsyncInferenceFinish(&domain.WuKongAsyncInferenceFinishEvent{
		Account: user,
		TaskId:  tid,
		Links:   links,
	})
}

func (s *asyncBigModelService) GetIdleEndpoint(bid string) (c int, err error) {
	return s.fm.GetIdleEndpoint(bid)
}
