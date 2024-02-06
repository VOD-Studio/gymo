package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var Db *gorm.DB

func InitPostgres() {
	dsn := "host=192.168.1.57 user=dev password=qwer1234 dbname=dev port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("connect to postgress failed", err)
	}
}
