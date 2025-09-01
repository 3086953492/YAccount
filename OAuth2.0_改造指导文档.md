# YAccount OAuth2.0 改造指导文档

## 概述

本文档指导如何将 YAccount 项目改造为标准的 OAuth2.0 授权服务器，使其能够为其他应用提供统一的用户认证和授权服务。改造完成后，YAccount 将作为中央认证服务，其他新开发的服务只需接入 OAuth2.0 标准接口即可完成用户登录。

## 目标架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   第三方应用A    │    │   第三方应用B    │    │   第三方应用C    │
│   (Client App)  │    │   (Client App)  │    │   (Client App)  │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                         ┌───────▼───────┐
                         │   YAccount    │
                         │ OAuth2.0 服务 │
                         │ (Authorization│
                         │ + Resource    │
                         │   Server)     │
                         └───────────────┘
```

## 改造计划

### 阶段一：数据模型扩展 (1-2天)
- 新增 OAuth 客户端管理模型
- 新增授权码管理模型
- 新增访问令牌管理模型
- 新增权限范围(Scope)管理模型

### 阶段二：配置和基础设施 (0.5-1天)
- 扩展配置结构支持 OAuth
- 添加 OAuth 相关错误类型
- 扩展 JWT 工具支持多种令牌类型

### 阶段三：OAuth 核心服务 (2-3天)
- 实现客户端管理服务
- 实现授权码流程服务
- 实现令牌管理服务
- 实现权限验证服务

### 阶段四：OAuth 标准端点 (2-3天)
- 实现授权端点 `/oauth/authorize`
- 实现令牌端点 `/oauth/token`
- 实现用户信息端点 `/oauth/userinfo`
- 实现令牌内省端点 `/oauth/introspect`

### 阶段五：授权页面和中间件 (1-2天)
- 实现授权确认页面
- 实现 OAuth 中间件
- 更新现有中间件兼容 OAuth

### 阶段六：客户端管理功能 (1天)
- 实现客户端注册功能
- 实现客户端管理后台

---

## 详细实施步骤

## 阶段一：数据模型扩展

### 1.1 新增 OAuth 客户端模型

在 `server/models/` 目录下创建 `oauth_client.go`:

```go
package models

import (
	"time"
	"gorm.io/gorm"
)

// OAuth客户端表
type OAuthClient struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 客户端基本信息
	ClientID     string `gorm:"unique;not null;size:255" json:"client_id"`
	ClientSecret string `gorm:"not null;size:255" json:"-"` // 不在 JSON 中返回
	Name         string `gorm:"not null;size:255" json:"name"`
	Description  string `gorm:"size:500" json:"description"`
	
	// 重定向URI
	RedirectURIs string `gorm:"type:text" json:"redirect_uris"` // JSON 数组存储
	
	// 授权类型
	GrantTypes string `gorm:"not null;size:255;default:'authorization_code'" json:"grant_types"` // 逗号分隔
	
	// 权限范围
	Scopes string `gorm:"size:255;default:'read'" json:"scopes"` // 逗号分隔
	
	// 客户端类型: public/confidential
	ClientType string `gorm:"not null;size:50;default:'confidential'" json:"client_type"`
	
	// 状态: active/inactive
	Status string `gorm:"not null;size:50;default:'active'" json:"status"`
	
	// 令牌有效期(秒)
	AccessTokenTTL  int `gorm:"default:3600" json:"access_token_ttl"`   // 1小时
	RefreshTokenTTL int `gorm:"default:604800" json:"refresh_token_ttl"` // 7天
	
	// 所有者信息
	OwnerID   uint   `gorm:"not null" json:"owner_id"`
	OwnerType string `gorm:"not null;size:50;default:'user'" json:"owner_type"` // user/system
}

func (OAuthClient) TableName() string {
	return "oauth_clients"
}

