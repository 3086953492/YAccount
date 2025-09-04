package controllers

import (
	"YAccount/models"
	"github.com/3086953492/YaBase/response"
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

func BatchUpdateSystemInfoHandler(c *gin.Context) {
	var req models.BatchUpdateSystemInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, err)
		return
	}
	if err := services.BatchUpdateSystemInfo(&req); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "批量更新系统信息成功", nil)
}