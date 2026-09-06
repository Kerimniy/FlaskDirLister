package transport

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"kerimniy.qzz.io/dirlister/internal/config"
	"kerimniy.qzz.io/dirlister/internal/services"
)

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

func ListenAndServeHTTP() {
	mux := http.NewServeMux()

	mux.HandleFunc("/s/{path...}", services.GetDirHandle)
	mux.HandleFunc("/search", services.SearchHandle)

	mux.HandleFunc("/auth/register", services.Register)
	mux.HandleFunc("/auth/login", services.Login)
	mux.HandleFunc("/auth/logout", services.Logout)
	mux.HandleFunc("/auth/reset", services.ResetPassword)
	mux.HandleFunc("/auth/change", services.Change_password)
	mux.HandleFunc("/auth/me", services.GetUser)
	mux.HandleFunc("/auth/code", services.RequestConfirmCode)

	mux.HandleFunc("/rules/create", services.CreateRule)
	mux.HandleFunc("/rules/delete", services.DeleteRule)
	mux.HandleFunc("/rules/get", services.GetRules)

	mux.HandleFunc("/manage/upload", services.UploadHandle)
	mux.HandleFunc("/manage/upload-multiple", services.UploadMultipleHandle)
	mux.HandleFunc("/manage/delete", services.DeleteHandle)
	mux.HandleFunc("/manage/delete-all", services.DeleteAllHandle)
	mux.HandleFunc("/manage/rename", services.RenameHandle)

	mux.HandleFunc("/info/upload-limit", getUploadLimit)
	mux.HandleFunc("/info/search-limit", getSearchLimit)
	mux.HandleFunc("/info/result-limit", getResultLimit)
	mux.HandleFunc("/info/disk-usage", services.GetUsageHandle)

	mux.HandleFunc("/", indexPage)

	if os.Getenv("CERT") != "" && os.Getenv("KEY") != "" {
		fmt.Println("Listening at: ", os.Getenv("HOST"), ", TLS enabled")

		err := http.ListenAndServeTLS(os.Getenv("HOST"), os.Getenv("CERT"), os.Getenv("KEY"), corsMiddleware(mux))

		if err != nil {
			log.Fatal(0, err)
		}
	} else {
		fmt.Println("Listening at: ", os.Getenv("HOST"), ", TLS disabled")

		err := http.ListenAndServe(os.Getenv("HOST"), corsMiddleware(mux))

		if err != nil {
			log.Fatal(0, err)
		}
	}
}

func getUploadLimit(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, strconv.Itoa(config.AppConf.UploadLimit))

}

func getSearchLimit(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, strconv.Itoa(config.AppConf.SearchResultCount))
}

func getResultLimit(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, strconv.Itoa(config.AppConf.ResultCount))
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, os.Getenv("FRONTEND"), 308)
}
