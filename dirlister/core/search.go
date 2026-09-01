package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID      uint   `gorm:"primaryKey"`
	Dir     string `gorm:"not null;index:idx_dir"`
	Name    string `gorm:"not null;index:idx_name"`
	Size    int64
	ModTime time.Time
	IsDir   bool
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
		rel, err := filepath.Rel(AppConf.ExposingDir, filepath.Dir(path))
		if err != nil {
			return err
		}

		if rel == "." {
			rel = ""
		}

		info, err := d.Info()

		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		if err != nil {
			fmt.Println("search.go:57 entry skipped", err)
		}

		file := File{Name: info.Name(), Dir: rel, IsDir: info.IsDir(), Size: info.Size(), ModTime: info.ModTime()}

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

func search(query string, user string, page int) ([]EntryInfo, error) {

	skip := page * AppConf.SearchResultCount

	remaining := AppConf.SearchResultCount + skip
	var match1 []File
	var match2 []File
	var match3 []File

	err := db.
		Where("name = ?", query).
		Limit(remaining).
		Find(&match1).Error

	if err != nil {
		return []EntryInfo{}, err
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
			return []EntryInfo{}, err
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
				return []EntryInfo{}, err
			}
		}

	}

	match := append(match1, append(match2, match3...)...)
	match = match[skip:]

	res := []EntryInfo{}

	for _, el := range match {
		_type := ""
		if el.IsDir {
			_type = "folder"
		} else {

			_type = "file"
		}
		res = append(res, EntryInfo{Name: el.Name, Type: _type, Size: el.Size, FullName: filepath.Join(el.Dir, el.Name), ModTime: el.ModTime})
	}

	return res, nil

}

func searchHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("p"))

	if err != nil {
		w.WriteHeader(400)
		return
	}

	b, e := io.ReadAll(r.Body)

	if e != nil {
		w.WriteHeader(500)
		return
	}
	query := r.URL.Query().Get("q")

	res, err := search(query, getSignedCookie(r, w), page)

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
