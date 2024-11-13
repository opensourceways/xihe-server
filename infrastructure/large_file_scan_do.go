package infrastructure

import "time"

type largeFileScanDO struct {
	Id               int64     `gorm:"column:id;primaryKey;autoIncrement"`
	Hash             string    `gorm:"column:hash;index"`
	Status           string    `gorm:"column:status"`
	Virus            string    `gorm:"column:virus"`
	SensitiveItem    string    `gorm:"column:sensitive_item"`
	ModerationStatus string    `gorm:"column:moderation_status"`
	ModerationResult string    `gorm:"column:moderation_result"`
	CreatedAt        time.Time `gorm:"column:created_at;<-:create"`
	UpdatedAt        time.Time `gorm:"column:updated_at"`
}

func (largeFileScanDO) TableName() string {
	return "large_file_scan"
}
