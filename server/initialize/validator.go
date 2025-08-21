package initialize

import (
	"YAccount/global"
	validator_pkg "YAccount/pkg/validator"

	"github.com/go-playground/validator/v10"
)

func InitValidator() error {

	global.Validate = validator.New()

	if err := global.Validate.RegisterValidation("usernameUnique", validator_pkg.VerifyUsernameUnique); err != nil {
		return err
	}

	return nil
}
