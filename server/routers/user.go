package routers

import (
	"YAccount/controllers"
	"github.com/3086953492/YaBase/middleware"
	"github.com/gin-gonic/gin"
)

func LoadUserRouters(router *gin.Engine) {
	m := middleware.NewManager()
	userRouters := router.Group("/api/account/v1/users")
	{
		userRouters.PUT("/:user_id", m.OAuth(), controllers.UpdateHandler)
		userRouters.GET("", m.OAuth(),m.RequiredScopes("admin"), controllers.UserListHandler)
		userRouters.GET("/:user_id", m.OAuth(), controllers.UserProfileHandler)
	}
}
