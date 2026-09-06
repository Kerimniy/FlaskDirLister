package services

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"kerimniy.qzz.io/dirlister/internal/config"
	db "kerimniy.qzz.io/dirlister/internal/database"
	"kerimniy.qzz.io/dirlister/internal/models"
)

func listDir(dir string, urlPath string, user string, page int) ([]models.EntryInfo, error) {

	dirInfo := []models.EntryInfo{}

	result := []models.File{}

	db.Db.Where("dir = ?", urlPath).Offset(page * config.AppConf.ResultCount).Limit(config.AppConf.ResultCount).Find(&result)

	for _, entry := range result {

		if len(entry.Name) > 0 && entry.Name[0] == '.' {
			continue
		}

		_, _, match := config.RulesTree.Tree.LongestPrefix(entry.Dir)

		if match == true && user == "" {
			continue
		}

		var _type string
		if entry.IsDir {
			_type = "folder"
		} else {

			_type = "file"
		}

		entryInfo := models.EntryInfo{Name: entry.Name, Type: _type, ModTime: entry.ModTime, Size: entry.Size, FullName: filepath.Join(urlPath, entry.Name)}

		dirInfo = append(dirInfo, entryInfo)
	}

	return dirInfo, nil
}

func GetDirHandle(w http.ResponseWriter, r *http.Request) {

	urlPath := r.PathValue("path")

	path := filepath.Join(config.AppConf.ExposingDir, urlPath)

	rel, err := filepath.Rel(config.AppConf.ExposingDir, path)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	if strings.HasPrefix(rel, "..") {

		w.WriteHeader(400)
		return
	}

	_, _, match := config.RulesTree.Tree.LongestPrefix(urlPath)

	if match == true && !checkAdmin(w, r) {
		w.WriteHeader(403)
		return
	}

	info, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(204)
			return
		} else {
			w.WriteHeader(500)
			return
		}
	}

	if info.IsDir() {

		if r.URL.Query().Get("open") == "true" {
			w.WriteHeader(400)
			io.WriteString(w, "Folder is not editable")
			return
		}

		page, err := strconv.Atoi(r.URL.Query().Get("p"))

		if err != nil {
			w.WriteHeader(400)
			return
		}

		res, err := listDir(path, urlPath, getSignedCookie(r, w), page)

		if len(res) == 0 {
			w.WriteHeader(204)
			w.Write([]byte("[]"))
			return
		}

		if err != nil {
			w.WriteHeader(500)
			return
		}

		b, err := json.Marshal(res)

		if err != nil {
			w.WriteHeader(500)
			return
		}
		//w.Header().Set("Cache-Control", "private, max-age=60")
		w.Write(b)
	} else {

		fileName := filepath.Base(path)

		if r.URL.Query().Get("open") != "true" {
			w.Header().Set("Content-Disposition", "attachment; filename*=UTF-8''"+url.QueryEscape(fileName))
		}

		http.ServeFile(w, r, path)
	}

}

func listDirFS(dir string, urlPath string, user string) ([]models.EntryInfo, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	dirInfo := []models.EntryInfo{}

	for _, entry := range entries {

		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		_, _, match := config.RulesTree.Tree.LongestPrefix(urlPath)

		if len(info.Name()) > 0 && info.Name()[0] == '.' {
			continue
		}

		if match == true && user == "" {
			continue
		}

		var _type string
		if info.IsDir() {
			_type = "folder"
		} else {

			_type = "file"
		}

		entryInfo := models.EntryInfo{Name: info.Name(), Type: _type, ModTime: info.ModTime(), Size: info.Size(), FullName: filepath.Join(urlPath, info.Name())}

		dirInfo = append(dirInfo, entryInfo)
	}

	return dirInfo, nil
}
