package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	"github.com/opensourceways/xihe-server/bigmodel/domain/bigmodel"
	"github.com/opensourceways/xihe-server/bigmodel/domain/message"
	types "github.com/opensourceways/xihe-server/domain"
)

type AsyncBigModelService interface {
	WuKong(uint64, types.Account, *WuKongCmd) error
	GetIdleEndpoint(bid string) (int, error)
}

func NewAsyncBigModelService(
	fm bigmodel.BigModel,
	sender message.AsyncMessageProducer,
) AsyncBigModelService {
	return &asyncBigModelService{
		fm:     fm,
		sender: sender,
	}
}

type asyncBigModelService struct {
	fm     bigmodel.BigModel
	sender message.AsyncMessageProducer
}

func (s *asyncBigModelService) WuKong(tid uint64, user types.Account, cmd *WuKongCmd) (err error) {
	// 1. inference
	_ = s.sender.AddOperateLogForAccessBigModel(user, domain.BigmodelWuKong)

	msg := new(message.MsgTask)
	msg.WuKongAsyncTaskStart(tid, user.Account())
	s.sender.SendBigModelMsg(msg)

	// chose
	var estype string
	switch cmd.ImgQuantity {
	case 2:
		estype = string(domain.BigmodelWuKong)
	case 4:
		estype = string(domain.BigmodelWuKong4Img)
	}

	links, err := s.fm.GenPicturesByWuKong(user, &cmd.WuKongPictureMeta, estype)
	if err != nil {
		if !bigmodel.IsErrorSensitiveInfo(err) {
			err = errors.New("internal error")
		}

		msg_error := new(message.MsgTask)
		msg_error.WuKongInferenceError(tid, user.Account(), err.Error())
		s.sender.SendBigModelMsg(msg_error)

		return
	}

	// 3. send msg
	msg_finish := new(message.MsgTask)
	msg_finish.WuKongAsyncInferenceFinish(tid, user.Account(), links)

	return s.sender.SendBigModelMsg(msg_finish)
}

func (s *asyncBigModelService) GetIdleEndpoint(bid string) (c int, err error) {
	return s.fm.GetIdleEndpoint(bid)
}
