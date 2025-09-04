package services

import (
	"YAccount/global"
	"YAccount/models"
	apperrors "github.com/3086953492/YaBase/errors"
	"YAccount/repositories"
	"YAccount/utils/logger"
	cache_utils "YAccount/utils/cache"

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
				if !apperrors.IsNotFoundError(err) {
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

func GetSystemInfoByKey(key string) (*models.SystemInfo, error) {
	var systemInfo models.SystemInfo

	if err := global.Cache.Once(&cache.Item{
		Key:   "system:info:" + key,
		Value: &systemInfo,
		Do: func(*cache.Item) (any, error) {
			systemInfo, err := repositories.GetSystemInfo(key)
			if err != nil {
				if !apperrors.IsNotFoundError(err) {
					logger.LogError("GetSystemInfoByKey", "database query", "从数据库中获取系统配置失败", err)
				}
				return nil, err
			}
			return systemInfo, nil
		},
	}); err != nil {
		return nil, apperrors.ErrSystemInfoNotFound
	}
	return &systemInfo, nil
}

func BatchUpdateSystemInfo(req *models.BatchUpdateSystemInfoRequest) error {
	// 开启数据库事务
	tx := global.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 确保在函数返回时处理事务
	defer func() {
		if r := recover(); r != nil {
			// 发生 panic 时回滚事务
			tx.Rollback()
			logger.LogError("BatchUpdateSystemInfo", "panic recovery", "批量更新系统配置时发生 panic，事务已回滚", nil)
		}
	}()

	// 批量更新系统配置
	for _, systemInfo := range req.SystemInfos {
		if err := repositories.UpdateSystemInfo(&systemInfo, tx); err != nil {
			// 更新失败时回滚事务
			tx.Rollback()
			logger.LogError("BatchUpdateSystemInfo", "database update", "批量更新系统配置失败，事务已回滚", err)
			return err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		logger.LogError("BatchUpdateSystemInfo", "transaction commit", "提交事务失败", err)
		return err
	}

	// 清空缓存
	cache_utils.DeleteCacheKeysByPrefix("system:info")

	return nil
}
