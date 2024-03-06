package db

import (
	"fmt"
	"log"
	"os"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"rua.plus/gymo/models"
)

var Db *gorm.DB

func InitPostgres() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbMdoe := os.Getenv("DB_SSLMODE")
	dbTimeZone := os.Getenv("DB_TIMEZONE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host,
		user,
		password,
		dbName,
		dbPort,
		dbMdoe,
		dbTimeZone,
	)
	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("connect to postgress failed", err)
	}
	MigrateDb()
}

func MigrateDb() {
	if err := Db.AutoMigrate(&models.User{}); err != nil {
		log.Println(err)
		return
	}
	if err := Db.AutoMigrate(&models.Contact{}); err != nil {
		log.Println(err)
		return
	}
	if err := Db.AutoMigrate(&models.FirendRequest{}); err != nil {
		log.Println(err)
		return
	}
}

func NewMockDB() sqlmock.Sqlmock {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	Db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening gorm database", err)
	}

	return mock
}
