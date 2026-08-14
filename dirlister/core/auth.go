package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"gorm.io/gorm"
)

type Email struct {
	Email string `json:"email"`
}

type Register struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type ChangePassword struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

type ResetPassword struct {
	Email string `json:"email"`
	Code  string `json:"code"`

	NewPassword string `json:"newPassword"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func checkAdmin(w http.ResponseWriter, r *http.Request) bool {
	fmt.Println(getSignedCookie(r, w), admin.Email)
	return getSignedCookie(r, w) == admin.Email
}

func getUser(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Origin") != os.Getenv("FRONTEND") {
		w.WriteHeader(403)
		return
	}

	if admin.Exist == false {
		_, err := io.WriteString(w, "no-user")
		if err != nil {
			fmt.Println(err)
		}

		return
	}

	email := getSignedCookie(r, w)

	fmt.Println(email, 1)

	user := User{}
	err := db.Where("email= ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(200)

		_, err := io.WriteString(w, "no-user")
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(UserResponse{Email: user.Email, CreatedAt: user.CreatedAt})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(b)

}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	if admin.Exist == true {
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
	payload := Register{}
	err := json.Unmarshal(b, &payload)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	fmt.Println(payload.Code)
	if !validate_code(payload.Email, payload.Code) {
		w.WriteHeader(400)

		return
	}

	password_hash, err := hashPassword(payload.Password)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	user := User{
		Email:    payload.Email,
		Password: password_hash,
	}

	res := db.Create(&user)
	if res.Error != nil {
		w.WriteHeader(500)
		return
	}

	setSignedCookie(w, payload.Email)
	admin.Exist = true
	admin.Email = payload.Email

	w.WriteHeader(200)

}

func reset_password(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	b, e := io.ReadAll(r.Body)

	if e != nil {
		w.WriteHeader(500)
		return
	}
	payload := ResetPassword{}
	err := json.Unmarshal(b, &payload)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	if !validate_code(payload.Email, payload.Code) {
		w.WriteHeader(400)
		return
	}

	newPwd, err := hashPassword(payload.NewPassword)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	res := db.Model(&User{}).Where("email= ?", payload.Email).Update("password", newPwd)

	if res.Error != nil {
		w.WriteHeader(500)
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	if admin.Exist == false {
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
	payload := Login{}
	err := json.Unmarshal(b, &payload)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	qres := User{}
	res := db.Where("email= ?", payload.Email).First(&qres)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			w.WriteHeader(400)
		}

		if res.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err := io.WriteString(w, "User does not exist or invalid password")
		if err != nil {
			fmt.Println(err)
		}
		return
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

func logout(w http.ResponseWriter, r *http.Request) {
	deleteCookie(w)
}

func change_password(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	b, e := io.ReadAll(r.Body)

	if e != nil {
		w.WriteHeader(500)
		return
	}
	payload := ChangePassword{}
	err := json.Unmarshal(b, &payload)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}

	usermail := getSignedCookie(r, w)

	user := User{}
	res := db.Where("email= ?", usermail).First(&user)
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
	db.Model(&User{}).Where("email= ?", usermail).Update("password", newPwd)
	w.WriteHeader(200)
}