// 创建客户端请求
type CreateOAuthClientRequest struct {
	Name         string   `json:"name" validate:"required,max=255"`
	Description  string   `json:"description" validate:"max=500"`
	RedirectURIs []string `json:"redirect_uris" validate:"required,min=1"`
	GrantTypes   []string `json:"grant_types" validate:"required"`
	Scopes       []string `json:"scopes" validate:"required"`
	ClientType   string   `json:"client_type" validate:"required,oneof=public confidential"`
}
```

### 1.2 新增授权码模型

继续在 `server/models/oauth_client.go` 中添加:

```go
// OAuth授权码表
type OAuthAuthorizationCode struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	Code        string `gorm:"unique;not null;size:255" json:"code"`
	ClientID    string `gorm:"not null;size:255;index" json:"client_id"`
	UserID      uint   `gorm:"not null;index" json:"user_id"`
	RedirectURI string `gorm:"not null;size:500" json:"redirect_uri"`
	Scopes      string `gorm:"size:255" json:"scopes"`
	ExpiresAt   time.Time `gorm:"not null" json:"expires_at"`
	Used        bool   `gorm:"default:false" json:"used"`
	
	// PKCE 支持
	CodeChallenge       string `gorm:"size:255" json:"code_challenge"`
	CodeChallengeMethod string `gorm:"size:50" json:"code_challenge_method"`
}

func (OAuthAuthorizationCode) TableName() string {
	return "oauth_authorization_codes"
}
```

### 1.3 新增访问令牌模型

继续添加:

```go
// OAuth访问令牌表
type OAuthAccessToken struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	AccessToken  string `gorm:"unique;not null;size:255" json:"access_token"`
	RefreshToken string `gorm:"unique;size:255" json:"refresh_token"`
	ClientID     string `gorm:"not null;size:255;index" json:"client_id"`
	UserID       uint   `gorm:"not null;index" json:"user_id"`
	Scopes       string `gorm:"size:255" json:"scopes"`
	ExpiresAt    time.Time `gorm:"not null" json:"expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
	Revoked      bool   `gorm:"default:false" json:"revoked"`
}

func (OAuthAccessToken) TableName() string {
	return "oauth_access_tokens"
}
```

### 1.4 新增权限范围模型

继续添加:

```go
// OAuth权限范围表
type OAuthScope struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	Name        string `gorm:"unique;not null;size:100" json:"name"`
	Description string `gorm:"size:500" json:"description"`
	IsDefault   bool   `gorm:"default:false" json:"is_default"`
	Status      string `gorm:"not null;size:50;default:'active'" json:"status"`
}

func (OAuthScope) TableName() string {
	return "oauth_scopes"
}
```

### 1.5 数据库迁移

在 `server/initialize/database.go` 中添加新表的自动迁移:

```go
// 在现有的 AutoMigrate 中添加新模型
func AutoMigrate() error {
	return global.DB.AutoMigrate(
		&models.User{},
		&models.SystemInfo{},
		// 新增 OAuth 相关表
		&models.OAuthClient{},
		&models.OAuthAuthorizationCode{},
		&models.OAuthAccessToken{},
		&models.OAuthScope{},
	)
}
```

## 阶段二：配置和基础设施

### 2.1 扩展配置结构

在 `server/configs/config.go` 中添加 OAuth 配置:

```go
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Log        LogConfig        `mapstructure:"log"`
	Middleware MiddlewareConfig `mapstructure:"middleware"`
	OAuth      OAuthConfig      `mapstructure:"oauth"` // 新增
}
```

创建 `server/configs/oauth.go`:

```go
package configs

import "time"

type OAuthConfig struct {
	// 授权服务器基本配置
	Issuer      string `mapstructure:"issuer"`
	AuthorizeUI string `mapstructure:"authorize_ui"` // 授权页面URL
	
	// 授权码配置
	AuthorizationCodeTTL time.Duration `mapstructure:"authorization_code_ttl"`
	
	// 访问令牌配置
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
	
	// 支持的授权类型
	SupportedGrantTypes []string `mapstructure:"supported_grant_types"`
	
	// 支持的响应类型
	SupportedResponseTypes []string `mapstructure:"supported_response_types"`
	
	// 默认权限范围
	DefaultScopes []string `mapstructure:"default_scopes"`
	
	// PKCE 支持
	RequirePKCE           bool     `mapstructure:"require_pkce"`
	SupportedChallengeMethod []string `mapstructure:"supported_challenge_methods"`
}
```

### 2.2 更新配置示例文件

更新 `server/configs/config.example.json`:

