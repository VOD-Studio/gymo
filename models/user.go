package models

import "gorm.io/gorm"

type User struct {
	gorm.Model `json:"gorm.Model"`
	Username   string `json:"username"`
	Password   string `json:"password"`
}

func (user User) GetSingle(username string) {

}
