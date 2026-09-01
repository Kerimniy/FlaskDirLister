package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type uploadData struct {
	NewName string `json:"newName"`
	Content string `json:"content"`
}

func uploadHandle(w http.ResponseWriter, r *http.Request) {

	if !checkAdmin(w, r) {
		w.WriteHeader(403)
		return
	}

	filename := r.URL.Query().Get("file")

	isEdit := r.URL.Query().Get("edit") == "true"

	if !isEdit {
		path := filepath.Join(AppConf.ExposingDir, filename)

		rel, err := filepath.Rel(AppConf.ExposingDir, path)
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
		payload := uploadData{}
		err := json.Unmarshal(b, &payload)

		if err != nil {
			fmt.Println(err)

			w.WriteHeader(500)
			return
		}

		path := filepath.Join(AppConf.ExposingDir, filename)

		rel, err := filepath.Rel(AppConf.ExposingDir, path)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(400)
			return
		}

		newPath := filepath.Join(AppConf.ExposingDir, payload.NewName)

		rel, err = filepath.Rel(AppConf.ExposingDir, newPath)

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

	}
}

func deleteHandle(w http.ResponseWriter, r *http.Request) {

	if !checkAdmin(w, r) {
		w.WriteHeader(403)
		return
	}

	filename := r.URL.Query().Get("file")

	path := filepath.Join(AppConf.ExposingDir, filename)

	rel, err := filepath.Rel(AppConf.ExposingDir, path)

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
}

func deleteAllHandle(w http.ResponseWriter, r *http.Request) {

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

		path := filepath.Join(AppConf.ExposingDir, filename)

		rel, err := filepath.Rel(AppConf.ExposingDir, path)

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
	}

}

func renameHandle(w http.ResponseWriter, r *http.Request) {

	if !checkAdmin(w, r) {
		w.WriteHeader(403)
		return
	}

	query := r.URL.Query()

	fileName := filepath.Join(AppConf.ExposingDir, query.Get("file"))
	newName := filepath.Join(AppConf.ExposingDir, query.Get("name"))

	err := os.Rename(fileName, newName)

	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, "Rename error: "+err.Error())
	}

}

func CreateFile(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	return os.Create(path)
}

func uploadMultipleHandle(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	if !checkAdmin(w, r) {
		w.WriteHeader(403)
		return
	}

	dir := r.URL.Query().Get("dir")

	dir=strings.TrimLeft(dir,"/")
	dir=strings.Trim(dir,"/")

	err := r.ParseMultipartForm(int64(AppConf.UploadLimit))
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

		fullDir := filepath.Join(AppConf.ExposingDir, dir)
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

		stat, err:=newFile.Stat()

		if err != nil {
			w.WriteHeader(500)
			io.WriteString(w, "03 Couldn't index file:"+file.Filename)
			return
		}

		err=db.Create(File{Name: newFile.Name(),Dir: dir, Size: stat.Size(), IsDir: false, ModTime: time.Now()}).Error

		if err != nil {
			w.WriteHeader(500)
			io.WriteString(w, "03 Couldn't index file:"+file.Filename)
			return
		}

	}
}
