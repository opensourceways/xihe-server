package app

import (
	"io"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

func (s bigModelService) LuoJiaUploadPicture(f io.Reader, user types.Account) error {
	return s.fm.LuoJiaUploadPicture(f, user)
}

func (s bigModelService) LuoJia(user types.Account) (v string, err error) {
	_ = s.sender.AddOperateLogForAccessBigModel(user, domain.BigmodelLuoJia)

	if v, err = s.fm.LuoJia(user.Account()); err != nil {
		return
	}

	record := domain.UserLuoJiaRecord{User: user}
	record.CreatedAt = utils.Now()

	s.luojia.Save(&record)

	return
}

func (s bigModelService) ListLuoJiaRecord(user types.Account) (
	dtos []LuoJiaRecordDTO, err error,
) {
	v, err := s.luojia.List(user)
	if err != nil || len(v) == 0 {
		return
	}

	dtos = append(dtos, LuoJiaRecordDTO{
		CreatedAt: utils.ToDate(v[0].CreatedAt),
		Id:        v[0].Id,
	})

	return
}
