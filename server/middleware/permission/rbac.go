package permission

import (
	"YAccount/pkg/apperrors"
	"YAccount/pkg/response"

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
