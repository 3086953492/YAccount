package services

import (
	"YAccount/global"
	"YAccount/models"
	"YAccount/pkg/apperrors"
	"YAccount/pkg/auth"
	"YAccount/repositories"
	"YAccount/utils/logger"
	"YAccount/utils/redis"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/cache/v9"
	"go.uber.org/zap"
)

func LoginService(req *models.LoginRequest) (*models.UserResponse, string, error) {

	user, err := repositories.Login(req)
	if err != nil {
		if apperrors.IsNotFoundError(err) {
			logger.LogError("LoginService", "database query", "从数据库中获取用户失败", err, zap.String("username", req.Username))
		}
		return nil, "", apperrors.ErrUsernameOrPasswordError
	}

	// 生成JWT令牌
	token, err := auth.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, "", err
	}

	logger.Info("用户登录成功", zap.String("username", user.Username))

	return &models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Status:   user.Status,
	}, token, nil
}

func RegisterService(req *models.RegisterRequest) (*models.User, error) {

	// 对用户名和邮箱加锁
	lockKey := fmt.Sprintf("user:register:%s:%s", req.Username, req.Email)
	lock := redis.NewDistributedLock(lockKey, 10*time.Second)

	if err := lock.Acquire(); err != nil {
		return nil, apperrors.ErrUsernameOrEmailExists
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

func ListUsersPage(query string, page, pageSize int) (models.PaginationResponse, error) {
	var paginationResponse models.PaginationResponse
	if err := global.Cache.Once(&cache.Item{
		Key:   fmt.Sprintf("users:list:%s:%d:%d", query, page, pageSize),
		Value: &paginationResponse,
		Do: func(*cache.Item) (any, error) {
			users, total, err := repositories.GetUserList(page, pageSize, query)
			if err != nil {
				if apperrors.IsNotFoundError(err) {
					logger.LogError("ListUsersPage", "database query", "从数据库中获取用户列表失败", err, zap.String("query", query), zap.Int("page", page), zap.Int("pageSize", pageSize))
				}
				return nil, err
			}
			userList := make([]models.UserList, len(users))
			for i, user := range users {
				userList[i] = models.UserList{
					ID:       user.ID,
					Username: user.Username,
					Email:    user.Email,
					Nickname: user.Nickname,
					Avatar:   user.Avatar,
				}	
			}
			return models.PaginationResponse{
				Items: userList,
				Total: total,
				Page:  page,
				PageSize: pageSize,
				TotalPages: int(total / int64(pageSize)),
			}, nil
		},
	}); err != nil {
		return models.PaginationResponse{}, apperrors.ErrUserListNotFound
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
				if apperrors.IsNotFoundError(err) {
					logger.LogError("GetUserProfile", "database query", "从数据库中获取用户信息失败", err, zap.Uint("userID", userID))
				}
				return nil, err
			}
			return &models.UserResponse{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
				Role:     user.Role,
				Nickname: user.Nickname,
				Avatar:   user.Avatar,
				Status:   user.Status,
			}, nil
		},
	}); err != nil {
		return nil, apperrors.ErrUserNotFound
	}
	return &user, nil
}