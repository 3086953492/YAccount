package services

import (
	"YAccount/global"
	"YAccount/models"
	"YAccount/pkg/apperrors"
	"YAccount/repositories"
	"YAccount/utils/logger"

	"github.com/go-redis/cache/v9"
)

func GetSystemConfigList() ([]models.SystemConfigList, error) {
	var systemConfigList []models.SystemConfigList

	if err := global.Cache.Once(&cache.Item{
		Key:   "system:config:list",
		Value: &systemConfigList,
		Do: func(*cache.Item) (any, error) {
			systemConfig, err := repositories.GetSystemConfigList()
			if err != nil {
				if apperrors.IsNotFoundError(err) {
					logger.LogError("GetSystemConfigList", "database query", "从数据库中获取系统配置列表失败", err)
				}
				return nil, err
			}
			for _, systemConfig := range systemConfig {
				systemConfigList = append(systemConfigList, models.SystemConfigList{
					ID:          systemConfig.ID,
					ConfigKey:   systemConfig.ConfigKey,
					ConfigValue: systemConfig.ConfigValue,
					ConfigType:  systemConfig.ConfigType,
					Description: systemConfig.Description,
					Status:      systemConfig.Status,
				})
			}
			return systemConfigList, nil
		},
	}); err != nil {
		return nil, apperrors.ErrSystemConfigListNotFound
	}
	return systemConfigList, nil
}
