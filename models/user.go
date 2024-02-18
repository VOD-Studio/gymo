package models

import (
	"time"

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

func (u *User) Create(db *gorm.DB) error {
	return db.Create(u).Error
}
