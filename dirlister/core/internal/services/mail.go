package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"math/rand"

	"kerimniy.qzz.io/dirlister/internal/models"
	"kerimniy.qzz.io/dirlister/pkg/tgbot"
)

var ctx = context.Background()

func validateCode(address string, code string) bool {

	res := authCode.Check(code)
	if res == true {
		authCode.Clear()
	}

	return res

}

func sendConfirm(address string) error {

	code := fmt.Sprintf("%06d", rand.Intn(999999))

	authCode.Set(code, time.Minute*10)

	return tgbot.SendCode(code)

}

func RequestConfirmCode(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}

	b, e := io.ReadAll(r.Body)

	if e != nil {
		w.WriteHeader(500)
		return
	}
	email := models.Email{}
	json.Unmarshal(b, &email)

	err := sendConfirm(email.Email)

	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, fmt.Sprintf("%s", err))
		return
	}

	w.WriteHeader(200)

}
