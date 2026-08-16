package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ExposingDir       string
	SearchResultCount int
}

type Admin struct {
	Exist bool
	Email string
}

// /home/kerimniy/Dev/dirlister/core/x/

var AppConf = Config{}
var admin Admin

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	err := godotenv.Load("./../.env")

	if err != nil {
		panic(err)
	}

	count, err := strconv.Atoi(os.Getenv("SEARCH_RESULT_COUNT"))

	if err != nil {
		log.Fatal("ERROR 57 (get count) ", err)
	}

	AppConf = Config{ExposingDir: os.Getenv("EXPDIR"), SearchResultCount: count}

	InitSecretKey()
	init_db()
	init_mail()
	initSearch()

	admin = _admin_exist()

	mux := http.NewServeMux()

	mux.HandleFunc("/s/{path...}", getDirHandle)
	mux.HandleFunc("/search", searchHandle)

	mux.HandleFunc("/auth/register", register)
	mux.HandleFunc("/auth/login", login)
	mux.HandleFunc("/auth/logout", logout)
	mux.HandleFunc("/auth/reset", reset_password)
	mux.HandleFunc("/auth/change", change_password)
	mux.HandleFunc("/auth/me", getUser)
	mux.HandleFunc("/auth/code", request_confirm_code)

	mux.HandleFunc("/rules/create", createRule)
	mux.HandleFunc("/rules/delete", deleteRule)

	mux.HandleFunc("/manage/upload", uploadHandle)
	mux.HandleFunc("/manage/delete", deleteHandle)
	mux.HandleFunc("/manage/delete-all", deleteAllHandle)
	mux.HandleFunc("/manage/rename", renameHandle)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-store")
		w.Write([]byte("hffhdfh"))
	})

	fmt.Println("Listening at: ", os.Getenv("HOST"))

	err = http.ListenAndServeTLS(os.Getenv("HOST"), os.Getenv("CERT"), os.Getenv("KEY"), corsMiddleware(mux))

	if err != nil {
		log.Fatal(0, err)
	}

}
