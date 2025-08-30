package auth

import (
	"YAccount/configs"
	"YAccount/pkg/apperrors"
	"YAccount/pkg/auth"
	"YAccount/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWT认证中间件
func NewAuthMiddleware(config *configs.MiddlewareConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否需要跳过认证
		if shouldSkip(c.Request.URL.Path, config.JWT.SkipPaths) {
			c.Next()
			return
		}

		// 提取token
		token, err := extractToken(c)
		if err != nil {
			response.Error(c, err)
			c.Abort() // 终止请求链
			return
		}

		// 验证token
		claims, err := auth.ParseToken(token)
		if err != nil {
			// 直接使用response.Error处理，它会根据错误类型返回正确的HTTP状态码
			response.Error(c, err)
			c.Abort() // 终止请求链
			return
		}

		// 设置用户信息到上下文
		setUserContext(c, claims)
		c.Next()
	}
}

// 检查是否应该跳过认证
func shouldSkip(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}

// 提取token
func extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", apperrors.ErrUnauthorized
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return "", apperrors.ErrTokenInvalid
	}

	return tokenParts[1], nil
}

// 设置用户上下文
func setUserContext(c *gin.Context, claims *auth.Claims) {
	c.Set("user_id", claims.UserID)
	c.Set("username", claims.Username)
	c.Set("role", claims.Role)
}
