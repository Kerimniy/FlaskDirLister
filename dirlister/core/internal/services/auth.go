package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"gorm.io/gorm"
	"kerimniy.qzz.io/dirlister/internal/config"
	db "kerimniy.qzz.io/dirlister/internal/database"
	"kerimniy.qzz.io/dirlister/internal/models"
)

func checkAdmin(w http.ResponseWriter, r *http.Request) bool {

	fmt.Println(getSignedCookie(r, w), config.Admin.Email)
	return getSignedCookie(r, w) == config.Admin.Email
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("77", getSignedCookie(r, w))

	if r.Header.Get("Origin") != os.Getenv("FRONTEND") {
		w.WriteHeader(403)
		return
	}

	if config.Admin.Exist == false {
		_, err := io.WriteString(w, "no-user")
		if err != nil {
			fmt.Println(err)
		}

		return
	}

	email := getSignedCookie(r, w)
	fmt.Println("11", email)
	user := models.User{}
	err := db.Db.Where("email= ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(401)

		_, err := io.WriteString(w, "unauthorized")
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(models.UserResponse{Email: user.Email, CreatedAt: user.CreatedAt})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(b)

}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	if config.Admin.Exist == true {
		w.WriteHeader(412)
		_, err := io.WriteString(w, "Admin already exists")
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	b, e := io.ReadAll(r.Body)

	if e != nil {
		w.WriteHeader(500)
		return
	}
	payload := models.Register{}
	err := json.Unmarshal(b, &payload)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	fmt.Println(payload.Code)
	if !validateCode(payload.Email, payload.Code) {
		w.WriteHeader(400)

		return
	}

	password_hash, err := hashPassword(payload.Password)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	user := models.User{
		Email:    payload.Email,
		Password: password_hash,
	}

	res := db.Db.Create(&user)
	if res.Error != nil {
		w.WriteHeader(500)
		return
	}

	setSignedCookie(w, payload.Email)
	config.Admin.Exist = true
	config.Admin.Email = payload.Email

	w.WriteHeader(200)

}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	b, e := io.ReadAll(r.Body)

	if e != nil {
		w.WriteHeader(500)
		return
	}
	payload := models.ResetPassword{}
	err := json.Unmarshal(b, &payload)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	if !validateCode(payload.Email, payload.Code) {
		w.WriteHeader(400)
		return
	}

	newPwd, err := hashPassword(payload.NewPassword)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	res := db.Db.Model(&models.User{}).Where("email= ?", payload.Email).Update("password", newPwd)

	if res.Error != nil {
		w.WriteHeader(500)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	if config.Admin.Exist == false {
		w.WriteHeader(412)
		_, err := io.WriteString(w, "Admin not exists")
		if err != nil {
			fmt.Println(err)
		}

		return
	}

	b, e := io.ReadAll(r.Body)

	if e != nil {
		w.WriteHeader(500)
		return
	}
	payload := models.Login{}
	err := json.Unmarshal(b, &payload)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	qres := models.User{}
	res := db.Db.Where("email= ?", payload.Email).First(&qres)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			w.WriteHeader(400)

			_, err := io.WriteString(w, "User does not exist or invalid password")
			if err != nil {
				fmt.Println(err)
			}

			return
		}

		if res.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}
	verified := checkPassword(qres.Password, payload.Password)

	if !verified {

		w.WriteHeader(400)
		_, err := io.WriteString(w, "User does not exist or invalid password")
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	setSignedCookie(w, payload.Email)

	w.WriteHeader(200)

}

func Logout(w http.ResponseWriter, r *http.Request) {
	deleteCookie(w)
}

func Change_password(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	b, e := io.ReadAll(r.Body)

	if e != nil {
		w.WriteHeader(500)
		return
	}
	payload := models.ChangePassword{}
	err := json.Unmarshal(b, &payload)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	usermail := getSignedCookie(r, w)

	user := models.User{}
	res := db.Db.Where("email= ?", usermail).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if res.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if !checkPassword(user.Password, payload.CurrentPassword) {
		w.WriteHeader(400)
		return
	}
	newPwd, err := hashPassword(payload.NewPassword)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
	db.Db.Model(&models.User{}).Where("email= ?", usermail).Update("password", newPwd)
	w.WriteHeader(200)
}
