package services

import (
	"errors"
	"time"

	"sync"

	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"kerimniy.qzz.io/dirlister/internal/config"
	db "kerimniy.qzz.io/dirlister/internal/database"

	"kerimniy.qzz.io/dirlister/internal/models"
)

type AuthCode struct {
	mu      sync.RWMutex
	code    string
	expires time.Time
}

var authCode = AuthCode{}

func Admin_exist() config.AdminStruct {
	var user models.User

	err := db.Db.Limit(1).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return config.AdminStruct{Exist: false}
	}

	return config.AdminStruct{Email: user.Email, Exist: true}
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
