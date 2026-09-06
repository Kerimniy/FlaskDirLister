package db

import (
	"log"
	"sync"
	"time"

	"github.com/armon/go-radix"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"kerimniy.qzz.io/dirlister/internal/config"
	"kerimniy.qzz.io/dirlister/internal/models"

)

var Db *gorm.DB

func InitDb() {

	_db, err := gorm.Open(sqlite.Open(".main.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Db error: %v", err)
	}
	_db.AutoMigrate(&models.User{})
	_db.AutoMigrate(&models.Rule{})
	_db.AutoMigrate(&models.File{})

	config.RulesTree = models.RulesTree{
		Mu:   sync.RWMutex{},
		Tree: radix.New(),
	}

	rules := []models.Rule{}
	res := _db.Find(&rules)

	if res.Error != nil {
		log.Fatal("db.go:33 ", res.Error)
	}
	config.RulesTree.Dates = make(map[string]time.Time)
	for _, rule := range rules {

		config.RulesTree.Tree.Insert(rule.Path, true)
		config.RulesTree.Dates[rule.Path] = rule.CreatedAt
	}

	Db = _db
}
