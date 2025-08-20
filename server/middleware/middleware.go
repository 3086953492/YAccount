package middleware

import (
	"YAccount/pkg/apperrors"
	"YAccount/pkg/auth"
	"YAccount/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 Authorization 头获取令牌（正确处理 Bearer 格式）
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, apperrors.ErrUnauthorized)
			return
		}

		// 2. 验证并提取 Bearer Token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			response.Error(c, apperrors.ErrInvalidInput)
			return
		}
		token := tokenParts[1]

		claims, err := auth.ParseToken(token)
		if err != nil {
			response.Error(c, err)
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			response.Error(c, apperrors.ErrPermissionDenied)
			return
		}
		c.Next()
	}
}
