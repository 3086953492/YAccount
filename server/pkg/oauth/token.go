package oauth

import (
	"YAccount/global"
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
	UserID    uint     `json:"sub"`
	ClientID  string   `json:"client_id"`
	Scopes    []string `json:"scope"`
	TokenType string   `json:"token_type"` // access_token, refresh_token
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
