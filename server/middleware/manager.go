package middleware

import (
	"YAccount/configs"
	"YAccount/global"
	"YAccount/middleware/auth"
	"YAccount/middleware/oauth"
	"YAccount/middleware/permission"
	"YAccount/middleware/security"

	"github.com/gin-gonic/gin"
)

// 中间件管理器 - 先做个简单版本
type Manager struct {
	config *configs.MiddlewareConfig
}

// 创建管理器
func NewManager() *Manager {

	config := global.Cfg.Middleware

	return &Manager{
		config: &config,
	}
}

// 加载所有全局中间件
func (m *Manager) LoadGlobal(engine *gin.Engine) {
	// CORS中间件
	engine.Use(m.CORS())

	// 添加Recovery中间件防止panic
	engine.Use(gin.Recovery())

	// 这里以后可以加更多全局中间件
	// engine.Use(m.Logger())
}

// 获取认证中间件
func (m *Manager) Auth() gin.HandlerFunc {
	return auth.NewAuthMiddleware(m.config)
}

// 获取管理员权限中间件
func (m *Manager) AdminPermission() gin.HandlerFunc {
	return permission.NewAdminPermissionMiddleware()
}

// CORS中间件
func (m *Manager) CORS() gin.HandlerFunc {
	return security.NewCORSMiddleware(security.CORSConfig{
		AllowOrigins: m.config.CORS.AllowOrigins,
		AllowMethods: m.config.CORS.AllowMethods,
		AllowHeaders: m.config.CORS.AllowHeaders,
	})
}

// 获取 OAuth 中间件
func (m *Manager) OAuth(requiredScopes ...string) gin.HandlerFunc {
	return oauth.OAuthTokenMiddleware(requiredScopes...)
}