package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