```json
{
  "server": {
    "port": 8081,
    "mode": "debug",
    "id": 1
  },
  "Database": {
    "host": "localhost",
    "port": 3306,
    "user": "root",
    "password": "123456",
    "dbname": "yaccount",
    "charset": "utf8mb4",
    "parseTime": true,
    "loc": "Local"
  },
  "redis": {
    "host": "localhost",
    "port": 6379,
    "password": "",
    "db": 0
  },
  "jwt": {
    "secret": "examplekey",
    "expire": "1h"
  },
  "oauth": {
    "issuer": "http://localhost:8081",
    "authorize_ui": "http://localhost:5173/oauth/authorize",
    "authorization_code_ttl": "10m",
    "access_token_ttl": "1h",
    "refresh_token_ttl": "168h",
    "supported_grant_types": [
      "authorization_code",
      "refresh_token",
      "client_credentials"
    ],
    "supported_response_types": [
      "code"
    ],
    "default_scopes": [
      "read"
    ],
    "require_pkce": false,
    "supported_challenge_methods": [
      "S256",
      "plain"
    ]
  },
  "log": {
    "level": "info",
    "filename": "logs/example.log",
    "maxSize": 10,
    "maxBackups": 3,
    "maxAge": 7,
    "compress": true,
    "rotateDaily": true
  },
  "middleware": {
    "jwt": {
      "secret": "examplekey",
      "expire": "1h",
      "skip_paths": [
        "/api/account/v1/auth/login",
        "/api/account/v1/auth/register",
        "/oauth/authorize",
        "/oauth/token",
        "/oauth/userinfo"
      ]
    },
    "cors": {
      "allow_origins": [
        "*"
      ],
      "allow_methods": [
        "GET",
        "POST",
        "PUT",
        "DELETE",
        "OPTIONS"
      ],
      "allow_headers": [
        "Authorization",
        "Content-Type"
      ]
    }
  }
}
```

### 2.3 扩展错误类型

在 `server/pkg/apperrors/` 目录下创建 `errors_oauth.go`:

```go
package apperrors

import "net/http"

// OAuth 相关错误
var (
	// 客户端错误
	ErrInvalidClient = &AppError{
		Code:       "INVALID_CLIENT",
		Message:    "客户端认证失败",
		HTTPStatus: http.StatusUnauthorized,
	}
	
	ErrUnsupportedGrantType = &AppError{
		Code:       "UNSUPPORTED_GRANT_TYPE",
		Message:    "不支持的授权类型",
		HTTPStatus: http.StatusBadRequest,
	}
	
	ErrInvalidGrant = &AppError{
		Code:       "INVALID_GRANT",
		Message:    "授权码无效或已过期",
		HTTPStatus: http.StatusBadRequest,
	}
	
	ErrInvalidScope = &AppError{
		Code:       "INVALID_SCOPE",
		Message:    "请求的权限范围无效",
		HTTPStatus: http.StatusBadRequest,
	}
	
	ErrInvalidRequest = &AppError{
		Code:       "INVALID_REQUEST",
		Message:    "请求参数无效",
		HTTPStatus: http.StatusBadRequest,
	}
	
	ErrInvalidRedirectURI = &AppError{
		Code:       "INVALID_REDIRECT_URI",
		Message:    "重定向URI无效",
		HTTPStatus: http.StatusBadRequest,
	}
	
	ErrAccessDenied = &AppError{
		Code:       "ACCESS_DENIED",
		Message:    "用户拒绝授权",
		HTTPStatus: http.StatusForbidden,
	}
	
	ErrUnsupportedResponseType = &AppError{
		Code:       "UNSUPPORTED_RESPONSE_TYPE",
		Message:    "不支持的响应类型",
		HTTPStatus: http.StatusBadRequest,
	}
	
	// 令牌相关错误
	ErrInvalidToken = &AppError{
		Code:       "INVALID_TOKEN",
		Message:    "访问令牌无效",
		HTTPStatus: http.StatusUnauthorized,
	}
	
	ErrTokenExpired = &AppError{
		Code:       "TOKEN_EXPIRED",
		Message:    "访问令牌已过期",
		HTTPStatus: http.StatusUnauthorized,
	}
	
	ErrInsufficientScope = &AppError{
		Code:       "INSUFFICIENT_SCOPE",
		Message:    "权限范围不足",
		HTTPStatus: http.StatusForbidden,
	}
)
```

### 2.4 扩展 JWT 工具

创建 `server/pkg/oauth/token.go`:

