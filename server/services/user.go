package services

import (
	"YAccount/global"
	"YAccount/models"
	apperrors "github.com/3086953492/YaBase/errors"
	"YAccount/pkg/oauth"
	"YAccount/repositories"
	"github.com/3086953492/YaBase/logger"
	"YAccount/utils/redis"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/cache/v9"
	"go.uber.org/zap"
)

func LoginService(req *models.LoginRequest, clientID string) (*models.UserResponse, *models.TokenResponse, error) {

	_, err := GetOAuthClientByID(clientID)
	if err != nil {
		return nil, nil, err
	}

	user, err := repositories.Login(req)
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
		ExpiresAt:        time.Now().Add(global.Cfg.OAuth.AccessTokenTTL),
		RefreshExpiresAt: time.Now().Add(global.Cfg.OAuth.RefreshTokenTTL),
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
		ExpiresIn:    int(global.Cfg.OAuth.AccessTokenTTL.Seconds()),
		RefreshToken: refreshToken,
		Scope:        strings.Join(scopes, " "),
	}

	return userResponse, tokenResponse, nil
}

func RegisterService(req *models.RegisterRequest) (*models.User, error) {

	// 对用户名加锁
	lockKey := fmt.Sprintf("user:register:%s", req.Username)
	lock := redis.NewDistributedLock(lockKey, 10*time.Second)

	if err := lock.Acquire(); err != nil {
		return nil, apperrors.ErrUsernameExists
	}
	defer lock.Release()

	user, err := repositories.Register(req)
	if err != nil {
		return nil, apperrors.ErrUserRegisterFailed
	}

	logger.Info("用户注册成功", zap.String("username", user.Username))

	return user, nil
}

func UpdateService(req *models.UpdateUserRequest, userID uint) error {

	if _, err := repositories.UpdateUser(userID, req); err != nil {
		logger.LogError("UpdateService", "database query", "更新用户信息失败", err, zap.Uint("userID", userID))
		return apperrors.ErrUpdateFailed
	}

	global.Cache.Delete(context.Background(), fmt.Sprintf("user:profile:%d", userID))

	logger.Info("用户信息更新成功", zap.Uint("userID", userID))

	return nil
}

func ListUsersPage(query string, page, pageSize int) (models.PaginationResponse[models.UserList], error) {
	var paginationResponse models.PaginationResponse[models.UserList]
	if err := global.Cache.Once(&cache.Item{
		Key:   fmt.Sprintf("users:list:%s:%d:%d", query, page, pageSize),
		Value: &paginationResponse,
		Do: func(*cache.Item) (any, error) {
			users, total, err := repositories.GetUserList(page, pageSize, query)
			if err != nil {
				if !apperrors.IsNotFoundError(err) {
					logger.LogError("ListUsersPage", "database query", "从数据库中获取用户列表失败", err, zap.String("query", query), zap.Int("page", page), zap.Int("pageSize", pageSize))
				}
				return nil, err
			}
			userList := make([]models.UserList, len(users))
			for i, user := range users {
				userList[i] = models.UserList{
					ID:       user.ID,
					Username: user.Username,
					Nickname: user.Nickname,
					Avatar:   user.Avatar,
				}
			}
			return models.PaginationResponse[models.UserList]{
				Items:      userList,
				Total:      total,
				Page:       page,
				PageSize:   pageSize,
				TotalPages: int(total / int64(pageSize)),
			}, nil
		},
	}); err != nil {
		return models.PaginationResponse[models.UserList]{}, apperrors.ErrUserListNotFound
	}
	return paginationResponse, nil
}

func GetUserProfile(userID uint) (*models.UserResponse, error) {
	var user models.UserResponse
	if err := global.Cache.Once(&cache.Item{
		Key:   fmt.Sprintf("user:profile:%d", userID),
		Value: &user,
		Do: func(*cache.Item) (any, error) {
			user, err := repositories.GetUserByID(userID)
			if err != nil {
				if !apperrors.IsNotFoundError(err) {
					logger.LogError("GetUserProfile", "database query", "从数据库中获取用户信息失败", err, zap.Uint("userID", userID))
				}
				return nil, err
			}
			return &models.UserResponse{
				ID:        user.ID,
				Username:  user.Username,
				Role:      user.Role,
				Nickname:  user.Nickname,
				Avatar:    user.Avatar,
				Status:    user.Status,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			}, nil
		},
	}); err != nil {
		return nil, apperrors.ErrUserNotFound
	}
	return &user, nil
}
