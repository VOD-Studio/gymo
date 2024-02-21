package models

import (
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
	LastLogin int64     `                              json:"last_login"`
}

func (u *User) GetSingle(username string, db *gorm.DB) error {
	return db.Where("username = ?", username).First(u).Error
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if hash, err := HashPassword(u.Password); err != nil {
		return err
	} else {
		u.Password = hash
	}
	return
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
