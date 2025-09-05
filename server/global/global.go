package global

import (
	"github.com/3086953492/YaBase/configs"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var (
	Cfg      *configs.Config
	DB       *gorm.DB
	Validate *validator.Validate
)
