package services

import (
	"YAccount/models"
	"YAccount/pkg/oauth"
	"YAccount/repositories"
	"fmt"
	"strings"
	"time"

	"github.com/3086953492/YaBase/config"
	"github.com/3086953492/YaBase/config/types"
	apperrors "github.com/3086953492/YaBase/errors"
	"github.com/3086953492/YaBase/logger"

	"go.uber.org/zap"
)

// 授权请求参数
type AuthorizeRequest struct {
	ResponseType        string `form:"response_type" binding:"required"`
	ClientID            string `form:"client_id" binding:"required"`
	RedirectURI         string `form:"redirect_uri" binding:"required"`
	Scope               string `form:"scope"`
	State               string `form:"state"`
	CodeChallenge       string `form:"code_challenge"`
	CodeChallengeMethod string `form:"code_challenge_method"`
}

// 令牌请求参数
type TokenRequest struct {
	GrantType    string `form:"grant_type" binding:"required"`
	ClientID     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
	Code         string `form:"code"`
	RedirectURI  string `form:"redirect_uri"`
	RefreshToken string `form:"refresh_token"`
	Scope        string `form:"scope"`
	CodeVerifier string `form:"code_verifier"`
}

func oauthCfg() *types.OAuthConfig {
	return &config.GetGlobalConfig().OAuth
}

// 处理授权请求
func HandleAuthorizeRequest(req *AuthorizeRequest, userID uint) (string, error) {
	// 验证响应类型
	if req.ResponseType != "code" {
		return "", apperrors.ErrUnsupportedResponseType
	}

	// 验证客户端
	client, err := GetOAuthClientByID(req.ClientID)
	if err != nil {
		return "", err
	}

	// 验证重定向URI
	if !ValidateRedirectURI(client, req.RedirectURI) {
		return "", apperrors.ErrInvalidRedirectURI
	}

	// 验证权限范围
	requestedScopes := strings.Fields(req.Scope)
	if len(requestedScopes) == 0 {
		requestedScopes = oauthCfg().DefaultScopes
	}

	clientScopes := strings.Split(client.Scopes, ",")
	if !oauth.ValidateScopes(requestedScopes, clientScopes) {
		return "", apperrors.ErrInvalidScope
	}

	// 生成授权码
	code, err := oauth.GenerateAuthorizationCode()
	if err != nil {
		return "", apperrors.ErrServerInternal
	}

	// 保存授权码
	authCode := &models.OAuthAuthorizationCode{
		Code:                code,
		ClientID:            req.ClientID,
		UserID:              userID,
		RedirectURI:         req.RedirectURI,
		Scopes:              strings.Join(requestedScopes, " "),
		ExpiresAt:           time.Now().Add(oauthCfg().AuthorizationCodeTTL),
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: req.CodeChallengeMethod,
	}

	if err := repositories.CreateOAuthAuthorizationCode(authCode); err != nil {
		logger.LogError("HandleAuthorizeRequest", "database", "保存授权码失败", err, zap.Any("authCode", authCode))
		return "", apperrors.ErrServerInternal
	}

	// 构建重定向URL
	redirectURL := fmt.Sprintf("%s?code=%s", req.RedirectURI, code)
	if req.State != "" {
		redirectURL += "&state=" + req.State
	}

	return redirectURL, nil
}

// 处理令牌请求
func HandleTokenRequest(req *TokenRequest) (map[string]any, error) {
	switch req.GrantType {
	case "authorization_code":
		return handleAuthorizationCodeGrant(req)
	case "refresh_token":
		return handleRefreshTokenGrant(req)
	case "client_credentials":
		return handleClientCredentialsGrant(req)
	default:
		return nil, apperrors.ErrUnsupportedGrantType
	}
}

