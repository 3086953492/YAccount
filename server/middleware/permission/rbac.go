package permission

import (
	apperrors "github.com/3086953492/YaBase/errors"
	"github.com/3086953492/YaBase/response"

	"github.com/gin-gonic/gin"
)


func NewAdminPermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户信息
		role := c.GetString("role")
		if role != "admin" {
			response.Error(c, apperrors.ErrPermissionDenied)
			c.Abort() // 终止请求链
			return
		}
		c.Next()
	}
}
