package repositoryimpl

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type TAsyncTask struct {
	Id        uint64  `gorm:"primaryKey;column:id"`
	User      string  `gorm:"column:username"`
	TaskType  string  `gorm:"column:task_type"`
	Status    string  `gorm:"column:status"`
	CreatedAt int64   `gorm:"column:created_at;default:extract(epoch from now())"`
	MetaData  JSONMap `gorm:"column:metadata;type:json;default: '{}'::json"`
}

func NewTAsyncTask() *TAsyncTask {
	return &TAsyncTask{
		MetaData: make(JSONMap),
	}
}

// JSONMap
type JSONMap map[string]interface{}

// Value implement driver.Valuer interface
func (j JSONMap) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	b, err := json.Marshal(j)
	return string(b), err
}

// Scan implement sql.Scanner interfacne
func (j *JSONMap) Scan(src interface{}) error {
	if src == nil {
		*j = nil
		return nil
	}

	// string data
	var jsonData string

	if s, ok := src.(string); ok {
		jsonData = s
	} else if b, ok := src.([]byte); ok {
		jsonData = string(b)
	} else {
		return fmt.Errorf("incompatible type for JSONMap")
	}

	// string data marshal to JSONMap
	err := json.Unmarshal([]byte(jsonData), &j)
	if err != nil {
		return fmt.Errorf("invalid JSON data: %v", err)
	}

	return nil
}

func (w TAsyncTask) TableName() string {
	return "async_task"
}
