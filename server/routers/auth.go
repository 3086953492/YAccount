package routers

import (
	"YAccount/controllers"
	"github.com/3086953492/YaBase/middleware"
	"github.com/gin-gonic/gin"
)

func LoadAuthRouters(router *gin.Engine) {
	m := middleware.NewManager()
	authRouters := router.Group("/api/account/v1/auth")
	{
		authRouters.POST("/login", controllers.LoginHandler)
		authRouters.POST("/register", controllers.RegisterHandler)
		authRouters.PUT("/token", m.OAuth(), controllers.OAuthRefreshHandler)
	}
}
