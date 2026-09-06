package main

import (
	"fmt"
	"time"

	"kerimniy.qzz.io/dirlister/internal/config"
	db "kerimniy.qzz.io/dirlister/internal/database"
	"kerimniy.qzz.io/dirlister/internal/services"
	"kerimniy.qzz.io/dirlister/internal/transport"
)

func main() {

	fmt.Println(time.Now(), "Starting... ")

	services.InitSecretKey()
	fmt.Println(time.Now(), "Initialized secret key ")

	db.InitDb()
	fmt.Println(time.Now(), "Initialized database ")

	services.InitMail()
	fmt.Println(time.Now(), "Initialized mail service ")

	services.InitSearch()
	fmt.Println(time.Now(), "Initialized search")

	config.Admin = services.Admin_exist()
	fmt.Println(time.Now(), "Initialized admin status")

	transport.ListenAndServeHTTP()
}
