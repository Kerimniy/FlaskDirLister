package services

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

	"gorm.io/gorm"

	"kerimniy.qzz.io/dirlister/internal/config"
	db "kerimniy.qzz.io/dirlister/internal/database"

	"kerimniy.qzz.io/dirlister/internal/models"
)

func InitSearch() {

	err := db.Db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.File{}).Error

	if err != nil {
		log.Fatal("ERROR 29 (clean search db) ", err)
	}

	err = filepath.WalkDir(config.AppConf.ExposingDir, func(path string, d os.DirEntry, err error) error {

		if err != nil {
			return err
		}
		rel, err := filepath.Rel(config.AppConf.ExposingDir, filepath.Dir(path))
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

		file := models.File{Name: info.Name(), Dir: rel, IsDir: info.IsDir(), Size: info.Size(), ModTime: info.ModTime()}

		res := db.Db.Create(&file)

		if res.Error != nil {
			return res.Error
		}

		return nil
	})
	if err != nil {
		log.Fatal("ERROR 001 ", err)
	}

}

func search(query string, user string, page int) ([]models.EntryInfo, error) {

	skip := page * config.AppConf.SearchResultCount

	remaining := config.AppConf.SearchResultCount + skip
	var match1 []models.File
	var match2 []models.File
	var match3 []models.File

	err := db.Db.
		Where("name = ?", query).
		Limit(remaining).
		Find(&match1).Error

	if err != nil {
		return []models.EntryInfo{}, err
	}

	match1 = filter(match1, user)

	if remaining > len(match1) {
		remaining -= len(match1)

		var ids []uint
		for _, f := range match1 {
			ids = append(ids, f.ID)
		}

		q := db.Db.
			Where("name LIKE ?", query+"%").
			Limit(remaining)

		if len(ids) > 0 {
			q = q.Where("id NOT IN ?", ids)
		}

		err = q.Find(&match2).Error

		if err != nil {
			return []models.EntryInfo{}, err
		}

		match2 = filter(match2, user)

		if remaining > len(match2) {
			remaining -= len(match2)

			for _, f := range match2 {
				ids = append(ids, f.ID)
			}

			q := db.Db.
				Where("name LIKE ?", "%"+query+"%").
				Limit(remaining)

			if len(ids) > 0 {
				q = q.Where("id NOT IN ?", ids)
			}

			err = q.Find(&match3).Error

			match3 = filter(match3, user)

			if err != nil {
				return []models.EntryInfo{}, err
			}
		}

	}

	match := append(match1, append(match2, match3...)...)
	match = match[skip:]

	res := []models.EntryInfo{}

	for _, el := range match {
		_type := ""
		if el.IsDir {
			_type = "folder"
		} else {

			_type = "file"
		}
		res = append(res, models.EntryInfo{Name: el.Name, Type: _type, Size: el.Size, FullName: filepath.Join(el.Dir, el.Name), ModTime: el.ModTime})
	}

	return res, nil

}

func SearchHandle(w http.ResponseWriter, r *http.Request) {
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

func filter(match []models.File, user string) []models.File {
	filtered := make([]models.File, 0, len(match))

	for _, file := range match {
		_, _, exclude := config.RulesTree.Tree.LongestPrefix(file.Dir)

		if exclude && user == "" {
			continue
		}

		filtered = append(filtered, file)
	}

	return filtered
}
