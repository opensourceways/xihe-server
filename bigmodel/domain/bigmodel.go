package domain

import (
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

// luojia
type UserLuoJiaRecord struct {
	User types.Account

	LuoJiaRecord
}

type LuoJiaRecord struct {
	Id        string
	CreatedAt int64
}

type LuoJiaRecordIndex struct {
	User types.Account
	Id   string
}

// wukong
type WuKongPicture struct {
	Id        string
	Owner     types.Account
	OBSPath   OBSPath
	Level     WuKongPictureLevel
	Diggs     []string
	DiggCount int
	Version   int
	CreatedAt string

	WuKongPictureMeta
}

type WuKongPictureMeta struct {
	Style string
	Desc  WuKongPictureDesc
}

func (r *WuKongPicture) IsOfficial() bool {
	return r.Level.IsOfficial()
}

func (r WuKongPicture) IsHot() bool {
	return r.Level.IsHot()
}

func (r *WuKongPicture) SetDefaultDiggs() {
	r.Diggs = []string{}
}

// ai detector
type AIDetectorInput struct {
	Lang Lang
	Text AIDetectorText
}

func (r AIDetectorInput) IsTextLengthOK() bool {
	if r.Lang.IsEN() {
		return utils.StrLen(r.Text.AIDetectorText()) <= 2000
	}

	if r.Lang.IsZH() {
		return utils.StrLen(r.Text.AIDetectorText()) <= 2000
	}

	return false
}

// taichu
type GenPictureInput struct {
	Desc Desc
}

// BaiChuan
type BaiChuanInput struct {
	Text              BaiChuanText
	Sampling          bool
	TopK              TopK
	TopP              TopP
	Temperature       Temperature
	RepetitionPenalty RepetitionPenalty
}

// glm2
type GLM2Input struct {
	Text              GLM2Text
	Sampling          bool
	History           []History
	TopK              TopK
	TopP              TopP
	Temperature       Temperature
	RepetitionPenalty RepetitionPenalty
}

// llama2
type LLAMA2Input struct {
	Text              LLAMA2Text
	Sampling          bool
	History           []History
	TopK              TopK
	TopP              TopP
	Temperature       Temperature
	RepetitionPenalty RepetitionPenalty
}

// skywork 13b
type SkyWorkInput struct {
	Text              SkyWorkText
	Sampling          bool
	History           []History
	TopK              TopK
	TopP              TopP
	Temperature       Temperature
	RepetitionPenalty RepetitionPenalty
}

// iflytekspark
type IFlytekSparkInput struct {
	Text              IFlytekSparkText
	Sampling          bool
	History           []History
	TopK              TopK
	TopP              TopP
	Temperature       Temperature
	RepetitionPenalty RepetitionPenalty
}
