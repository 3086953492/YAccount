package auth

import (
	"YAccount/global"
	"YAccount/pkg/apperrors"
	"YAccount/utils/logger"
	"errors"
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
		// 检查是否是token过期错误
		if errors.Is(err, jwt.ErrTokenExpired) {
			logger.LogError("ParseToken", "jwt.ParseWithClaims", "token已过期", err, zap.String("token", tokenString))
			return nil, apperrors.ErrTokenExpired
		}

		// 检查是否是token无效错误
		if errors.Is(err, jwt.ErrTokenMalformed) {
			logger.LogError("ParseToken", "jwt.ParseWithClaims", "token格式错误", err, zap.String("token", tokenString))
			return nil, apperrors.ErrTokenInvalid
		}

		// 检查是否是token未生效错误
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			logger.LogError("ParseToken", "jwt.ParseWithClaims", "token还未生效", err, zap.String("token", tokenString))
			return nil, apperrors.ErrTokenInvalid
		}

		// 检查是否是签名无效错误
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			logger.LogError("ParseToken", "jwt.ParseWithClaims", "token签名无效", err, zap.String("token", tokenString))
			return nil, apperrors.ErrTokenInvalid
		}

		// 其他未知错误
		logger.LogError("ParseToken", "jwt.ParseWithClaims", "解析token失败", err, zap.String("token", tokenString))
		return nil, apperrors.ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, apperrors.ErrTokenInvalid
}

func RefreshToken(tokenString string) (string, error) {
	// 尝试解析token，即使是过期的token也要能够解析出用户信息用于刷新
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(global.Cfg.JWT.Secret), nil
	})

	var claims *Claims
	if err != nil {
		// 只有token过期错误才允许刷新
		if errors.Is(err, jwt.ErrTokenExpired) {
			if c, ok := token.Claims.(*Claims); ok {
				claims = c
			} else {
				logger.LogError("RefreshToken", "jwt.ParseWithClaims", "无法从过期token中提取claims", err, zap.String("token", tokenString))
				return "", apperrors.ErrTokenInvalid
			}
		} else {
			// 其他错误不允许刷新
			logger.LogError("RefreshToken", "jwt.ParseWithClaims", "token无效，无法刷新", err, zap.String("token", tokenString))
			return "", apperrors.ErrTokenInvalid
		}
	} else {
		// token有效的情况
		if c, ok := token.Claims.(*Claims); ok {
			claims = c
		} else {
			return "", apperrors.ErrTokenInvalid
		}
	}

	// 生成新的token
	return GenerateToken(claims.UserID, claims.Username, claims.Role)
}
