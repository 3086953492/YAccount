package repositories

import (
	"YAccount/global"
	"YAccount/models"
)

func GetSystemConfig(key string) (*models.SystemConfig, error) {
	var system models.SystemConfig
	if err := global.DB.Where("config_key = ?", key).First(&system).Error; err != nil {
		return nil, err
	}
	return &system, nil
}

func GetSystemConfigList() ([]models.SystemConfig, error) {
	var system []models.SystemConfig
	if err := global.DB.Find(&system).Error; err != nil {
		return nil, err
	}
	return system, nil
}