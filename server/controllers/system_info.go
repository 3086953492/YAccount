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
	response.Success(c, "获取系统信息列表成功", systemInfoList)
}

func SystemInfoByKeyHandler(c *gin.Context) {
	key := c.Param("key")
	systemInfo, err := services.GetSystemInfoByKey(key)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "获取系统信息成功", systemInfo)
}