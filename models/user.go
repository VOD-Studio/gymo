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
	Secret
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

// 发送的好友请求
type FirendRequest struct {
	ID          uint `gorm:"primaryKey"     json:"id,omitempty"`
	FromUserUID uint // 好友的 UID
	ToUserUID   uint // 自身的 UID
	FromUser    User `gorm:"references:UID"`
	ToUser      User `gorm:"references:UID"`
	Accept      bool
}

// 好友表
type Contact struct {
	ID        uint `gorm:"primaryKey"     json:"id,omitempty"`
	FirendUID uint // 好友的 UID
	UserUID   uint // 自身的 UID
	User      User `gorm:"references:UID"`
	Firend    User
}

// 用户表
type User struct {
	ID          uint      `gorm:"primaryKey"                     json:"id,omitempty"`
	Email       string    `gorm:"unique;not null"                json:"email"`
	Username    string    `gorm:"not null"                       json:"username"`
	Password    string    `gorm:"not null"                       json:"-"`
	Description string    `                                      json:"description"`
	Avatar      string    `                                      json:"avatar"`
	Gender      int8      `                                      json:"gender"`
	UID         uint      `gorm:"unique;not null;default:100000" json:"uid"`
	CreatedAt   time.Time `gorm:"default:NOW();not null"         json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"default:NOW();not null"         json:"updated_at,omitempty"`
	LastLogin   int64     `                                      json:"last_login,omitempty"`
	Onlie       bool
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
	user := &User{}
	tx.Model(&User{}).Order("uid desc").First(user, "")
	if user.UID > 0 {
		u.UID = user.UID + 1
	}
	return u.HashPassword()
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
