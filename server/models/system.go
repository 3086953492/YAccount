package models

import "time"

type System struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	ConfigKey   string    `gorm:"not null;size:100" json:"config_key"`
	ConfigValue string    `gorm:"not null" json:"config_value"`
	ConfigType  string    `gorm:"not null;size:50" json:"config_type"`
	Description string    `gorm:"size:255" json:"description"`
	Status      int       `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   uint      `json:"created_by"`
	UpdatedBy   uint      `json:"updated_by"`
}

func (System) TableName() string {
	return "system"
}