```go
package oauth

import (
	"YAccount/global"
	"YAccount/models"
	"YAccount/pkg/apperrors"
	"YAccount/utils/logger"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type OAuthClaims struct {
	UserID   uint     `json:"sub"`
	ClientID string   `json:"client_id"`
	Scopes   []string `json:"scope"`
	TokenType string  `json:"token_type"` // access_token, refresh_token
	jwt.RegisteredClaims
}

// 生成授权码
func GenerateAuthorizationCode() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// 生成访问令牌
func GenerateAccessToken(userID uint, clientID string, scopes []string) (string, error) {
	claims := OAuthClaims{
		UserID:    userID,
		ClientID:  clientID,
		Scopes:    scopes,
		TokenType: "access_token",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(global.Cfg.OAuth.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    global.Cfg.OAuth.Issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(global.Cfg.JWT.Secret))
	if err != nil {
		logger.LogError("GenerateAccessToken", "jwt.NewWithClaims", "生成访问令牌失败", err, zap.Any("claims", claims))
		return "", apperrors.ErrTokenGenerateFailed
	}
	return tokenString, nil
}

// 生成刷新令牌
func GenerateRefreshToken(userID uint, clientID string, scopes []string) (string, error) {
	claims := OAuthClaims{
		UserID:    userID,
		ClientID:  clientID,
		Scopes:    scopes,
		TokenType: "refresh_token",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(global.Cfg.OAuth.RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    global.Cfg.OAuth.Issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(global.Cfg.JWT.Secret))
	if err != nil {
		logger.LogError("GenerateRefreshToken", "jwt.NewWithClaims", "生成刷新令牌失败", err, zap.Any("claims", claims))
		return "", apperrors.ErrTokenGenerateFailed
	}
	return tokenString, nil
}

// 解析 OAuth 令牌
func ParseOAuthToken(tokenString string) (*OAuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &OAuthClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(global.Cfg.JWT.Secret), nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return nil, apperrors.ErrTokenExpired
		}
		return nil, apperrors.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*OAuthClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, apperrors.ErrInvalidToken
}

// 验证权限范围
func ValidateScopes(requestedScopes, allowedScopes []string) bool {
	for _, requested := range requestedScopes {
		found := false
		for _, allowed := range allowedScopes {
			if requested == allowed {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
```

## 阶段三：OAuth 核心服务

### 3.1 创建 OAuth 客户端服务

在 `server/services/` 目录下创建 `oauth_client.go`:

```go
package services

import (
	"YAccount/global"
	"YAccount/models"
	"YAccount/pkg/apperrors"
	"YAccount/repositories"
	"YAccount/utils/logger"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// 创建 OAuth 客户端
func CreateOAuthClient(req *models.CreateOAuthClientRequest, ownerID uint) (*models.OAuthClient, error) {
	// 生成客户端ID和密钥
	clientID, err := generateClientID()
	if err != nil {
		return nil, apperrors.ErrServerInternal
	}

	clientSecret, err := generateClientSecret()
	if err != nil {
		return nil, apperrors.ErrServerInternal
	}

	// 加密客户端密钥
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.ErrServerInternal
	}

	// 序列化数组字段
	redirectURIsJSON, _ := json.Marshal(req.RedirectURIs)
	grantTypesStr := strings.Join(req.GrantTypes, ",")
	scopesStr := strings.Join(req.Scopes, ",")

	client := &models.OAuthClient{
		ClientID:     clientID,
		ClientSecret: string(hashedSecret),
		Name:         req.Name,
		Description:  req.Description,
		RedirectURIs: string(redirectURIsJSON),
		GrantTypes:   grantTypesStr,
		Scopes:       scopesStr,
		ClientType:   req.ClientType,
		OwnerID:      ownerID,
		OwnerType:    "user",
	}

	// 保存到数据库
	if err := global.DB.Create(client).Error; err != nil {
		logger.LogError("CreateOAuthClient", "database", "创建OAuth客户端失败", err, zap.Any("client", client))
		return nil, apperrors.ErrServerInternal
	}

	// 在响应中包含未加密的客户端密钥（仅此一次）
	client.ClientSecret = clientSecret

	logger.Info("OAuth客户端创建成功", zap.String("clientID", clientID), zap.String("name", req.Name))
	return client, nil
}

// 验证客户端凭证
func ValidateClientCredentials(clientID, clientSecret string) (*models.OAuthClient, error) {
	var client models.OAuthClient
	if err := global.DB.Where("client_id = ? AND status = ?", clientID, "active").First(&client).Error; err != nil {
		return nil, apperrors.ErrInvalidClient
	}

	// 验证客户端密钥
	if err := bcrypt.CompareHashAndPassword([]byte(client.ClientSecret), []byte(clientSecret)); err != nil {
		return nil, apperrors.ErrInvalidClient
	}

	return &client, nil
}

// 验证重定向URI
func ValidateRedirectURI(client *models.OAuthClient, redirectURI string) bool {
	var redirectURIs []string
	if err := json.Unmarshal([]byte(client.RedirectURIs), &redirectURIs); err != nil {
		return false
	}

	for _, uri := range redirectURIs {
		if uri == redirectURI {
			return true
		}
	}
	return false
}

// 生成客户端ID
func generateClientID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("client_%s", base64.URLEncoding.EncodeToString(bytes)[:22]), nil
}

// 生成客户端密钥
func generateClientSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// 获取客户端信息
func GetOAuthClientByID(clientID string) (*models.OAuthClient, error) {
	var client models.OAuthClient
	if err := global.DB.Where("client_id = ? AND status = ?", clientID, "active").First(&client).Error; err != nil {
		return nil, apperrors.ErrInvalidClient
	}
	return &client, nil
}
```

