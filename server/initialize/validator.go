// server/initialize/validator.go
package initialize

import (
	"fmt"

	"YAccount/validations/user"

	"github.com/3086953492/YaBase/validator"
)

func InitValidator() error {
	registry := validator.NewAutoRegistry()

	err := registry.
		RegisterPackage("user", &user.Validators{}).
		Apply()

	if err != nil {
		return fmt.Errorf("验证器注册失败: %w", err)
	}

	// 手动测试验证器是否注册成功
	v := validator.GetValidator()
	if v == nil {
		return fmt.Errorf("获取验证器实例失败")
	}

	return nil
}
