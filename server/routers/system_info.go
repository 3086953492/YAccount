package routers

import (
	"YAccount/controllers"

	"github.com/gin-gonic/gin"
)

func LoadSystemInfoRouters(r *gin.Engine) {
	systemInfoRouters := r.Group("api/account/v1/system")
	{
		systemInfoRouters.GET("/infos", controllers.SystemInfoListHandler)
		systemInfoRouters.GET("/infos/:key", controllers.SystemInfoByKeyHandler)
	}
}
