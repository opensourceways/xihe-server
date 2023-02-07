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

type WuKongPictureListOption struct {
	CountPerPage int
	PageNum      int
}

type WuKongListPublicGlobalCmd struct {
	User domain.Account
	WuKongPictureListOption
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

type WuKongAddDiggCmd struct {
	User  domain.Account
	Owner domain.Account
	Id    string
}

type WuKongGetPublicCmd = WuKongAddLikeFromPublicCmd

type WuKongAddPublicFromTempCmd = WuKongAddLikeFromTempCmd

type WuKongAddPublicFromLikeCmd = WuKongAddLikeFromPublicCmd

type WuKongCancelDiggCmd WuKongAddDiggCmd

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

	WuKongPictureBaseDTO
}

type WuKongPublicDTO struct { // public
	Avatar    string `json:"avatar"`
	IsLike    bool   `json:"is_like"`
	IsDigg    bool   `json:"is_digg"`
	DiggCount int    `json:"digg_count"`

	WuKongPictureBaseDTO
}

func (dto *WuKongPublicDTO) toWuKongPublicDTO(
	p *domain.WuKongPicture, avatar string,
	isLike bool, isDigg bool, link string,
) {
	*dto = WuKongPublicDTO{
		Avatar:    avatar,
		IsLike:    isLike,
		IsDigg:    isDigg,
		DiggCount: p.DiggCount,

		WuKongPictureBaseDTO: WuKongPictureBaseDTO{
			Id:        p.Id,
			Owner:     p.Owner.Account(),
			Desc:      p.Desc.WuKongPictureDesc(),
			Style:     p.Style,
			Link:      link,
			CreatedAt: p.CreatedAt,
		},
	}
}

type WuKongPublicGlobalDTO struct {
	Pictures []WuKongPublicDTO `json:"pictures"`
	Total    int               `json:"total"`
}
