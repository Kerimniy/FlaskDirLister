package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"kerimniy.qzz.io/dirlister/internal/config"
	db "kerimniy.qzz.io/dirlister/internal/database"
	"kerimniy.qzz.io/dirlister/internal/models"
)

func UploadHandle(w http.ResponseWriter, r *http.Request) {

	if !checkAdmin(w, r) {
		w.WriteHeader(403)
		return
	}

	filename := r.URL.Query().Get("file")
	filename = strings.Trim(strings.TrimLeft(filename, "/"), "/")

	isEdit := r.URL.Query().Get("edit") == "true"

	if !isEdit {
		path := filepath.Join(config.AppConf.ExposingDir, filename)

		rel, err := filepath.Rel(config.AppConf.ExposingDir, path)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		if strings.HasPrefix(rel, "..") {

			w.WriteHeader(400)
			return
		}

		file, err := CreateFile(path)

		if err != nil {
			w.WriteHeader(500)
			return
		}

		defer file.Close()

		_, err = io.Copy(file, r.Body)

		if err != nil {
			w.WriteHeader(500)
			return
		}
		stat, err := file.Stat()
		err = indexNewDirs(filename, stat.Size(), false)

		if err != nil {
			fmt.Println(err)

			w.WriteHeader(500)
			return
		}

	} else {

		if r.Method == "OPTIONS" {
			return
		}

		b, e := io.ReadAll(r.Body)

		if e != nil {
			fmt.Println(e)

			w.WriteHeader(500)
			return
		}
		payload := models.UploadData{}
		err := json.Unmarshal(b, &payload)

		if err != nil {
			fmt.Println(err)

			w.WriteHeader(500)
			return
		}

		path := filepath.Join(config.AppConf.ExposingDir, filename)

		rel, err := filepath.Rel(config.AppConf.ExposingDir, path)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}

		newPath := filepath.Join(config.AppConf.ExposingDir, payload.NewName)

		rel, err = filepath.Rel(config.AppConf.ExposingDir, newPath)

		if err != nil {
			fmt.Println(err)

			w.WriteHeader(400)
			return
		}

		if strings.HasPrefix(rel, "..") {

			w.WriteHeader(400)
			return
		}

		file, err := CreateFile(path)

		if err != nil {
			fmt.Println(err)

			w.WriteHeader(500)
			return
		}

		_, err = io.WriteString(file, payload.Content)

		if err != nil {
			fmt.Println(err)

			w.WriteHeader(500)
			return
		}
		file.Close()

		if newPath != "" {
			err = os.Rename(path, newPath)

			if err != nil {
				fmt.Println(err)

				w.WriteHeader(500)
				return
			}

		}

		stat, err := file.Stat()
		err = indexNewDirs(rel, stat.Size(), false)

		if err != nil {
			fmt.Println(err)

			w.WriteHeader(500)
			return
		}

	}

}

func DeleteHandle(w http.ResponseWriter, r *http.Request) {

	if !checkAdmin(w, r) {
		w.WriteHeader(403)
		return
	}

	filename := r.URL.Query().Get("file")
	filename = strings.Trim(strings.TrimLeft(filename, "/"), "/")

	path := filepath.Join(config.AppConf.ExposingDir, filename)

	rel, err := filepath.Rel(config.AppConf.ExposingDir, path)

	if err != nil {
		w.WriteHeader(400)
		return
	}

	if strings.HasPrefix(rel, "..") {

		w.WriteHeader(400)
		return
	}

	err = os.RemoveAll(path)

	if err != nil {
		w.WriteHeader(500)
		return
	}

	err = indexNewDirs(filename, 0, true)
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func DeleteAllHandle(w http.ResponseWriter, r *http.Request) {

	if !checkAdmin(w, r) {
		w.WriteHeader(403)
		return
	}

	filenames := []string{}

	b, e := io.ReadAll(r.Body)

	if e != nil {
		w.WriteHeader(500)
		return
	}

	err := json.Unmarshal(b, &filenames)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	for _, filename := range filenames {

		path := filepath.Join(config.AppConf.ExposingDir, filename)
		filename = strings.Trim(strings.TrimLeft(filename, "/"), "/")

		rel, err := filepath.Rel(config.AppConf.ExposingDir, path)

		if err != nil {
			w.WriteHeader(400)
			return
		}

		if strings.HasPrefix(rel, "..") {

			w.WriteHeader(400)
			return
		}

		err = os.RemoveAll(path)
		if err != nil {
			w.WriteHeader(500)
			io.WriteString(w, filename)
			return
		}

		err = indexNewDirs(filename, 0, true)
		if err != nil {
			fmt.Println(err)

			w.WriteHeader(500)
			return
		}

	}

}

func RenameHandle(w http.ResponseWriter, r *http.Request) {

	if !checkAdmin(w, r) {
		w.WriteHeader(403)
		return
	}

	query := r.URL.Query()

	_f := strings.Trim(strings.TrimLeft(query.Get("file"), "/"), "/")
	_n := strings.Trim(strings.TrimLeft(query.Get("name"), "/"), "/")

	fileName := filepath.Join(config.AppConf.ExposingDir, _f)
	newName := filepath.Join(config.AppConf.ExposingDir, _n)

	err := os.Rename(fileName, newName)

	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, "Rename error: "+err.Error())
		return
	}

	file, err := os.Open(newName)

	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, "Rename error: "+err.Error())
		return
	}

	stat, err := file.Stat()

	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, "Rename error: "+err.Error())
		return
	}

	err = indexNewDirs(_n, stat.Size(), false)

	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, "Rename error: "+err.Error())
		return
	}

	err = indexNewDirs(_f, 0, true)

	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, "Rename error: "+err.Error())
		return
	}

}

