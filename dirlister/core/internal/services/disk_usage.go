package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/shirou/gopsutil/disk"

	"regexp"

	"kerimniy.qzz.io/dirlister/internal/config"
)

type Usage struct {
	Base    uint64  `json:"base"`
	Free    uint64  `json:"free"`
	Percent float64 `json:"percent"`
}

var re = regexp.MustCompile(`^(?:/mnt/[^/]+/|[A-Za-z]:\\)`)

func getUsage(path string) (Usage, error) {

	path, err := filepath.Abs(path)
	if err != nil {
		return Usage{}, err
	}

	root := re.FindString(path)
	if root == "" {
		root = "/"
	}

	usage, err := disk.Usage(root)

	if err != nil {
		return Usage{}, err
	}

	return Usage{Percent: usage.UsedPercent, Base: usage.Used + usage.Free, Free: usage.Free}, nil

}

func GetUsageHandle(w http.ResponseWriter, r *http.Request) {

	if !checkAdmin(w, r) {
		w.WriteHeader(403)
		return
	}

	res, err := getUsage(config.AppConf.ExposingDir)

	if err != nil {
		w.WriteHeader(500)
		fmt.Println("diskUsage.go:55", err)
		return
	}

	b, err := json.Marshal(res)

	w.Header().Set("Cache-Control", "private, max-age=60")

	w.Write(b)
}
