package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string `gorm:"unique;not null"`
	CreatedAt time.Time
	Password  []byte `gorm:"not null"`
}

type UserResponse struct {
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

type Email struct {
	Email string `json:"email"`
}

type Register struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type ChangePassword struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type ResetPassword struct {
	Email string `json:"email"`
	Code  string `json:"code"`

	NewPassword string `json:"newPassword"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
