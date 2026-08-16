package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"gorm.io/gorm"
)

type File struct {
	ID   uint   `gorm:"primaryKey"`
	Dir  string `gorm:"not null;index:idx_dir"`
	Name string `gorm:"not null;index:idx_name"`
}

type SearchResponse struct {
	Match1 []File
	Match2 []File
	Match3 []File
}

func initSearch() {

	err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&File{}).Error

	if err != nil {
		log.Fatal("ERROR 29 (clean search db) ", err)
	}

	err = filepath.WalkDir(AppConf.ExposingDir, func(path string, d os.DirEntry, err error) error {

		if err != nil {
			return err
		}
		rel, err := filepath.Rel(AppConf.ExposingDir, path)

		if err != nil {
			return err
		}

		if rel[0] == '.' {
			return nil
		}

		file := File{Name: filepath.Base(rel), Dir: rel}

		res := db.Create(&file)

		if res.Error != nil {
			return res.Error
		}

		return nil
	})
	if err != nil {
		log.Fatal("ERROR 001 ", err)
	}

}

func search(query string, user string) (SearchResponse, error) {

	remaining := AppConf.SearchResultCount

	var match1 []File
	var match2 []File
	var match3 []File

	err := db.
		Where("name = ?", query).
		Limit(remaining).
		Find(&match1).Error

	if err != nil {
		return SearchResponse{}, err
	}

	match1 = filter(match1, user)

	if remaining > len(match1) {
		remaining -= len(match1)

		var ids []uint
		for _, f := range match1 {
			ids = append(ids, f.ID)
		}

		q := db.
			Where("name LIKE ?", query+"%").
			Limit(remaining)

		if len(ids) > 0 {
			q = q.Where("id NOT IN ?", ids)
		}

		err = q.Find(&match2).Error

		if err != nil {
			return SearchResponse{}, err
		}

		match2 = filter(match2, user)

		if remaining > len(match2) {
			remaining -= len(match2)

			for _, f := range match2 {
				ids = append(ids, f.ID)
			}

			q := db.
				Where("name LIKE ?", "%"+query+"%").
				Limit(remaining)

			if len(ids) > 0 {
				q = q.Where("id NOT IN ?", ids)
			}

			err = q.Find(&match3).Error

			match3 = filter(match3, user)

			if err != nil {
				return SearchResponse{}, err
			}
		}

	}

	return SearchResponse{Match1: match1, Match2: match2, Match3: match3}, nil

}

func searchHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}

	b, e := io.ReadAll(r.Body)

	if e != nil {
		w.WriteHeader(500)
		return
	}
	query := r.URL.Query().Get("q")

	res, err := search(query, getSignedCookie(r, w))

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	b, err = json.Marshal(res)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Write(b)
}

func filter(match []File, user string) []File {
	filtered := make([]File, 0, len(match))

	for _, file := range match {
		_, _, exclude := rulesTree.Tree.LongestPrefix(file.Dir)

		if exclude && user == "" {
			continue
		}

		filtered = append(filtered, file)
	}

	return filtered
}
