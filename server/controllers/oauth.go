package controllers

import (
	"YAccount/models"
	"YAccount/pkg/apperrors"
	"YAccount/pkg/oauth"
	"YAccount/pkg/pagination"
	"YAccount/pkg/response"
	"YAccount/pkg/validator"
	"YAccount/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// OAuth 授权端点
func OAuthAuthorizeHandler(c *gin.Context) {
	var req services.AuthorizeRequest
	if !validator.ValidateStruct(c, &req) {
		response.Error(c, apperrors.ErrInvalidInput)
		return
	}

	// 检查用户是否已登录
	userID := c.GetUint("user_id")
	if userID == 0 {
		// 返回需要登录的响应，前端根据此响应跳转到登录页面
		response.Error(c, apperrors.ErrUnauthorized)
		return
	}

	// 处理授权请求
	redirectURL, err := services.HandleAuthorizeRequest(&req, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 返回授权成功响应，包含重定向URL，前端负责处理重定向
	response.Success(c, "授权成功", map[string]interface{}{
		"redirect_uri": redirectURL,
		"state":        req.State,
	})
}

// OAuth 令牌端点
func OAuthTokenHandler(c *gin.Context) {
	var req services.TokenRequest
	if !validator.ValidateStruct(c, &req) {
		response.Error(c, apperrors.ErrInvalidInput)
		return
	}

	// 处理令牌请求
	tokenResponse, err := services.HandleTokenRequest(&req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, "令牌获取成功", tokenResponse)
}

// OAuth 令牌内省端点
func OAuthIntrospectHandler(c *gin.Context) {
	token := c.PostForm("token")
	if token == "" {
		response.Error(c, apperrors.ErrInvalidInput)
		return
	}

	// 解析令牌
	claims, err := oauth.ParseOAuthToken(token)
	if err != nil {
		// 令牌无效时返回 active: false
		c.JSON(http.StatusOK, map[string]interface{}{
			"active": false,
		})
		return
	}

	// 返回令牌信息
	introspectResponse := map[string]interface{}{
		"active":     true,
		"sub":        claims.UserID,
		"client_id":  claims.ClientID,
		"scope":      strings.Join(claims.Scopes, " "),
		"exp":        claims.ExpiresAt.Unix(),
		"iat":        claims.IssuedAt.Unix(),
		"token_type": claims.TokenType,
	}

	c.JSON(http.StatusOK, introspectResponse)
}

// 客户端注册端点
func OAuthClientRegisterHandler(c *gin.Context) {
	var req models.CreateOAuthClientRequest
	if !validator.ValidateStruct(c, &req) {
		response.Error(c, apperrors.ErrInvalidInput)
		return
	}

	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Error(c, apperrors.ErrUnauthorized)
		return
	}

	client, err := services.CreateOAuthClient(&req, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, "OAuth客户端创建成功", client)
}

// OAuth 授权确认端点（用于前后端分离架构）
func OAuthAuthorizeConfirmHandler(c *gin.Context) {
	var req models.OAuthAuthorizeConfirmRequest

	if !validator.ValidateStruct(c, &req) {
		response.Error(c, apperrors.ErrInvalidInput)
		return
	}

	userID := c.GetUint("user_id")
	if userID == 0 {
		response.Error(c, apperrors.ErrUnauthorized)
		return
	}

	if !req.Approved {
		// 用户拒绝授权，返回错误响应
		response.Error(c, apperrors.ErrAccessDenied)
		return
	}

	// 构建授权请求
	authorizeReq := services.AuthorizeRequest{
		ResponseType:        "code",
		ClientID:            req.ClientID,
		RedirectURI:         req.RedirectURI,
		Scope:               req.Scope,
		State:               req.State,
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: req.CodeChallengeMethod,
	}

	// 处理授权请求
	redirectURL, err := services.HandleAuthorizeRequest(&authorizeReq, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 返回授权成功响应
	response.Success(c, "授权成功", map[string]interface{}{
		"redirect_uri": redirectURL,
		"state":        req.State,
	})
}

// 获取客户端列表
func ListOAuthClientsHandler(c *gin.Context) {
	page, pageSize := pagination.GetPageAndPageSize(c)
	var query string
	name := c.Query("name")
	if name != "" {
		query += " AND name LIKE '%" + name + "%'"
	}
	clients, err := services.ListOAuthClients(page, pageSize, query)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "获取客户端列表成功", clients)
}

// 获取客户端详情
func GetOAuthClientHandler(c *gin.Context) {
	clientID := c.Param("client_id")
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	client, err := services.GetOAuthClientDetail(clientID, userID, role)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "获取客户端详情成功", client)
}

func UpdateOAuthClientHandler(c *gin.Context) {
	clientID := c.Param("client_id")
	var req models.UpdateOAuthClientRequest
	if !validator.ValidateStruct(c, &req) {
		response.Error(c, apperrors.ErrInvalidInput)
		return
	}
	err := services.UpdateOAuthClient(clientID, req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "更新客户端成功", nil)
}

func DeleteOAuthClientHandler(c *gin.Context) {
	clientID := c.Param("client_id")
	err := services.DeleteOAuthClient(clientID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, "删除客户端成功", nil)
}
