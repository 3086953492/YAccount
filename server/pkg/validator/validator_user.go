package validator

import (
	"YAccount/repositories"

	"github.com/go-playground/validator/v10"
)

func VerifyUsernameUnique(fl validator.FieldLevel) bool {

	username := fl.Field().Interface().(string)

	// 如果有错误，则说明用户名未被使用，返回true
	if _, err := repositories.GetUserByUsername(username); err != nil {
		return true
	}

	return false
}
