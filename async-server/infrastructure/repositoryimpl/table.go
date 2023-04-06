package repositoryimpl

type WukongRequest struct {
	Id        uint64 `gorm:"primaryKey;column:id"`
	User      string `gorm:"column:user"`
	Style     string `gorm:"column:style"`
	Desc      string `gorm:"column:desc"`
	CreatedAt int64  `gorm:"column:created_at;default:extract(epoch from now())"`
}

func (w WukongRequest) TableName() string {
	return "wukong_request"
}
