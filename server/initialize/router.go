package initialize

import (
	"YAccount/middleware"
	"YAccount/routers"

	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	router := gin.Default()

	middlewareManager := middleware.NewManager()
	middlewareManager.LoadGlobal(router)

	// 注册路由
	routers.LoadUserRouters(router)
	routers.LoadAuthRouters(router)
	
	return router
}
