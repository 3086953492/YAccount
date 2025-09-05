package global

import (
	"YAccount/configs"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var (
	Cfg      *configs.Config
	DB       *gorm.DB
	Validate *validator.Validate
)
