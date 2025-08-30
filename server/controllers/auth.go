package controllers

import (
	"YAccount/pkg/apperrors"
	"YAccount/pkg/auth"
	"YAccount/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func RefreshHandler(c *gin.Context) {
	// 从请求头中获取旧令牌
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Error(c, apperrors.ErrUnauthorized)
		return
	}

	// 解析Bearer token
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		response.Error(c, apperrors.ErrTokenInvalid)
		return
	}

	token := tokenParts[1]

	// 调用刷新令牌的方法
	newToken, err := auth.RefreshToken(token)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 返回新令牌
	response.Success(c, "刷新令牌成功", gin.H{
		"token": newToken,
	})
}
