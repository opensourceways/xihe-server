package repositoryimpl

type TAsyncTask struct {
	Id        uint64                 `gorm:"primaryKey;column:id"`
	User      string                 `gorm:"column:username"`
	TaskType  string                 `gorm:"column:task_type"`
	Status    string                 `gorm:"column:status"`
	CreatedAt int64                  `gorm:"column:created_at;default:extract(epoch from now())"`
	MetaData  map[string]interface{} `gorm:"type:json"`
}

func (w TAsyncTask) TableName() string {
	return "async_task"
}
