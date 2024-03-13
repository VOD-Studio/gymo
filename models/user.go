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

type BaseID struct {
	ID        uint           `gorm:"primaryKey"             json:"id,omitempty"`
	CreatedAt time.Time      `gorm:"default:NOW();not null" json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"default:NOW();not null" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index"                  json:"-"`
}

// 发送的好友请求
type FirendRequest struct {
	BaseID
	FromUserUID uint `json:"from_user_uid"` // 自身的 UID
	ToUserUID   uint `json:"to_user_uid"`   // 好友的 UID
	FromUser    User `json:"from_user"     gorm:"references:UID"`
	ToUser      User `json:"to_user"       gorm:"references:UID"`
	Accept      bool
}

// 好友表
type Contact struct {
	BaseID
	FirendUID uint `json:"firend_uid"` // 好友的 UID
	UserUID   uint `json:"user_uid"`   // 自身的 UID
	User      User `json:"user"       gorm:"references:UID"`
	Firend    User `json:"firend"     gorm:"references:UID"`
}

// 用户表
type User struct {
	BaseID
	Email       string `gorm:"not null"                       json:"email"`
	Username    string `gorm:"not null"                       json:"username"`
	Password    string `gorm:"not null"                       json:"-"`
	Description string `                                      json:"description"`
	Avatar      string `                                      json:"avatar"`
	Gender      int8   `                                      json:"gender"`
	UID         uint   `gorm:"unique;not null;default:100000" json:"uid"`
	LastLogin   int64  `                                      json:"last_login,omitempty"`
	Onlie       bool   `                                      json:"online"`
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
	tx.Model(&User{}).Unscoped().Order("uid desc").First(user, "")
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
