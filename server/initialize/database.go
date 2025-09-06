package initialize

import (
	"YAccount/global"
	"YAccount/models"
	"fmt"
	apperrors "github.com/3086953492/YaBase/errors"
	ybase_global "github.com/3086953492/YaBase/global"
	logger_utils "github.com/3086953492/YaBase/logger"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logger_gorm "gorm.io/gorm/logger"
)

func InitDB() error {

	cfg := ybase_global.GetGlobalConfig().Database

	userDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
		cfg.ParseTime,
		cfg.Loc,
	)

	var err error
	global.DB, err = gorm.Open(mysql.Open(userDsn), &gorm.Config{
		Logger: logger_gorm.Default.LogMode(logger_gorm.Info),
	})
	if err != nil {
		logger_utils.Error("数据库连接失败",
			zap.String("operation", "database_connection"),
			zap.Error(err),
			zap.String("host", cfg.Host),
			zap.Int("port", cfg.Port),
		)
		return fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := global.DB.DB()
	if err != nil {
		logger_utils.Error("获取底层数据库连接失败",
			zap.String("operation", "get_sql_db"),
			zap.Error(err),
		)
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := AutoMigrate(); err != nil {
		logger_utils.Error("数据库自动迁移失败",
			zap.String("operation", "auto_migrate"),
			zap.Error(err),
		)
		return apperrors.ErrAutoMigrate
	}

	return nil
}

func AutoMigrate() error {
	modelsList := []any{
		&models.User{},
		&models.SystemInfo{},
		// 新增 OAuth 相关表
		&models.OAuthClient{},
		&models.OAuthAuthorizationCode{},
		&models.OAuthAccessToken{},
		&models.OAuthScope{},
	}

	logger_utils.Info("开始迁移数据表结构",
		zap.String("operation", "auto_migrate_tables"),
		zap.Int("tables_count", len(modelsList)),
	)

	err := global.DB.AutoMigrate(modelsList...)
	if err != nil {
		logger_utils.Error("数据表结构迁移失败",
			zap.String("operation", "auto_migrate_tables"),
			zap.Error(err),
		)
		return err
	}

	logger_utils.Info("数据表结构迁移成功",
		zap.String("operation", "auto_migrate_tables"),
		zap.Strings("tables", []string{"User", "SystemInfo", "OAuthClient", "OAuthAuthorizationCode", "OAuthAccessToken", "OAuthScope"}),
	)

	return nil
}
