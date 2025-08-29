package repositories

import (
	"YAccount/global"
	"YAccount/models"

	"gorm.io/gorm"
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

func UpdateSystemInfo(req *models.UpdateSystemInfoRequest, tx *gorm.DB) error {
	if tx == nil {
		tx = global.DB
	}
	if err := tx.Model(&models.SystemInfo{}).Where("id = ?", req.ID).Updates(req).Error; err != nil {
		return err
	}
	return nil
}