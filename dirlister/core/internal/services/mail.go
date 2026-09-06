package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"bytes"
	"math/rand"

	"github.com/wneessen/go-mail"
	"kerimniy.qzz.io/dirlister/internal/models"
)

type Mail struct {
	From string

	Username string
	Password string
	SMTP     string
}

var mail_conf Mail
var ctx = context.Background()

//var tmpl = template.Must(template.ParseFiles("templates/recovery.html"))

func InitMail() {

	mail_conf = Mail{
		Username: os.Getenv("USER"),
		From:     os.Getenv("FROM"),
		Password: os.Getenv("PASSWORD"),
		SMTP:     os.Getenv("SMTP"),
	}

}

func validateCode(address string, code string) bool {

	res := authCode.Check(code)
	if res == true {
		authCode.Clear()
	}

	return res

}

func sendConfirm(address string) error {

	code := fmt.Sprintf("%06d", rand.Intn(999999))

	fmt.Println(code)

	authCode.Set(code, time.Minute*10)

	var tpl bytes.Buffer

	//tmpl.Execute(&tpl, code)

	from := mail_conf.From
	host := mail_conf.SMTP
	username := mail_conf.Username
	password := mail_conf.Password

	message := mail.NewMsg()
	if err := message.From(from); err != nil {
		return fmt.Errorf("failed to set From address: %s", err)
	}
	if err := message.To(address); err != nil {
		return fmt.Errorf("failed to set To address: %s", err)
	}

	message.Subject("Email confirmation")
	message.SetBodyString(mail.TypeTextHTML, tpl.String())

	client, err := mail.NewClient(host,

		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(username),
		mail.WithPassword(password),
		mail.WithTimeout(time.Duration(10)*time.Second),
	)

	if err != nil {
		return fmt.Errorf("failed to create mail client: %s", err)
	}
	if err := client.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send mail: %s", err)
	}
	return nil
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
