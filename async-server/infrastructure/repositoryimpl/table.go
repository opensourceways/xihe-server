package repositoryimpl

type TWukongTask struct {
	Id        uint64 `gorm:"primaryKey;column:id"`
	User      string `gorm:"column:username"`
	Style     string `gorm:"column:picture_style"`
	Desc      string `gorm:"column:description"`
	Status    string `gorm:"column:status"`
	CreatedAt int64  `gorm:"column:created_at;default:extract(epoch from now())"`
}

func (w TWukongTask) TableName() string {
	return "wukong_task"
}
