package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primaryKey"             json:"id"`
	Email     string    `gorm:"unique;not null"        json:"email"`
	Username  string    `gorm:"not null"               json:"username"`
	Password  string    `gorm:"not null"               json:"password"`
	CreatedAt time.Time `gorm:"default:NOW();not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:NOW();not null" json:"updated_at"`
}

func (u *User) GetSingle(username string, db *gorm.DB) error {
	return db.Where("username = ?", username).First(u).Error
}

var UserAlreadyExist = errors.New("user already exist")

func (u *User) Create(db *gorm.DB) error {
	if hash, err := HashPassword(u.Password); err != nil {
		return err
	} else {
		u.Password = hash
	}

	res := db.FirstOrCreate(u)
	if err := res.Error; err != nil {
		return err
	}

	if res.RowsAffected == 1 {
		return nil
	} else {
		return UserAlreadyExist
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
