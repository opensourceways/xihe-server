package repositoryimpl

type TPod struct {
	Id        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CloudId   string `gorm:"column:cloud_id;not null"`
	Owner     string `gorm:"column:owner;not null"`
	Status    string `gorm:"column:status;not null"`
	Expiry    int64  `gorm:"column:expiry;not null"`
	Error     string `gorm:"column:error"`
	AccessURL string `gorm:"column:access_url"`
	CreatedAt int64  `gorm:"column:created_at;not null;default:extract(epoch from now())"`
	CardsNum  int    `gorm:"column:cards_num;not null"`
}

func (TPod) TableName() string {
	return "pod"
}
