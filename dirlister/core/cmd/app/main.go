package main

import (
	"kerimniy.qzz.io/dirlister/internal/config"
	db "kerimniy.qzz.io/dirlister/internal/database"
	"kerimniy.qzz.io/dirlister/internal/services"
	"kerimniy.qzz.io/dirlister/internal/transport"
)

func main() {

	services.InitSecretKey()
	db.InitDb()
	services.InitMail()
	services.InitSearch()
	config.Admin = services.Admin_exist()

	transport.ListenAndServeHTTP()
}