### 3.2 创建 OAuth 授权服务

创建 `server/services/oauth_auth.go`:

```go
package services

import (
	"YAccount/global"
	"YAccount/models"
	"YAccount/pkg/apperrors"
	"YAccount/pkg/oauth"
	"YAccount/utils/logger"
	"encoding/json"
	"strings"
	"time"

	"go.uber.org/zap"
)

// 授权请求参数
type AuthorizeRequest struct {
	ResponseType    string `form:"response_type" binding:"required"`
	ClientID        string `form:"client_id" binding:"required"`
	RedirectURI     string `form:"redirect_uri" binding:"required"`
	Scope           string `form:"scope"`
	State           string `form:"state"`
	CodeChallenge   string `form:"code_challenge"`
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
		requestedScopes = global.Cfg.OAuth.DefaultScopes
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
		ExpiresAt:           time.Now().Add(global.Cfg.OAuth.AuthorizationCodeTTL),
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: req.CodeChallengeMethod,
	}

	if err := global.DB.Create(authCode).Error; err != nil {
		logger.LogError("HandleAuthorizeRequest", "database", "保存授权码失败", err, zap.Any("authCode", authCode))
		return "", apperrors.ErrServerInternal
	}

	// 构建重定向URL
	redirectURL := fmt.Sprintf("%s?code=%s", req.RedirectURI, code)
	if req.State != "" {
		redirectURL += "&state=" + req.State
	}

	logger.Info("授权码生成成功", zap.String("clientID", req.ClientID), zap.Uint("userID", userID))
	return redirectURL, nil
}

// 处理令牌请求
func HandleTokenRequest(req *TokenRequest) (map[string]interface{}, error) {
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
func handleAuthorizationCodeGrant(req *TokenRequest) (map[string]interface{}, error) {
	// 验证客户端
	client, err := ValidateClientCredentials(req.ClientID, req.ClientSecret)
	if err != nil {
		return nil, err
	}

	// 查找授权码
	var authCode models.OAuthAuthorizationCode
	if err := global.DB.Where("code = ? AND client_id = ? AND used = ?", req.Code, req.ClientID, false).First(&authCode).Error; err != nil {
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
	if err := global.DB.Model(&authCode).Update("used", true).Error; err != nil {
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
		ExpiresAt:        time.Now().Add(global.Cfg.OAuth.AccessTokenTTL),
		RefreshExpiresAt: time.Now().Add(global.Cfg.OAuth.RefreshTokenTTL),
	}

	if err := global.DB.Create(tokenRecord).Error; err != nil {
		logger.LogError("handleAuthorizationCodeGrant", "database", "保存令牌失败", err)
		return nil, apperrors.ErrServerInternal
	}

	return map[string]interface{}{
		"access_token":  accessToken,
		"token_type":    "Bearer",
		"expires_in":    int(global.Cfg.OAuth.AccessTokenTTL.Seconds()),
		"refresh_token": refreshToken,
		"scope":         authCode.Scopes,
	}, nil
}

// 处理刷新令牌模式
func handleRefreshTokenGrant(req *TokenRequest) (map[string]interface{}, error) {
	// 验证客户端
	_, err := ValidateClientCredentials(req.ClientID, req.ClientSecret)
	if err != nil {
		return nil, err
	}

	// 查找刷新令牌
	var tokenRecord models.OAuthAccessToken
	if err := global.DB.Where("refresh_token = ? AND client_id = ? AND revoked = ?", req.RefreshToken, req.ClientID, false).First(&tokenRecord).Error; err != nil {
		return nil, apperrors.ErrInvalidGrant
	}

	// 检查刷新令牌是否过期
	if time.Now().After(tokenRecord.RefreshExpiresAt) {
		return nil, apperrors.ErrInvalidGrant
	}

	// 撤销旧令牌
	if err := global.DB.Model(&tokenRecord).Update("revoked", true).Error; err != nil {
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
		ExpiresAt:        time.Now().Add(global.Cfg.OAuth.AccessTokenTTL),
		RefreshExpiresAt: time.Now().Add(global.Cfg.OAuth.RefreshTokenTTL),
	}

	if err := global.DB.Create(newTokenRecord).Error; err != nil {
		logger.LogError("handleRefreshTokenGrant", "database", "保存新令牌失败", err)
		return nil, apperrors.ErrServerInternal
	}

	return map[string]interface{}{
		"access_token":  accessToken,
		"token_type":    "Bearer",
		"expires_in":    int(global.Cfg.OAuth.AccessTokenTTL.Seconds()),
		"refresh_token": refreshToken,
		"scope":         tokenRecord.Scopes,
	}, nil
}

// 处理客户端凭证模式
func handleClientCredentialsGrant(req *TokenRequest) (map[string]interface{}, error) {
	// 验证客户端
	client, err := ValidateClientCredentials(req.ClientID, req.ClientSecret)
	if err != nil {
		return nil, err
	}

	// 验证权限范围
	requestedScopes := strings.Fields(req.Scope)
	if len(requestedScopes) == 0 {
		requestedScopes = global.Cfg.OAuth.DefaultScopes
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
		ExpiresAt:   time.Now().Add(global.Cfg.OAuth.AccessTokenTTL),
	}

	if err := global.DB.Create(tokenRecord).Error; err != nil {
		logger.LogError("handleClientCredentialsGrant", "database", "保存令牌失败", err)
		return nil, apperrors.ErrServerInternal
	}

	return map[string]interface{}{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   int(global.Cfg.OAuth.AccessTokenTTL.Seconds()),
		"scope":        strings.Join(requestedScopes, " "),
	}, nil
}
```

