package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"kerimniy.qzz.io/dirlister/internal/config"
	db "kerimniy.qzz.io/dirlister/internal/database"

	"kerimniy.qzz.io/dirlister/internal/models"
)


func CreateRule(w http.ResponseWriter, r *http.Request) {

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

	err := db.Db.Create(&models.Rule{Path: path, CreatedAt: _time}).Error

	if err != nil {
		w.WriteHeader(500)
		fmt.Println(err)
		return
	}

	config.RulesTree.Tree.Insert(path, true)
	config.RulesTree.Dates[path] = _time
}

func DeleteRule(w http.ResponseWriter, r *http.Request) {

	if !checkAdmin(w, r) {
		w.WriteHeader(403)
		return
	}

	if r.Method != "DELETE" {
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
	payload := []string{}
	err := json.Unmarshal(b, &payload)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	for _, entry := range payload {

		result := db.Db.Unscoped().Where("path = ?", entry).Delete(models.Rule{})

		if result.Error != nil {
			w.WriteHeader(500)
			fmt.Println(result.Error)
		}

		config.RulesTree.Tree.Delete(entry)
		delete(config.RulesTree.Dates, entry)
	}
}

func GetRules(w http.ResponseWriter, r *http.Request) {

	if !checkAdmin(w, r) {
		w.WriteHeader(403)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("p"))

	if err != nil || page < 0 {

		w.WriteHeader(400)
		io.WriteString(w, "invalid page param")
		return
	}

	rulesList := []models.RuleResponse{}

	offset := page * config.AppConf.ResultCount

	var i = 0

	config.RulesTree.Tree.Walk(func(s string, v interface{}) bool {
		i++
		if i <= offset {
			return false
		} else if i <= offset+config.AppConf.ResultCount {

			rulesList = append(rulesList, models.RuleResponse{Path: s, CreatedAt: config.RulesTree.Dates[s]})

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
