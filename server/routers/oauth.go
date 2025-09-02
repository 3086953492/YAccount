package routers

import (
	"YAccount/controllers"
	"YAccount/middleware"

	"github.com/gin-gonic/gin"
)

func LoadOAuthRouters(router *gin.Engine) {
	m := middleware.NewManager()

	// OAuth 标准端点
	oauthGroup := router.Group("/api/account/v1/oauth")
	{
		// 授权端点（需要用户登录）
		oauthGroup.GET("/authorize", m.Auth(), controllers.OAuthAuthorizeHandler)

		// 令牌端点（公开）
		oauthGroup.POST("/token", controllers.OAuthTokenHandler)

		// 令牌内省端点（公开）
		oauthGroup.POST("/introspect", controllers.OAuthIntrospectHandler)
	}

	clientGroup := router.Group("/api/account/v1/oauth/clients")
	{
		clientGroup.POST("", m.Auth(), m.AdminPermission(), controllers.OAuthClientRegisterHandler)
		clientGroup.GET("", m.Auth(), m.AdminPermission(), controllers.ListOAuthClientsHandler)
		clientGroup.GET("/:client_id", m.Auth(), m.AdminPermission(), controllers.GetOAuthClientHandler)
	}

}
