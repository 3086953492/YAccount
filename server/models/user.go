package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"unique;not null;size:50" json:"username"`
	Password  string         `gorm:"not null;size:255" json:"-"`
	Status    int            `gorm:"default:1" json:"status"` // 1:正常 0:禁用
	Nickname  string         `gorm:"size:50" json:"nickname"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Role      string         `gorm:"size:50" json:"role"`
}

func (User) TableName() string {
	return "users"
}

type RegisterRequest struct {
	Username        string `json:"username" validate:"required,min=3,max=15,usernameUnique"`
	Password        string `json:"password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	Nickname        string `json:"nickname" validate:"required,max=50"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=15"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	Username        string `json:"username" validate:"omitempty,min=3,max=15,usernameUnique"`
	Password        string `json:"password" validate:"omitempty,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"omitempty,eqfield=Password"`
	Nickname        string `json:"nickname" validate:"omitempty,max=50"`
	Avatar          string `json:"avatar"`
	Status          int    `json:"status" validate:"oneof=1 0"`
	Role            string `json:"role" validate:"oneof=admin user"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Status   int    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserList struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}
