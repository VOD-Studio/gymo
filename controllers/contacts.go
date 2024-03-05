package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"

	"rua.plus/gymo/models"
	"rua.plus/gymo/utils"
)

type Contacts struct {
	Db *gorm.DB
}

type MakeFirendJson struct {
	Uid uint `json:"uid" binding:"required"`
}

func (contacts Contacts) MakeFirend(c *gin.Context) {
	// response
	resp := &utils.BasicRes{}
	u := utils.GetContextUser(c, resp)

	var info MakeFirendJson
	if err := c.ShouldBindWith(&info, binding.JSON); err != nil {
		resp.Status = "error"
		resp.Message = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// find target user
	firend := &models.User{}
	dbRes := contacts.Db.Model(firend).Find(firend, "uid = ?", info.Uid)
	if dbRes.Error != nil {
		resp.Status = "error"
		resp.Message = dbRes.Error.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	if dbRes.RowsAffected == 0 {
		resp.Status = "error"
		resp.Message = "target user not exist"
		c.JSON(http.StatusUnprocessableEntity, resp)
		return
	}

	// check is already in contect
	contact := &models.Contact{}
	dbRes = contacts.Db.Model(contact).Find(contact, "user_uid = ? AND firend = ?", u.UID, info.Uid)
	if dbRes.Error != nil {
		resp.Status = "error"
		resp.Message = dbRes.Error.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	if dbRes.RowsAffected != 0 {
		resp.Status = "error"
		resp.Message = "target user is already firend"
		c.JSON(http.StatusUnprocessableEntity, resp)
		return
	}

	// save
	contact.UserUID = u.UID
	contact.Firend = firend.UID
	dbRes = contacts.Db.Save(contact)
	if dbRes.Error != nil {
		resp.Status = "error"
		resp.Message = dbRes.Error.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	log.Println(contact, dbRes.RowsAffected)
}
