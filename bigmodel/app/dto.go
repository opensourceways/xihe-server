package app

import (
	"errors"
	"fmt"
	"io"

	"github.com/opensourceways/xihe-server/bigmodel/domain"
	bigmodeldomain "github.com/opensourceways/xihe-server/bigmodel/domain"
	"github.com/opensourceways/xihe-server/bigmodel/domain/bigmodel"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	userdomain "github.com/opensourceways/xihe-server/user/domain"
	"github.com/opensourceways/xihe-server/utils"
)

type DescribePictureCmd struct {
	User    types.Account
	Picture io.Reader
	Name    string
	Length  int64
}

func (cmd *DescribePictureCmd) Validate() error {
	if !utils.IsSafeFileName(cmd.Name) {
		return errors.New("file name invalid")
	}

	return nil
}

type VQAHFCmd struct {
	User    types.Account
	Picture io.Reader
	Ask     string
}

func (cmd *VQAHFCmd) Validate() error {
	if cmd.Picture == nil || cmd.Ask == "" {
		return errors.New("invalid cmd")
	}

	cmd.Ask = utils.XSSFilter(cmd.Ask)

	return nil
}

type LuoJiaHFCmd struct {
	User    userdomain.Account
	Picture io.Reader
}

func (cmd *LuoJiaHFCmd) Validate() error {
	if cmd.Picture == nil {
		return errors.New("invalid cmd")
	}

	return nil
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

	EsType string
}

func (cmd *WuKongCmd) Validate() error {
	cmd.Style = utils.XSSFilter(cmd.Style)

	if max := 4 * 4; utils.StrLen(cmd.Style) > max {
		return fmt.Errorf("style should less than %d", max)
	}

	return nil
}

type WuKongHFCmd struct {
	WuKongCmd

	EndPointType string
	User         types.Account
}

func (cmd *WuKongHFCmd) Validate() error {
	cmd.WuKongCmd.Style = utils.XSSFilter(cmd.WuKongCmd.Style)

	b := cmd.User == nil ||
		cmd.User.Account() != "wukong_hf" ||
		cmd.Desc == nil

	if b {
		return errors.New("invalid cmd")
	}

	return nil
}

type WuKongApiCmd struct {
	WuKongCmd

	EndPointType string
	User         types.Account
}

func (cmd *WuKongApiCmd) Validate() error {
	cmd.WuKongCmd.Style = utils.XSSFilter(cmd.WuKongCmd.Style)

	b := cmd.Desc == nil

	if b {
		return errors.New("invalid cmd")
	}

	return nil
}

type WuKongPicturesListCmd = repository.WuKongPictureListOption

type WuKongAddLikeFromTempCmd struct {
	User    types.Account
	OBSPath bigmodeldomain.OBSPath
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

type wukongPictureDTO struct {
	Link     string `json:"link"`
	IsPublic bool   `json:"is_public"`
	PublicID string `json:"public_id"`
	IsLike   bool   `json:"is_like"`
	LikeID   string `json:"like_id"`
}

type WuKongRankDTO struct {
	Rank int `json:"rank"`
}

type AIDetectorCmd struct {
	User types.Account         `json:"user"`
	Lang domain.Lang           `json:"lang"`
	Text domain.AIDetectorText `json:"text"`
}

func (cmd AIDetectorCmd) Validate() error {
	input := domain.AIDetectorInput{
		Lang: cmd.Lang,
		Text: cmd.Text,
	}

	if !input.IsTextLengthOK() {
		return errors.New("text length too long")
	}

	return nil
}

// taichu
type GenPictureCmd struct {
	User types.Account `json:"user"`
	Desc domain.Desc   `json:"desc"`
}

type ApplyRecordGetCmd struct {
	User      types.Account
	ModelName domain.ModelName
}

type ApiApplyRecordDTO struct {
	User      string `json:"user"`
	Name      string `json:"name"`
	Endpoint  string `json:"endpoint"`
	Token     string `json:"token"`
	ApplyAt   string `json:"apply_at"`
	ModelName string `json:"model_name"`
	Enabled   bool   `json:"enabled"`
}

func (s bigModelService) toApiApplyRecordDTO(
	v *domain.UserApiRecord,
	dto *ApiApplyRecordDTO,
	info domain.ApiInfo,
) {
	*dto = ApiApplyRecordDTO{
		User:      v.User.Account(),
		Name:      info.Name,
		ApplyAt:   v.ApplyAt,
		Token:     v.Token,
		ModelName: v.ModelName.ModelName(),
		Enabled:   v.Enabled,
		Endpoint:  info.Endpoint,
	}
}

type ApiInfoDTO struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Doc      string `json:"doc"`
}

func (s bigModelService) toApiInfoDTO(
	v *domain.ApiInfo,
	dto *ApiInfoDTO,
) {
	*dto = ApiInfoDTO{
		Id:       v.Id,
		Name:     v.Name,
		Endpoint: v.Endpoint,
		Doc:      v.Doc.URL(),
	}
}

// baichuan
type BaiChuanCmd struct {
	User              types.Account
	Sampling          bool
	Text              domain.BaiChuanText
	TopK              domain.TopK
	TopP              domain.TopP
	Temperature       domain.Temperature
	RepetitionPenalty domain.RepetitionPenalty
}

func (cmd *BaiChuanCmd) SetDefault() {
	cmd.TopK, _ = domain.NewTopK(5)
	cmd.TopP, _ = domain.NewTopP(0.85)
	cmd.Temperature, _ = domain.NewTemperature(0.3)
	cmd.RepetitionPenalty, _ = domain.NewRepetitionPenalty(1.05)
}

type BaiChuanDTO struct {
	Text string `json:"text"`
}
