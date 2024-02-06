package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"rua.plus/gymo/models"
)

var Db *gorm.DB

func InitPostgres() {
	dsn := "host=192.168.1.57 user=xfy password=qwer1234 dbname=dev port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("connect to postgress failed", err)
	}
	MigrateDb()
}

func MigrateDb() {
	err := Db.AutoMigrate(&models.User{})
	if err != nil {
		log.Println(err)
		return
	}
}