## 阶段四：OAuth 标准端点

### 4.1 创建 OAuth 控制器

在 `server/controllers/` 目录下创建 `oauth.go`:

```go
package controllers

import (
	"YAccount/models"
	"YAccount/pkg/apperrors"
	"YAccount/pkg/oauth"
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
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, apperrors.ErrInvalidRequest)
		return
	}

	// 检查用户是否已登录
	userID := c.GetUint("user_id")
	if userID == 0 {
		// 重定向到登录页面，并在登录后回调到授权页面
		c.Redirect(http.StatusFound, "/login?redirect="+c.Request.URL.String())
		return
	}

	// 处理授权请求
	redirectURL, err := services.HandleAuthorizeRequest(&req, userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 重定向到客户端
	c.Redirect(http.StatusFound, redirectURL)
}

// OAuth 令牌端点
func OAuthTokenHandler(c *gin.Context) {
	var req services.TokenRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, apperrors.ErrInvalidRequest)
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

// OAuth 用户信息端点
func OAuthUserInfoHandler(c *gin.Context) {
	// 从 Authorization 头中获取访问令牌
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Error(c, apperrors.ErrInvalidToken)
		return
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		response.Error(c, apperrors.ErrInvalidToken)
		return
	}

	accessToken := tokenParts[1]

	// 解析和验证令牌
	claims, err := oauth.ParseOAuthToken(accessToken)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 检查令牌类型
	if claims.TokenType != "access_token" {
		response.Error(c, apperrors.ErrInvalidToken)
		return
	}

	// 获取用户信息
	userInfo, err := services.GetUserProfile(claims.UserID)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 返回用户信息（根据权限范围过滤）
	userInfoResponse := map[string]interface{}{
		"sub":      claims.UserID,
		"username": userInfo.Username,
		"nickname": userInfo.Nickname,
	}

	// 根据权限范围添加额外信息
	for _, scope := range claims.Scopes {
		switch scope {
		case "profile":
			userInfoResponse["avatar"] = userInfo.Avatar
			userInfoResponse["role"] = userInfo.Role
		case "email":
			// 如果有邮箱字段，在这里添加
		}
	}

	response.Success(c, "获取用户信息成功", userInfoResponse)
}

// OAuth 令牌内省端点
func OAuthIntrospectHandler(c *gin.Context) {
	token := c.PostForm("token")
	if token == "" {
		response.Error(c, apperrors.ErrInvalidRequest)
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
		"active":    true,
		"sub":       claims.UserID,
		"client_id": claims.ClientID,
		"scope":     strings.Join(claims.Scopes, " "),
		"exp":       claims.ExpiresAt.Unix(),
		"iat":       claims.IssuedAt.Unix(),
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
```

