package models

import (
	"sync"
	"time"

	"github.com/armon/go-radix"
	"gorm.io/gorm"
)

type RulesTree struct {
	Mu    sync.RWMutex
	Tree  *radix.Tree
	Dates map[string]time.Time
}

type Rule struct {
	gorm.Model
	ID        int    `gorm:"primaryKey"`
	Path      string `gorm:"unique"`
	CreatedAt time.Time
}

type RuleResponse struct {
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"createdAt"`
}
