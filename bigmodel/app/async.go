package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	"github.com/opensourceways/xihe-server/bigmodel/domain/bigmodel"
	"github.com/opensourceways/xihe-server/bigmodel/domain/message"
	bigmodelmsg "github.com/opensourceways/xihe-server/bigmodel/domain/message"
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

	links, err := s.fm.GenPicturesByWuKong(user, (*domain.WuKongPictureMeta)(cmd))
	if err != nil {
		err = errors.New("access overload, please try again later")

		return
	}

	// 2. send msg
	return s.sender.UpdateWuKongTask(&bigmodelmsg.MsgWuKongLinks{
		Type:   "wukong",
		TaskId: tid,
		User:   user.Account(),
		Links:  links,
	})
}

func (s *asyncBigModelService) GetIdleEndpoint(bid string) (c int, err error) {
	return s.fm.GetIdleEndpoint(bid)
}
