package validator

import (
	"YAccount/global"

	"github.com/gin-gonic/gin"
)

func ValidateStruct(c *gin.Context, req any) bool {

	if err := c.ShouldBindJSON(req); err != nil {
		return false
	}

	if err := global.Validate.Struct(req); err != nil {
		return false
	}

	return true
}