// 处理授权码模式
func handleAuthorizationCodeGrant(req *TokenRequest) (map[string]any, error) {
	// 验证客户端
	_, err := ValidateClientCredentials(req.ClientID, req.ClientSecret)
	if err != nil {
		return nil, err
	}

	// 查找授权码
	var authCode models.OAuthAuthorizationCode
	if err := repositories.GetOAuthAuthorizationCodeByCode(req.Code, req.ClientID, false); err != nil {
		if !apperrors.IsNotFoundError(err) {
			logger.LogError("handleAuthorizationCodeGrant", "database", "获取授权码失败", err, zap.String("code", req.Code), zap.String("clientID", req.ClientID))
		}
		return nil, apperrors.ErrInvalidGrant
	}

	// 检查授权码是否过期
	if time.Now().After(authCode.ExpiresAt) {
		return nil, apperrors.ErrInvalidGrant
	}

	// 验证重定向URI
	if authCode.RedirectURI != req.RedirectURI {
		return nil, apperrors.ErrInvalidGrant
	}

	// 标记授权码为已使用
	if err := repositories.UpdateOAuthAuthorizationCodeUsed(authCode.ID, true); err != nil {
		logger.LogError("handleAuthorizationCodeGrant", "database", "更新授权码状态失败", err)
		return nil, apperrors.ErrServerInternal
	}

	// 生成访问令牌和刷新令牌
	scopes := strings.Fields(authCode.Scopes)
	accessToken, err := oauth.GenerateAccessToken(authCode.UserID, authCode.ClientID, scopes)
	if err != nil {
		return nil, err
	}

	refreshToken, err := oauth.GenerateRefreshToken(authCode.UserID, authCode.ClientID, scopes)
	if err != nil {
		return nil, err
	}

	// 保存令牌到数据库
	tokenRecord := &models.OAuthAccessToken{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ClientID:         authCode.ClientID,
		UserID:           authCode.UserID,
		Scopes:           authCode.Scopes,
		ExpiresAt:        time.Now().Add(oauthCfg().AccessTokenTTL),
		RefreshExpiresAt: time.Now().Add(oauthCfg().RefreshTokenTTL),
	}

	if err := repositories.CreateOAuthAccessToken(tokenRecord); err != nil {
		logger.LogError("handleAuthorizationCodeGrant", "database", "保存令牌失败", err)
		return nil, apperrors.ErrServerInternal
	}

	return map[string]any{
		"access_token":  accessToken,
		"token_type":    "Bearer",
		"expires_in":    int(oauthCfg().AccessTokenTTL.Seconds()),
		"refresh_token": refreshToken,
		"scope":         authCode.Scopes,
	}, nil
}

// 处理刷新令牌模式
func handleRefreshTokenGrant(req *TokenRequest) (map[string]any, error) {
	// 验证客户端
	_, err := ValidateClientCredentials(req.ClientID, req.ClientSecret)
	if err != nil {
		return nil, err
	}

	// 查找刷新令牌
	var tokenRecord models.OAuthAccessToken
	if err := repositories.GetOAuthAccessTokenByRefreshToken(req.RefreshToken, req.ClientID, false); err != nil {
		if !apperrors.IsNotFoundError(err) {
			logger.LogError("handleRefreshTokenGrant", "database", "获取刷新令牌失败", err, zap.String("refreshToken", req.RefreshToken), zap.String("clientID", req.ClientID))
		}
		return nil, apperrors.ErrInvalidGrant
	}

	// 检查刷新令牌是否过期
	if time.Now().After(tokenRecord.RefreshExpiresAt) {
		return nil, apperrors.ErrInvalidGrant
	}

	// 撤销旧令牌
	if err := repositories.UpdateOAuthAccessTokenRevoked(tokenRecord.ID, true); err != nil {
		logger.LogError("handleRefreshTokenGrant", "database", "撤销旧令牌失败", err)
		return nil, apperrors.ErrServerInternal
	}

	// 生成新的访问令牌和刷新令牌
	scopes := strings.Fields(tokenRecord.Scopes)
	accessToken, err := oauth.GenerateAccessToken(tokenRecord.UserID, tokenRecord.ClientID, scopes)
	if err != nil {
		return nil, err
	}

	refreshToken, err := oauth.GenerateRefreshToken(tokenRecord.UserID, tokenRecord.ClientID, scopes)
	if err != nil {
		return nil, err
	}

	// 保存新令牌
	newTokenRecord := &models.OAuthAccessToken{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ClientID:         tokenRecord.ClientID,
		UserID:           tokenRecord.UserID,
		Scopes:           tokenRecord.Scopes,
		ExpiresAt:        time.Now().Add(oauthCfg().AccessTokenTTL),
		RefreshExpiresAt: time.Now().Add(oauthCfg().RefreshTokenTTL),
	}

	if err := repositories.CreateOAuthAccessToken(newTokenRecord); err != nil {
		logger.LogError("handleRefreshTokenGrant", "database", "保存新令牌失败", err)
		return nil, apperrors.ErrServerInternal
	}

	return map[string]any{
		"access_token":  accessToken,
		"token_type":    "Bearer",
		"expires_in":    int(oauthCfg().AccessTokenTTL.Seconds()),
		"refresh_token": refreshToken,
		"scope":         tokenRecord.Scopes,
	}, nil
}

