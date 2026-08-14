package main

import (
	"errors"
	"time"

	"sync"

	_ "gorm.io/driver/sqlite"
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

type AuthCode struct {
	mu      sync.RWMutex
	code    string
	expires time.Time
}

var authCode = AuthCode{}

func _admin_exist() Admin {
	var user User

	err := db.Limit(1).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Admin{Exist: false}
	}
	
	return Admin{Email: user.Email, Exist: true}
}

func (a *AuthCode) Set(code string, ttl time.Duration) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.code = code
	a.expires = time.Now().Add(ttl)
}

func (a *AuthCode) Check(code string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.code == code && time.Now().Before(a.expires)
}

func (a *AuthCode) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.code = ""
	a.expires = time.Time{}
}
