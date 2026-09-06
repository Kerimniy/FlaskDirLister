package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"kerimniy.qzz.io/dirlister/internal/models"
)

type Config struct {
	ExposingDir       string
	ResultCount       int
	SearchResultCount int
	UploadLimit       int
}

type AdminStruct struct {
	Exist bool
	Email string
}

var AppConf = Config{}
var Admin AdminStruct
var RulesTree models.RulesTree

func initConf() {

	err := godotenv.Load("./../.env")

	if err != nil {
		panic(err)
	}

	resultCount, err := strconv.Atoi(os.Getenv("RESULT_COUNT"))
	searchResultCount, err := strconv.Atoi(os.Getenv("SEARCH_RESULT_COUNT"))

	if err != nil {
		log.Fatal("ERROR 57 (get count) ", err)
	}

	_upload_limit, err := strconv.Atoi(os.Getenv("UPLOAD_LIMIT"))

	if err != nil {
		log.Fatal("ERROR 62 (get limit) ", err)
	}

	AppConf = Config{ExposingDir: os.Getenv("EXPDIR"), ResultCount: resultCount, SearchResultCount: searchResultCount, UploadLimit: _upload_limit}

}
