package routers

import (
	"YAccount/controllers"
	"YAccount/middleware"
	"github.com/gin-gonic/gin"
)

func LoadSystemInfoRouters(r *gin.Engine) {
	m := middleware.NewManager()
	systemInfoRouters := r.Group("api/account/v1/system")
	{
		systemInfoRouters.GET("/infos", controllers.SystemInfoListHandler)
		systemInfoRouters.GET("/infos/:key", controllers.SystemInfoByKeyHandler)
		systemInfoRouters.POST("/infos", m.OAuth(), m.RequiredScopes("admin"), controllers.BatchUpdateSystemInfoHandler)
	}
}
