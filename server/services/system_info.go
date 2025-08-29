package services

import (
	"YAccount/global"
	"YAccount/models"
	"YAccount/pkg/apperrors"
	"YAccount/repositories"
	"YAccount/utils/logger"

	"github.com/go-redis/cache/v9"
)

func GetSystemInfoList() ([]models.SystemInfoList, error) {
	var systemInfoList []models.SystemInfoList

	if err := global.Cache.Once(&cache.Item{
		Key:   "system:info:list",
		Value: &systemInfoList,
		Do: func(*cache.Item) (any, error) {
			systemInfo, err := repositories.GetSystemInfoList()
			if err != nil {
				if apperrors.IsNotFoundError(err) {
					logger.LogError("GetSystemInfoList", "database query", "从数据库中获取系统配置列表失败", err)
				}
				return nil, err
			}
			for _, systemInfo := range systemInfo {
				systemInfoList = append(systemInfoList, models.SystemInfoList{
					ID:          systemInfo.ID,
					ConfigKey:   systemInfo.ConfigKey,
					ConfigValue: systemInfo.ConfigValue,
					ConfigType:  systemInfo.ConfigType,
					Description: systemInfo.Description,
					Status:      systemInfo.Status,
				})
			}
			return systemInfoList, nil
		},
	}); err != nil {
		return nil, apperrors.ErrSystemInfoListNotFound
	}
	return systemInfoList, nil
}
