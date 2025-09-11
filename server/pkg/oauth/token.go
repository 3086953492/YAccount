package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/3086953492/YaBase/config"
	apperrors "github.com/3086953492/YaBase/errors"
	"github.com/3086953492/YaBase/logger"

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

func cfg() *config.Config {
	return config.GetGlobalConfig()
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg().OAuth.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    cfg().OAuth.Issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg().JWT.Secret))
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg().OAuth.RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    cfg().OAuth.Issuer,
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg().JWT.Secret))
	if err != nil {
		logger.LogError("GenerateRefreshToken", "jwt.NewWithClaims", "生成刷新令牌失败", err, zap.Any("claims", claims))
		return "", apperrors.ErrTokenGenerateFailed
	}
	return tokenString, nil
}

// 解析 OAuth 令牌
func ParseOAuthToken(tokenString string) (*OAuthClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &OAuthClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(cfg().JWT.Secret), nil
	})

	if err != nil {
		// 检查是否是token过期错误
		if errors.Is(err, jwt.ErrTokenExpired) {
			logger.LogError("ParseOAuthToken", "jwt.ParseWithClaims", "OAuth token已过期", err, zap.String("token", tokenString))
			return nil, apperrors.ErrTokenExpired
		}

		// 检查是否是token无效错误
		if errors.Is(err, jwt.ErrTokenMalformed) {
			logger.LogError("ParseOAuthToken", "jwt.ParseWithClaims", "OAuth token格式错误", err, zap.String("token", tokenString))
			return nil, apperrors.ErrInvalidToken
		}

		// 检查是否是token未生效错误
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			logger.LogError("ParseOAuthToken", "jwt.ParseWithClaims", "OAuth token还未生效", err, zap.String("token", tokenString))
			return nil, apperrors.ErrInvalidToken
		}

		// 检查是否是签名无效错误
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			logger.LogError("ParseOAuthToken", "jwt.ParseWithClaims", "OAuth token签名无效", err, zap.String("token", tokenString))
			return nil, apperrors.ErrInvalidToken
		}

		// 其他未知错误
		logger.LogError("ParseOAuthToken", "jwt.ParseWithClaims", "解析OAuth token失败", err, zap.String("token", tokenString))
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