// 处理客户端凭证模式
func handleClientCredentialsGrant(req *TokenRequest) (map[string]any, error) {
	// 验证客户端
	client, err := ValidateClientCredentials(req.ClientID, req.ClientSecret)
	if err != nil {
		return nil, err
	}

	// 验证权限范围
	requestedScopes := strings.Fields(req.Scope)
	if len(requestedScopes) == 0 {
		requestedScopes = oauthCfg().DefaultScopes
	}

	clientScopes := strings.Split(client.Scopes, ",")
	if !oauth.ValidateScopes(requestedScopes, clientScopes) {
		return nil, apperrors.ErrInvalidScope
	}

	// 生成访问令牌（客户端凭证模式不需要用户ID）
	accessToken, err := oauth.GenerateAccessToken(0, client.ClientID, requestedScopes)
	if err != nil {
		return nil, err
	}

	// 保存令牌（客户端凭证模式不需要刷新令牌）
	tokenRecord := &models.OAuthAccessToken{
		AccessToken: accessToken,
		ClientID:    client.ClientID,
		UserID:      0, // 客户端凭证模式没有用户
		Scopes:      strings.Join(requestedScopes, " "),
		ExpiresAt:   time.Now().Add(oauthCfg().AccessTokenTTL),
	}

	if err := repositories.CreateOAuthAccessToken(tokenRecord); err != nil {
		logger.LogError("handleClientCredentialsGrant", "database", "保存令牌失败", err)
		return nil, apperrors.ErrServerInternal
	}

	return map[string]any{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   int(oauthCfg().AccessTokenTTL.Seconds()),
		"scope":        strings.Join(requestedScopes, " "),
	}, nil
}

func LoginService(req *models.LoginRequest, clientID string) (*models.UserResponse, *models.TokenResponse, error) {

	_, err := GetOAuthClientByID(clientID)
	if err != nil {
		return nil, nil, err
	}

	user, err := repositories.VerifyUserPassword(req)
	if err != nil {
		if !apperrors.IsNotFoundError(err) {
			logger.LogError("LoginService", "database query", "从数据库中获取用户失败", err, zap.String("username", req.Username))
		}
		return nil, nil, apperrors.ErrUsernameOrPasswordError
	}

	// 根据用户角色确定授权范围
	scopes := []string{"read"}
	if user.Role == "admin" {
		scopes = append(scopes, "write", "admin")
	}

	// 生成访问令牌和刷新令牌
	accessToken, err := oauth.GenerateAccessToken(user.ID, clientID, scopes)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := oauth.GenerateRefreshToken(user.ID, clientID, scopes)
	if err != nil {
		return nil, nil, err
	}

	// 保存令牌记录
	tokenRecord := &models.OAuthAccessToken{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ClientID:         clientID,
		UserID:           user.ID,
		Scopes:           strings.Join(scopes, " "),
		ExpiresAt:        time.Now().Add(oauthCfg().AccessTokenTTL),
		RefreshExpiresAt: time.Now().Add(oauthCfg().RefreshTokenTTL),
	}

	if err := repositories.CreateOAuthAccessToken(tokenRecord); err != nil {
		logger.LogError("OAuthLoginService", "database", "保存令牌失败", err)
		return nil, nil, apperrors.ErrServerInternal
	}

	logger.Info("用户登录成功", zap.String("username", user.Username))

	userResponse := &models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	tokenResponse := &models.TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int(oauthCfg().AccessTokenTTL.Seconds()),
		RefreshToken: refreshToken,
		Scope:        strings.Join(scopes, " "),
	}

	return userResponse, tokenResponse, nil
}
