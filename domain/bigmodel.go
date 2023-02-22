package domain

type UserLuoJiaRecord struct {
	User Account

	LuoJiaRecord
}

type LuoJiaRecord struct {
	Id        string
	CreatedAt int64
}

type LuoJiaRecordIndex struct {
	User Account
	Id   string
}

type WuKongPicture struct {
	Id        string
	Owner     Account
	OBSPath   string
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
