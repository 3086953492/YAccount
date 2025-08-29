package repositories

import (
	"YAccount/global"
	"YAccount/models"
)

func GetSystemInfo(key string) (*models.SystemInfo, error) {
	var system models.SystemInfo
	if err := global.DB.Where("config_key = ?", key).First(&system).Error; err != nil {
		return nil, err
	}
	return &system, nil
}

func GetSystemInfoList() ([]models.SystemInfo, error) {
	var system []models.SystemInfo
	if err := global.DB.Find(&system).Error; err != nil {
		return nil, err
	}
	return system, nil
}