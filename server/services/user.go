package services

import (
	"YAccount/models"
	"YAccount/repositories"
	"context"
	"fmt"
	"time"

	"github.com/3086953492/YaBase/cache"
	apperrors "github.com/3086953492/YaBase/errors"
	"github.com/3086953492/YaBase/logger"
	"github.com/3086953492/YaBase/redis"

	"go.uber.org/zap"
)

// userCache 优雅地获取用户缓存实例
func userCache() *cache.Cache {
	return cache.GetGlobalCache()
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

	userCache().Delete(context.Background(), fmt.Sprintf("user:profile:%d", userID))

	logger.Info("用户信息更新成功", zap.Uint("userID", userID))

	return nil
}

func ListUsersPage(query string, page, pageSize int) (models.PaginationResponse[models.UserList], error) {
	var paginationResponse models.PaginationResponse[models.UserList]
	if err := userCache().Once(&cache.Item{
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
	if err := userCache().Once(&cache.Item{
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
