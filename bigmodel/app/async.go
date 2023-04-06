package app

import (
	"github.com/opensourceways/xihe-server/bigmodel/domain/bigmodel"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
)

type AsyncBigModelService interface {
	WuKong(types.Account, *WuKongCmd) error
	GetIdleEndpoint(bid string) (int, error)
}

func NewAsyncBigModelService(
	fm bigmodel.BigModel,
	sender message.Sender,
) AsyncBigModelService {
	return &asyncBigModelService{
		fm:     fm,
		sender: sender,
	}
}

type asyncBigModelService struct {
	fm     bigmodel.BigModel
	sender message.Sender
}

func (s *asyncBigModelService) WuKong(types.Account, *WuKongCmd) (err error) {
	// 1. inference

	// 2. send msg

	return
}

func (s *asyncBigModelService) GetIdleEndpoint(bid string) (c int, err error) {
	return s.fm.GetIdleEndpoint(bid)
}
