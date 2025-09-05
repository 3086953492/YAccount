package global

import (
	"YAccount/configs"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"github.com/go-redis/cache/v9"
	"gorm.io/gorm"
)

var (
	Cfg      *configs.Config
	DB       *gorm.DB
	Validate *validator.Validate
	Redis    *redis.Client
	Cache    *cache.Cache
)
