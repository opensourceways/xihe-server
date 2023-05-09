package domain

import types "github.com/opensourceways/xihe-server/domain"

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

func (r *WuKongPicture) SetDefaultDiggs() {
	r.Diggs = []string{}
}
