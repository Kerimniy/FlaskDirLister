package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/armon/go-radix"
	"gorm.io/gorm"
)

type RulesTree struct {
	mu    sync.RWMutex
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
	Path      string `json:"path"`
	CreatedAt string `json:"createdAt"`
}

var rulesTree RulesTree

func createRule(w http.ResponseWriter, r *http.Request) {

	if !checkAdmin(w, r) {
		w.WriteHeader(403)
		return
	}
	path := r.URL.Query().Get("p")
	if strings.Trim(getSignedCookie(r, w), " ") == "" {
		w.WriteHeader(403)
		return
	}

	_time := time.Now()

	err := db.Create(&Rule{Path: path, CreatedAt: _time}).Error

	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	rulesTree.Tree.Insert(path, true)
	rulesTree.Dates[path] = _time
}

func deleteRule(w http.ResponseWriter, r *http.Request) {

	if !checkAdmin(w, r) {
		w.WriteHeader(403)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	if strings.Trim(getSignedCookie(r, w), " ") == "" {
		w.WriteHeader(403)
		return
	}

	b, e := io.ReadAll(r.Body)

	if e != nil {
		w.WriteHeader(500)
		return
	}
	payload := Rule{}
	err := json.Unmarshal(b, &payload)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	result := db.Unscoped().Delete(payload)

	if result.Error != nil {
		w.WriteHeader(500)
		fmt.Println(result.Error)
	}

	rulesTree.Tree.Delete(payload.Path)
}

func getRules(w http.ResponseWriter, r *http.Request) {

	page, err := strconv.Atoi(r.URL.Query().Get("p"))

	if err != nil || page < 0 {

		w.WriteHeader(400)
		io.WriteString(w, "invalid page param")
		return
	}

	rulesList := []RuleResponse{}

	offset := page * AppConf.SearchResultCount

	var i = 0

	rulesTree.Tree.Walk(func(s string, v interface{}) bool {
		i++
		if i <= offset {
			return false
		} else if i <= offset+AppConf.SearchResultCount {

			rulesList = append(rulesList, RuleResponse{Path: s, CreatedAt: rulesTree.Dates[s].String()})

			return false
		} else {
			return true
		}

	})

	b, err := json.Marshal(rulesList)

	if err != nil {
		fmt.Println("rules.go:143", err)
		w.WriteHeader(500)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		fmt.Println(err)
	}
}
