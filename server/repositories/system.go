package repositories

import (
	"YAccount/global"
	"YAccount/models"
)

func GetSystemConfig(key string) (*models.System, error) {
	var system models.System
	if err := global.DB.Where("config_key = ?", key).First(&system).Error; err != nil {
		return nil, err
	}
	return &system, nil
}

func GetSystemConfigList() ([]models.System, error) {
	var system []models.System
	if err := global.DB.Find(&system).Error; err != nil {
		return nil, err
	}
	return system, nil
}