package routers

import (
	"YAccount/controllers"
	"YAccount/middleware"

	"github.com/gin-gonic/gin"
)

func LoadUserRouters(router *gin.Engine) {
	m := middleware.NewManager()
	userRouters := router.Group("/api/account/v1/users")
	{
		userRouters.PUT("/:user_id", m.Auth(), controllers.UpdateHandler)
		userRouters.GET("", m.Auth(), m.AdminPermission(), controllers.UserListHandler)
		userRouters.GET("/:user_id", m.Auth(), controllers.UserProfileHandler)
	}
}
