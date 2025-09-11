package repositories

import (
	"YAccount/models"

	"github.com/3086953492/YaBase/database"
	"gorm.io/gorm"
)

func systemInfoDB() *gorm.DB {
	return database.GetGlobalDB()
}

func GetSystemInfo(key string) (*models.SystemInfo, error) {
	var system models.SystemInfo
	if err := systemInfoDB().Where("config_key = ?", key).First(&system).Error; err != nil {
		return nil, err
	}
	return &system, nil
}

func GetSystemInfoList() ([]models.SystemInfo, error) {
	var system []models.SystemInfo
	if err := systemInfoDB().Find(&system).Error; err != nil {
		return nil, err
	}
	return system, nil
}

func UpdateSystemInfo(req *models.UpdateSystemInfoRequest, tx *gorm.DB) error {
	if tx == nil {
		tx = systemInfoDB()
	}
	if err := tx.Model(&models.SystemInfo{}).Where("id = ?", req.ID).Updates(req).Error; err != nil {
		return err
	}
	return nil
}