func CreateFile(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	return os.Create(path)
}

func UploadMultipleHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	if !checkAdmin(w, r) {
		w.WriteHeader(403)
		return
	}

	dir := r.URL.Query().Get("dir")

	dir = strings.TrimLeft(dir, "/")
	dir = strings.Trim(dir, "/")

	err := r.ParseMultipartForm(int64(config.AppConf.UploadLimit))
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, "Parse form error:"+err.Error())
		return
	}
	files := r.MultipartForm.File["files"]

	for _, file := range files {
		f, err := file.Open()

		if err != nil {
			w.WriteHeader(500)
			io.WriteString(w, "00 Couldn't load file:"+file.Filename)
			return
		}

		defer f.Close()

		fullDir := filepath.Join(config.AppConf.ExposingDir, dir)

		rel, err := filepath.Rel(config.AppConf.ExposingDir, fullDir)
		if err != nil {
			w.WriteHeader(400)
			return
		}

		relName := filepath.Join(rel, file.Filename)

		fullName := filepath.Join(fullDir, file.Filename)

		err = os.MkdirAll(fullDir, os.ModePerm)

		if err != nil {
			w.WriteHeader(500)
			io.WriteString(w, "03 Couldn't load file:"+file.Filename)
			return
		}

		newFile, err := os.Create(fullName)

		if err != nil {
			w.WriteHeader(500)
			io.WriteString(w, "01 Couldn't load file:"+file.Filename)
			return
		}

		_, err = io.Copy(newFile, f)

		defer newFile.Close()

		if err != nil {
			w.WriteHeader(500)
			io.WriteString(w, "02 Couldn't load file:"+file.Filename)
			return
		}

		stat, err := newFile.Stat()

		if err != nil {
			w.WriteHeader(500)
			io.WriteString(w, "03 Couldn't index file:"+file.Filename)
			return
		}

		err = indexNewDirs(relName, stat.Size(), false)

		if err != nil {
			w.WriteHeader(500)
			io.WriteString(w, "03 Couldn't index file:"+file.Filename)
			return
		}

	}
}

func indexNewDirs(path string, fileSize int64, delete bool) error {
	f := false
	r := false
	for {
		if r {
			break
		}

		name := filepath.Base(path)
		path = filepath.Dir(path)

		if path == "." || path == string(filepath.Separator) {
			path = ""

			r = true
		}

		var _fileSize int64 = 0

		if !f {
			_fileSize = fileSize
		}

		var result *gorm.DB

		if delete {

			result = db.Db.Unscoped().Where("dir = ?", path).Where("name = ?", name).Delete(&models.File{})

		} else {

			result = db.Db.Clauses(clause.OnConflict{
				DoNothing: true,
			}).Create(&models.File{
				Name:    name,
				Dir:     path,
				IsDir:   f,
				ModTime: time.Now(),
				Size:    _fileSize,
			})
		}

		if f {
			f = true
		}

		if result.Error != nil {
			return result.Error
		}

	}

	return nil

}
