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

	s.sender.UpdateWuKongTask(&message.MsgTask{
		Type:   "wukong_update",
		TaskId: tid,
		Status: "running",
	})

	links, err := s.fm.GenPicturesByWuKong(user, (*domain.WuKongPictureMeta)(cmd))
	if err != nil {
		if !bigmodel.IsErrorSensitiveInfo(err) {
			err = errors.New("internal error")
		}

		msg := new(message.MsgTask)
		msg.SetErrorMsgTask(tid, user.Account(), err.Error())

		s.sender.UpdateWuKongTask(msg)

		return
	}

	// 3. send msg
	return s.sender.UpdateWuKongTask(&message.MsgTask{
		Type:    "wukong_update",
		TaskId:  tid,
		User:    user.Account(),
		Status:  "finished",
		Details: links,
	})
}

func (s *asyncBigModelService) GetIdleEndpoint(bid string) (c int, err error) {
	return s.fm.GetIdleEndpoint(bid)
}
