package controllers

import (
	"YAccount/models"
	"YAccount/pkg/apperrors"
	"YAccount/pkg/pagination"
	"YAccount/pkg/response"
	"YAccount/pkg/validator"
	"YAccount/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {

	var req models.LoginRequest

	if !validator.ValidateStruct(c, &req) {
		response.Error(c, apperrors.ErrInvalidInput)
		return
	}

	user, token, err := services.LoginService(&req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, "登录成功", gin.H{
		"user":  user,
		"token": token,
	})
}

func RegisterHandler(c *gin.Context) {

	var req models.RegisterRequest
	if !validator.ValidateStruct(c, &req) {
		response.Error(c, apperrors.ErrInvalidInput)
		return
	}

	user, err := services.RegisterService(&req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, "注册成功", gin.H{
		"user": user,
	})
}

func UpdateHandler(c *gin.Context) {

	// 从请求体中获取用户id
	userIDStr := c.Param("user_id")

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil || userID < 1 {
		response.Error(c, apperrors.ErrInvalidInput)
		return
	}

	currentUserID := c.GetUint("user_id") // 从上下文中获取用户ID
	role := c.GetString("role")

	// 判断当前用户是否在修改自己的信息，只有管理员可修改他人信息
	if uint(userID) != currentUserID && role != "admin" {
		response.Error(c, apperrors.ErrPermissionDenied)
		return
	}

	// 解析请求参数
	var req models.UpdateUserRequest
	if !validator.ValidateStruct(c, &req) {
		response.Error(c, apperrors.ErrInvalidInput)
		return
	}

	if role != "admin" { // 防止用户将自己设置为管理员
		req.Role = "user"
	}

	err = services.UpdateService(&req, uint(userID))
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, "用户信息更新成功", nil)
}

func UserListHandler(c *gin.Context) {

	var query string

	username := c.Query("username")
	email := c.Query("email")
	nickname := c.Query("nickname")

	if username != "" {
		query += " AND username LIKE '%" + username + "%'"
	}
	if email != "" {
		query += " AND email LIKE '%" + email + "%'"
	}
	if nickname != "" {
		query += " AND nickname LIKE '%" + nickname + "%'"
	}

	page, pageSize := pagination.GetPageAndPageSize(c)

	users, err := services.ListUsersPage(query, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "获取用户列表成功", users)
}
