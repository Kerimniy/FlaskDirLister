package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

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

var AppConf = Config{ExposingDir: "/home/kerimniy/Dev/dirlister/core/x/", SearchResultCount: 15}
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
	mux.HandleFunc("/manage/rename", renameHandle)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "204") })

	fmt.Println("Listening at", os.Getenv("HOST"))

	err = http.ListenAndServeTLS(os.Getenv("HOST"), "/home/kerimniy/localhost+2.pem", "/home/kerimniy/localhost+2-key.pem", corsMiddleware(mux))

	if err != nil {
		log.Fatal(err)
	}

}
