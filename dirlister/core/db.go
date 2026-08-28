package main

import (
	"log"
	"sync"
	"time"

	"github.com/armon/go-radix"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init_db() {

	_db, err := gorm.Open(sqlite.Open("main.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Db error: %v", err)
	}
	_db.AutoMigrate(&User{})
	_db.AutoMigrate(&Rule{})
	_db.AutoMigrate(&File{})

	rulesTree = RulesTree{
		mu:   sync.RWMutex{},
		Tree: radix.New(),
	}

	rules := []Rule{}
	res := _db.Find(&rules)

	if res.Error != nil {
		log.Fatal("db.go:33 ", res.Error)
	}
	rulesTree.Dates = make(map[string]time.Time)
	for _, rule := range rules {

		rulesTree.Tree.Insert(rule.Path, true)
		rulesTree.Dates[rule.Path] = rule.CreatedAt
	}

	db = _db
}
