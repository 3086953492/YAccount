package repositories

import (
	"YAccount/global"
	"YAccount/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(req *models.RegisterRequest) (*models.User, error) {

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Status:   1,
		Nickname: req.Nickname,
	}

	if err := global.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func Login(req *models.LoginRequest) (*models.User, error) {
	var user models.User
	if err := global.DB.Where("username = ? AND status = 1", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := global.DB.Where("id = ? AND status = 1", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(id uint, req *models.UpdateUserRequest) (*models.User, error) {
	var user models.User
	if err := global.DB.Where("id = ? AND status = 1", id).First(&user).Error; err != nil {
		return nil, err
	}

	// 存储用户表的更新
	userUpdates := make(map[string]interface{})

	// 非空字段加入更新
	if req.Username != "" {
		userUpdates["username"] = req.Username
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		userUpdates["password"] = string(hashedPassword)
	}
	if req.Nickname != "" {
		userUpdates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		userUpdates["avatar"] = req.Avatar
	}

	// 更新用户表
	if len(userUpdates) > 0 {
		if err := global.DB.Model(&user).Updates(userUpdates).Error; err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func DeleteUser(id uint) error {
	return global.DB.Delete(&models.User{}, id).Error
}

func GetUsers(page, pageSize int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	offset := (page - 1) * pageSize

	if err := global.DB.Model(&models.User{}).Where("status = 1").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := global.DB.Where("status = 1").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := global.DB.Where("username =? AND status = 1", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserList(page, pageSize int, query string) ([]models.User, int64, error) {
	var users []models.User
	if err := global.DB.Where("status = 1 " + query).
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := global.DB.Model(&models.User{}).Where("status = 1 " + query).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}