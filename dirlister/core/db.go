package main

import (
	"log"
	"sync"

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

	for _, rule := range rules {

		rulesTree.Tree.Insert(rule.Pattern, true)
	}

	db = _db
}
