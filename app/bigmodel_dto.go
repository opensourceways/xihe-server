package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/bigmodel"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type CodeGeexDTO = bigmodel.CodeGeexResp

type CodeGeexCmd bigmodel.CodeGeexReq

func (cmd *CodeGeexCmd) Validate() error {
	if cmd.Content == "" || cmd.Lang == "" {
		return errors.New("invalid cmd")
	}

	return nil
}

type LuoJiaRecordDTO struct {
	CreatedAt string `json:"created_at"`
	Id        string `json:"id"`
}

type WuKongCmd domain.WuKongPictureMeta

func (cmd *WuKongCmd) Validate() error {
	if cmd.Desc == nil {
		return errors.New("invalid cmd")
	}

	return nil
}

type WuKongPicturesListCmd = repository.WuKongPictureListOption

type WuKongAddLikeFromTempCmd struct {
	User    domain.Account
	OBSPath string
}

type WuKongAddLikeFromPublicCmd struct {
	User domain.Account
	Id   string
}

type WuKongGetPublicCmd = WuKongAddLikeFromPublicCmd

type WuKongPictureBaseDTO struct {
	Id        string `json:"id"`
	Owner     string `json:"owner"` // owner of picture
	Desc      string `json:"desc"`
	Style     string `json:"style"`
	Link      string `json:"link"`
	CreatedAt string `json:"created_at"`
}

type WuKongLikeDTO struct { // like
	IsPublic bool `json:"is_public"`
	IsDigg   bool `json:"is_digg"`

	WuKongPictureBaseDTO
}

type WuKongPublicDTO struct { // public
	IsLike bool `json:"is_like"`
	IsDigg bool `json:"is_digg"`

	WuKongPictureBaseDTO
}
