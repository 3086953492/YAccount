package initialize

import (
	"YAccount/routers"
	"github.com/3086953492/YaBase/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	router := gin.Default()

	middlewareManager := middleware.NewManager()
	middlewareManager.LoadGlobal(router)

	// 注册路由
	routers.LoadUserRouters(router)
	routers.LoadAuthRouters(router)
	routers.LoadSystemInfoRouters(router)
	routers.LoadOAuthRouters(router)
	return router
}
