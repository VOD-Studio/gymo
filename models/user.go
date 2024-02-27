package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID          uint      `gorm:"primaryKey"                                   json:"id"`
	Email       string    `gorm:"unique;not null"                              json:"email"`
	Username    string    `gorm:"not null"                                     json:"username"`
	Password    string    `gorm:"not null"                                     json:"password"`
	Description string    `                                                    json:"description"`
	Avatar      string    `                                                    json:"avatar"`
	Gender      string    `                                                    json:"gender"`
	UID         int       `gorm:"unique;not null;default:100000;autoIncrement" json:"uid"`
	CreatedAt   time.Time `gorm:"default:NOW();not null"                       json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:NOW();not null"                       json:"updated_at"`
	LastLogin   int64     `                                                    json:"last_login"`
}

func (u *User) GetSingle(username string, db *gorm.DB) error {
	return db.Where("username = ?", username).First(u).Error
}

func (u *User) HashPassword() (err error) {
	if hash, err := HashPassword(u.Password); err != nil {
		return err
	} else {
		u.Password = hash
	}
	return
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	return u.HashPassword()
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
