package controllers

import (
	"YAccount/pkg/response"
	"YAccount/services"

	"github.com/gin-gonic/gin"
)

func SystemInfoListHandler(c *gin.Context) {
	systemInfoList, err := services.GetSystemInfoList()
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "获取系统配置列表成功", systemInfoList)
}
