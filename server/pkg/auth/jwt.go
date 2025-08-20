package auth

import (
	"YAccount/global"
	"YAccount/pkg/apperrors"
	"YAccount/utils/logger"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, username string, role string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(global.Cfg.JWT.Expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(global.Cfg.JWT.Secret))
	if err != nil {
		logger.LogError("GenerateToken", "jwt.NewWithClaims", "生成token失败", err, zap.Any("claims", claims))
		return "", apperrors.ErrTokenGenerateFailed
	}
	return tokenString, nil
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(global.Cfg.JWT.Secret), nil
	})

	if err != nil {
		logger.LogError("ParseToken", "jwt.ParseWithClaims", "解析token失败", err, zap.String("token", tokenString))
		return nil, apperrors.ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, apperrors.ErrTokenInvalid
}

func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 检查token是否即将过期（还有10分钟过期）
	if time.Until(claims.ExpiresAt.Time) < 10*time.Minute {
		return GenerateToken(claims.UserID, claims.Username, claims.Role)
	}

	return tokenString, apperrors.ErrTokenRefreshTooEarly
}
