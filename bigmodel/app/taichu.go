package app

import (
	"io"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

func (s bigModelService) DescribePicture(
	user types.Account, picture io.Reader, name string, length int64,
) (string, error) {
	_ = s.sender.AddOperateLogForAccessBigModel(user, domain.BigmodelDescPicture)

	return s.fm.DescribePicture(picture, name, length, string(domain.BigmodelDescPicture))
}

func (s bigModelService) DescribePictureHF(
	cmd *DescribePictureCmd,
) (string, error) {
	_ = s.sender.AddOperateLogForAccessBigModel(cmd.User, domain.BigmodelDescPicture)

	return s.fm.DescribePicture(cmd.Picture, cmd.Name, cmd.Length, string(domain.BigmodelDescPictureHF))
}

func (s bigModelService) GenPicture(
	cmd GenPictureCmd,
) (link string, code string, err error) {
	_ = s.sender.AddOperateLogForAccessBigModel(cmd.User, domain.BigmodelGenPicture)

	if link, err = s.fm.GenPicture(cmd.User, cmd.Desc.Desc()); err != nil {
		code = s.setCode(err)
	}

	return
}

func (s bigModelService) GenPictures(
	cmd GenPictureCmd,
) (links []string, code string, err error) {
	_ = s.sender.AddOperateLogForAccessBigModel(cmd.User, domain.BigmodelGenPicture)

	if links, err = s.fm.GenPictures(cmd.User, cmd.Desc.Desc()); err != nil {
		code = s.setCode(err)
	}

	return
}

func (s bigModelService) Ask(
	u types.Account, q domain.Question, f string,
) (v string, code string, err error) {
	_ = s.sender.AddOperateLogForAccessBigModel(u, domain.BigmodelVQA)

	if v, err = s.fm.Ask(q, f); err != nil {
		code = s.setCode(err)
	}

	return
}

func (s bigModelService) VQAUploadPicture(f io.Reader, user types.Account, fileName string) error {
	return s.fm.VQAUploadPicture(f, user, fileName)
}

func (s bigModelService) VQAHF(cmd *VQAHFCmd) (v string, code string, err error) {
	_ = s.sender.AddOperateLogForAccessBigModel(cmd.User, domain.BigmodelVQA)

	if v, err = s.fm.AskHF(cmd.Picture, cmd.User, cmd.Ask); err != nil {
		code = s.setCode(err)
	}

	return
}