### 4.2 创建 OAuth 路由

在 `server/routers/` 目录下创建 `oauth.go`:

```go
package routers

import (
	"YAccount/controllers"
	"YAccount/middleware"

	"github.com/gin-gonic/gin"
)

func LoadOAuthRouters(router *gin.Engine) {
	m := middleware.NewManager()

	// OAuth 标准端点
	oauthGroup := router.Group("/oauth")
	{
		// 授权端点（需要用户登录）
		oauthGroup.GET("/authorize", m.Auth(), controllers.OAuthAuthorizeHandler)
		
		// 令牌端点（公开）
		oauthGroup.POST("/token", controllers.OAuthTokenHandler)
		
		// 用户信息端点（需要有效的访问令牌，由专门的中间件验证）
		oauthGroup.GET("/userinfo", controllers.OAuthUserInfoHandler)
		
		// 令牌内省端点（公开）
		oauthGroup.POST("/introspect", controllers.OAuthIntrospectHandler)
	}

	// OAuth 客户端管理端点
	clientGroup := router.Group("/api/oauth/v1/clients")
	clientGroup.Use(m.Auth()) // 需要登录
	{
		// 注册新客户端
		clientGroup.POST("", controllers.OAuthClientRegisterHandler)
		
		// 获取客户端列表（管理员）
		clientGroup.GET("", m.AdminPermission(), controllers.ListOAuthClientsHandler)
		
		// 获取客户端详情
		clientGroup.GET("/:client_id", controllers.GetOAuthClientHandler)
		
		// 更新客户端
		clientGroup.PUT("/:client_id", controllers.UpdateOAuthClientHandler)
		
		// 删除客户端
		clientGroup.DELETE("/:client_id", controllers.DeleteOAuthClientHandler)
	}
}
```

### 4.3 更新路由初始化

在 `server/initialize/router.go` 中添加 OAuth 路由:

```go
func InitRouters() *gin.Engine {
	router := gin.Default()

	middlewareManager := middleware.NewManager()
	middlewareManager.LoadGlobal(router)

	// 注册路由
	routers.LoadUserRouters(router)
	routers.LoadAuthRouters(router)
	routers.LoadSystemInfoRouters(router)
	routers.LoadOAuthRouters(router) // 新增
	
	return router
}
```

## 阶段五：OAuth 中间件

### 5.1 创建 OAuth 中间件

在 `server/middleware/oauth/` 目录下创建 `oauth.go`:

```go
package oauth

import (
	"YAccount/pkg/apperrors"
	"YAccount/pkg/oauth"
	"YAccount/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// OAuth 访问令牌验证中间件
func OAuthTokenMiddleware(requiredScopes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization 头中获取令牌
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, apperrors.ErrInvalidToken)
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			response.Error(c, apperrors.ErrInvalidToken)
			c.Abort()
			return
		}

		accessToken := tokenParts[1]

		// 解析和验证令牌
		claims, err := oauth.ParseOAuthToken(accessToken)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}

		// 检查令牌类型
		if claims.TokenType != "access_token" {
			response.Error(c, apperrors.ErrInvalidToken)
			c.Abort()
			return
		}

		// 检查权限范围
		if len(requiredScopes) > 0 && !oauth.ValidateScopes(requiredScopes, claims.Scopes) {
			response.Error(c, apperrors.ErrInsufficientScope)
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("client_id", claims.ClientID)
		c.Set("scopes", claims.Scopes)
		c.Set("token_type", "oauth")

		c.Next()
	}
}
```

### 5.2 更新中间件管理器

在 `server/middleware/manager.go` 中添加 OAuth 中间件:

```go
// 导入
import (
	"YAccount/middleware/oauth"
)

// 获取 OAuth 中间件
func (m *Manager) OAuth(requiredScopes ...string) gin.HandlerFunc {
	return oauth.OAuthTokenMiddleware(requiredScopes...)
}
```

## 阶段六：完善功能

### 6.1 添加客户端管理控制器方法

在 `server/controllers/oauth.go` 中添加缺失的控制器方法:

```go
// 获取客户端列表
func ListOAuthClientsHandler(c *gin.Context) {
	clients, err := services.ListOAuthClients()
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

// 更新客户端
func UpdateOAuthClientHandler(c *gin.Context) {
	// 实现客户端更新逻辑
	response.Success(c, "客户端更新成功", nil)
}

// 删除客户端
func DeleteOAuthClientHandler(c *gin.Context) {
	// 实现客户端删除逻辑
	response.Success(c, "客户端删除成功", nil)
}
```

