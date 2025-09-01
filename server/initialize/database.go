package initialize

import (
	"YAccount/global"
	"YAccount/models"
	"YAccount/pkg/apperrors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() error {
	userCfg := global.Cfg.Database

	userDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		userCfg.User,
		userCfg.Password,
		userCfg.Host,
		userCfg.Port,
		userCfg.DBName,
		userCfg.Charset,
		userCfg.ParseTime,
		userCfg.Loc,
	)

	var err error
	global.DB, err = gorm.Open(mysql.Open(userDsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := global.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := AutoMigrate(); err != nil {
		global.Logger.Error("failed to auto migrate", zap.Error(err))
		return apperrors.ErrAutoMigrate
	}

	return nil
}

func AutoMigrate() error {
	return global.DB.AutoMigrate(
		&models.User{},
		&models.SystemInfo{},
		// 新增 OAuth 相关表
		&models.OAuthClient{},
		&models.OAuthAuthorizationCode{},
		&models.OAuthAccessToken{},
		&models.OAuthScope{},
	)
}
