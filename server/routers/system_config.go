package routers

import (
	"YAccount/controllers"

	"github.com/gin-gonic/gin"
)

func LoadSystemConfigRouters(r *gin.Engine) {
	systemConfigRouters := r.Group("api/account/v1/system")
	{
		systemConfigRouters.GET("/configs", controllers.SystemConfigListHandler)
	}
}