### 6.2 添加对应的服务方法

在 `server/services/oauth_client.go` 中添加:

```go
// 获取客户端列表
func ListOAuthClients() ([]models.OAuthClient, error) {
	var clients []models.OAuthClient
	if err := global.DB.Where("status = ?", "active").Find(&clients).Error; err != nil {
		return nil, apperrors.ErrServerInternal
	}
	return clients, nil
}

// 获取客户端详情
func GetOAuthClientDetail(clientID string, userID uint, role string) (*models.OAuthClient, error) {
	var client models.OAuthClient
	query := global.DB.Where("client_id = ? AND status = ?", clientID, "active")
	
	// 非管理员只能查看自己的客户端
	if role != "admin" {
		query = query.Where("owner_id = ?", userID)
	}
	
	if err := query.First(&client).Error; err != nil {
		return nil, apperrors.ErrInvalidClient
	}
	return &client, nil
}
```

### 6.3 初始化默认权限范围

创建数据库初始化脚本，在 `server/initialize/database.go` 中添加:

```go
// 初始化默认数据
func InitDefaultData() error {
	// 初始化默认权限范围
	defaultScopes := []models.OAuthScope{
		{Name: "read", Description: "读取基本用户信息", IsDefault: true},
		{Name: "write", Description: "修改用户信息", IsDefault: false},
		{Name: "profile", Description: "访问完整用户资料", IsDefault: false},
		{Name: "admin", Description: "管理员权限", IsDefault: false},
	}

	for _, scope := range defaultScopes {
		var existingScope models.OAuthScope
		if err := global.DB.Where("name = ?", scope.Name).First(&existingScope).Error; err != nil {
			if err := global.DB.Create(&scope).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
```

在 `server/main.go` 中调用初始化:

```go
// 在数据库初始化后添加
if err := initialize.InitDefaultData(); err != nil {
	logger.Error("初始化默认数据失败", zap.Error(err))
	return
}
```

## 测试和验证

### 测试流程

1. **授权码模式测试**:
   ```bash
   # 步骤1: 获取授权码
   curl "http://localhost:8081/oauth/authorize?response_type=code&client_id=your_client_id&redirect_uri=http://localhost:3000/callback&scope=read&state=abc123"
   
   # 步骤2: 使用授权码获取访问令牌
   curl -X POST http://localhost:8081/oauth/token \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "grant_type=authorization_code&client_id=your_client_id&client_secret=your_client_secret&code=authorization_code_here&redirect_uri=http://localhost:3000/callback"
   ```

2. **客户端凭证模式测试**:
   ```bash
   curl -X POST http://localhost:8081/oauth/token \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "grant_type=client_credentials&client_id=your_client_id&client_secret=your_client_secret&scope=read"
   ```

3. **访问受保护资源**:
   ```bash
   curl -X GET http://localhost:8081/oauth/userinfo \
     -H "Authorization: Bearer your_access_token"
   ```

### 数据库表结构

执行以下 SQL 确认表结构正确创建:

```sql
-- 查看 OAuth 相关表
SHOW TABLES LIKE 'oauth_%';

-- 查看表结构
DESCRIBE oauth_clients;
DESCRIBE oauth_authorization_codes;
DESCRIBE oauth_access_tokens;
DESCRIBE oauth_scopes;
```

## 部署注意事项

1. **安全配置**:
   - 生产环境中确保使用强密钥
   - 启用 HTTPS
   - 配置适当的 CORS 策略

2. **性能优化**:
   - 为常用查询添加数据库索引
   - 启用 Redis 缓存
   - 考虑令牌清理任务

3. **监控**:
   - 添加日志记录
   - 监控令牌使用情况
   - 设置异常告警

## 总结

通过以上步骤，YAccount 项目将被成功改造为标准的 OAuth2.0 授权服务器。改造后的系统将支持:

- ✅ 标准的 OAuth2.0 授权码模式
- ✅ 客户端凭证模式  
- ✅ 刷新令牌机制
- ✅ 权限范围管理
- ✅ 客户端管理
- ✅ 令牌内省
- ✅ 用户信息端点

新开发的服务只需要实现 OAuth2.0 客户端逻辑，即可快速接入统一的认证系统，大大简化了多服务架构下的用户认证复杂度。
