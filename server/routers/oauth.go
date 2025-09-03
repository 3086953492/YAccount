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
		oauthGroup.GET("/authorize", m.OAuth(), controllers.OAuthAuthorizeHandler)

		// 令牌端点（公开）
		oauthGroup.POST("/token", controllers.OAuthTokenHandler)

		// 令牌内省端点（公开）
		oauthGroup.POST("/introspect", controllers.OAuthIntrospectHandler)
	}

	clientGroup := router.Group("/api/account/v1/oauth/clients")
	{
		clientGroup.POST("", m.OAuth(), m.RequiredScopes("admin"), controllers.OAuthClientRegisterHandler)
		clientGroup.GET("", m.OAuth(), m.RequiredScopes("admin"), controllers.ListOAuthClientsHandler)
		clientGroup.GET("/:client_id", m.OAuth(), m.RequiredScopes("admin"), controllers.GetOAuthClientHandler)
		clientGroup.PUT("/:client_id", m.OAuth(), m.RequiredScopes("admin"), controllers.UpdateOAuthClientHandler)
		clientGroup.DELETE("/:client_id", m.OAuth(), m.RequiredScopes("admin"), controllers.DeleteOAuthClientHandler)
	}

}
