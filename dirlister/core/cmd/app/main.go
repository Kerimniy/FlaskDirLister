package main

import (
	"fmt"
	"time"

	"kerimniy.qzz.io/dirlister/internal/config"
	db "kerimniy.qzz.io/dirlister/internal/database"
	"kerimniy.qzz.io/dirlister/internal/services"
	"kerimniy.qzz.io/dirlister/internal/transport"
	"kerimniy.qzz.io/dirlister/pkg/tgbot"
)

func main() {

	fmt.Println(time.Now(), "Starting... ")

	config.InitConf()
	fmt.Println(time.Now(), "Initialized config ")

	services.InitSecretKey()
	fmt.Println(time.Now(), "Initialized secret key ")

	db.InitDb()
	fmt.Println(time.Now(), "Initialized database ")

	tgbot.InitTGBot()
	fmt.Println(time.Now(), "Initialized tg bot service ")

	services.InitSearch()
	fmt.Println(time.Now(), "Initialized search")

	config.Admin = services.Admin_exist()
	fmt.Println(time.Now(), "Initialized admin status")

	transport.ListenAndServeHTTP()
}
