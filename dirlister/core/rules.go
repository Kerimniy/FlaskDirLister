package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/armon/go-radix"
)

type RulesTree struct {
	mu   sync.RWMutex
	Tree *radix.Tree
}

type Rule struct {
	ID      int    `gorm:"primaryKey"`
	Pattern string `gorm:"unique"`
}

var rulesTree RulesTree

func createRule(w http.ResponseWriter, r *http.Request) {

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

	result := db.Create(payload)

	if result.Error != nil {
		w.WriteHeader(500)
		fmt.Println(result.Error)
	}

	rulesTree.Tree.Insert(payload.Pattern, true)
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

	rulesTree.Tree.Delete(payload.Pattern)
}

func getRules(w http.ResponseWriter, r *http.Request) {

	page, err := strconv.Atoi(r.URL.Query().Get("p"))

	if err != nil || page < 1 {

		w.WriteHeader(400)
		io.WriteString(w, "invalid page param")
		return
	}

	rulesList := []string{}

	page -= 1

	offset := page * AppConf.SearchResultCount

	var i = 0

	rulesTree.Tree.Walk(func(s string, v interface{}) bool {
		i++
		if i <= offset {
			return false
		} else if i <= offset+AppConf.SearchResultCount {

			rulesList = append(rulesList, s)

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

	fmt.Println(err)
}
