package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type uploadData struct {
	NewName string `json:"newName"`
	Content string `json:"content"`
}

func uploadHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("75", getSignedCookie(r, w))

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
