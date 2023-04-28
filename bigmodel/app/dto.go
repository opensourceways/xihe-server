package app

import (
	"errors"
	"io"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	"github.com/opensourceways/xihe-server/bigmodel/domain/bigmodel"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type DescribePictureCmd struct {
	User    types.Account
	Picture io.Reader
	Name    string
	Length  int64
}

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
	User  types.Account
	Level domain.WuKongPictureLevel
	WuKongPictureListOption
}

func (cmd *WuKongListPublicGlobalCmd) Validate() error {
	if cmd.WuKongPictureListOption.PageNum < 1 {
		return errors.New("page_num less than 1")
	}

	if cmd.WuKongPictureListOption.CountPerPage < 1 {
		return errors.New("count_per_page less than 1")
	}

	return nil
}

type WuKongCmd struct {
	domain.WuKongPictureMeta

	ImgQuantity int
}

func (cmd *WuKongCmd) Validate() error {
	if cmd.Desc == nil {
		return errors.New("invalid cmd")
	}

	if cmd.ImgQuantity != 2 && cmd.ImgQuantity != 4 {
		return errors.New("inalid cmd")
	}

	return nil
}

type WuKongICBCCmd struct {
	WuKongCmd

	User types.Account
}

type WuKongHFCmd struct {
	WuKongCmd

	EndPointType string
	User         types.Account
}

func (cmd *WuKongHFCmd) Validate() error {
	b := cmd.User == nil ||
		cmd.User.Account() != "wukong_hf" ||
		cmd.Desc == nil

	if b {
		return errors.New("invalid cmd")
	}

	return nil
}

type WuKongPicturesListCmd = repository.WuKongPictureListOption

type WuKongAddLikeFromTempCmd struct {
	User    types.Account
	OBSPath string
}

type WuKongAddLikeFromPublicCmd struct {
	User  types.Account
	Owner types.Account
	Id    string
}

type WuKongAddDiggCmd struct {
	User  types.Account
	Owner types.Account
	Id    string
}

type WuKongAddPublicFromTempCmd = WuKongAddLikeFromTempCmd

type WuKongAddPublicFromLikeCmd struct {
	User types.Account
	Id   string
}

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
	IsPublic bool   `json:"is_public"`
	Avatar   string `json:"avatar"`

	WuKongPictureBaseDTO
}

type WuKongPublicDTO struct { // public
	Avatar    string `json:"avatar"`
	IsLike    bool   `json:"is_like"`
	LikeID    string `json:"like_id"`
	IsDigg    bool   `json:"is_digg"`
	DiggCount int    `json:"digg_count"`

	WuKongPictureBaseDTO
}

func (dto *WuKongPublicDTO) toWuKongPublicDTO(
	p *domain.WuKongPicture, avatar string,
	isLike bool, likeId string, isDigg bool, link string,
) {
	*dto = WuKongPublicDTO{
		Avatar:    avatar,
		IsLike:    isLike,
		LikeID:    likeId,
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

type WuKongIsLikeDTO struct {
	IsLike bool
	LikeID string
}

type WuKongPublicGlobalDTO struct {
	Pictures []WuKongPublicDTO `json:"pictures"`
	Total    int               `json:"total"`
}

type WuKongLinksDTO struct {
	Pictures []string `json:"pictures"`
}

type WuKongRankDTO struct {
	Rank int `json:"rank"`
}
