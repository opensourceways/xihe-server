package controller

type userBasicInfoModel struct {
	Nickname string `json:"nickname"`
	AvatarId string `json:"avatar_id"`
	Bio      string `json:"bio"`
}
