// server/validations/user/validators.go
package user

import (
	"YAccount/repositories"
	"fmt"

	"github.com/3086953492/YaBase/validator"
)

type Validators struct{}

// GetValidators 实现ValidatorPackage接口
func (v *Validators) GetValidators() map[string]validator.ValidatorFunc {
	methods := validator.ExtractValidatorMethods(v)

	for name := range methods {
		fmt.Printf("- %s\n", name)
	}

	return methods
}

// UsernameUnique 用户名唯一性验证 -> username_unique
func (v *Validators) UsernameUnique(fl validator.FieldLevel) bool {
	username := fl.Field().Interface().(string)
	if username == "" {
		return true
	}
	_, err := repositories.GetUserByUsername(username)
	return err != nil
}
