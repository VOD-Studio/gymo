package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Gender int

const (
	Male Gender = iota
	Female
)

func (g Gender) String() string {
	switch g {
	case Male:
		return "male"
	case Female:
		return "female"
	}
	return "unknown"
}

type User struct {
	ID          uint      `gorm:"primaryKey"                                   json:"id,omitempty"`
	Email       string    `gorm:"unique;not null"                              json:"email"`
	Username    string    `gorm:"not null"                                     json:"username"`
	Password    string    `gorm:"not null"                                     json:"-"`
	Description string    `                                                    json:"description"`
	Avatar      string    `                                                    json:"avatar"`
	Gender      int8      `                                                    json:"gender"`
	UID         int       `gorm:"unique;not null;default:100000;autoIncrement" json:"uid"`
	CreatedAt   time.Time `gorm:"default:NOW();not null"                       json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"default:NOW();not null"                       json:"updated_at,omitempty"`
	LastLogin   int64     `                                                    json:"last_login,omitempty"`
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
