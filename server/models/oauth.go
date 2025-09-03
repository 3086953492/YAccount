package models

import (
	"gorm.io/gorm"
	"time"
)

// OAuth客户端表
type OAuthClient struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 客户端基本信息
	ClientID     string `gorm:"unique;not null;size:255" json:"client_id"`
	ClientSecret string `gorm:"not null;size:255" json:"client_secret"`
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
	AccessTokenTTL  int `gorm:"default:3600" json:"access_token_ttl"`    // 1小时
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

// OAuth授权码表
type OAuthAuthorizationCode struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Code        string    `gorm:"unique;not null;size:255" json:"code"`
	ClientID    string    `gorm:"not null;size:255;index" json:"client_id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	RedirectURI string    `gorm:"not null;size:500" json:"redirect_uri"`
	Scopes      string    `gorm:"size:255" json:"scopes"`
	ExpiresAt   time.Time `gorm:"not null" json:"expires_at"`
	Used        bool      `gorm:"default:false" json:"used"`

	// PKCE 支持
	CodeChallenge       string `gorm:"size:255" json:"code_challenge"`
	CodeChallengeMethod string `gorm:"size:50" json:"code_challenge_method"`
}

func (OAuthAuthorizationCode) TableName() string {
	return "oauth_authorization_codes"
}

// OAuth访问令牌表
type OAuthAccessToken struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	AccessToken      string    `gorm:"not null;type:text" json:"access_token"`
	RefreshToken     string    `gorm:"type:text" json:"refresh_token"`
	ClientID         string    `gorm:"not null;size:255;index" json:"client_id"`
	UserID           uint      `gorm:"not null;index" json:"user_id"`
	Scopes           string    `gorm:"size:255" json:"scopes"`
	ExpiresAt        time.Time `gorm:"not null" json:"expires_at"`
	RefreshExpiresAt time.Time `json:"refresh_expires_at"`
	Revoked          bool      `gorm:"default:false" json:"revoked"`
}

func (OAuthAccessToken) TableName() string {
	return "oauth_access_tokens"
}

// OAuth权限范围表
type OAuthScope struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"unique;not null;size:100" json:"name"`
	Description string `gorm:"size:500" json:"description"`
	IsDefault   bool   `gorm:"default:false" json:"is_default"`
	Status      string `gorm:"not null;size:50;default:'active'" json:"status"`
}

func (OAuthScope) TableName() string {
	return "oauth_scopes"
}

type OAuthAuthorizeConfirmRequest struct {
	ClientID            string `json:"client_id" binding:"required"`
	RedirectURI         string `json:"redirect_uri" binding:"required"`
	Scope               string `json:"scope"`
	State               string `json:"state"`
	CodeChallenge       string `json:"code_challenge"`
	CodeChallengeMethod string `json:"code_challenge_method"`
	Approved            bool   `json:"approved" binding:"required"`
}

type UpdateOAuthClientRequest struct {
	Name         string   `json:"name" validate:"required,max=255"`
	Description  string   `json:"description" validate:"max=500"`
	RedirectURIs []string `json:"redirect_uris" validate:"required,min=1"`
	GrantTypes   []string `json:"grant_types" validate:"required"`
	Scopes       []string `json:"scopes" validate:"required"`
	ClientType   string   `json:"client_type" validate:"required,oneof=public confidential"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}
