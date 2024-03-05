package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"rua.plus/gymo/utils"
)

type Contacts struct {
	Db *gorm.DB
}

func (contacts Contacts) MakeFirend(c *gin.Context) {
	// response
	resp := &utils.BasicRes{}

	u := utils.GetContextUser(c, resp)

	log.Println(u)
}
