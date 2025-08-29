package models

import "time"

type SystemInfo struct {
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

func (SystemInfo) TableName() string {
	return "system_infos"
}

type SystemInfoList struct {
	ID          uint   `json:"id"`
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
	ConfigType  string `json:"config_type"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

type UpdateSystemInfoRequest struct {
	ID          uint   `json:"id"`
	ConfigValue string `json:"config_value"`
}

type BatchUpdateSystemInfoRequest struct {
	SystemInfos []UpdateSystemInfoRequest `json:"system_infos"`
}
