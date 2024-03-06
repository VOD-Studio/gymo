package controllers

import (
	"fmt"
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

// 向指定的用户发送好友请求
// 发送后将保存到 `firend_request` 表中
// 同时向对方发送通知
// TODO: 给对方用户发送通知
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

	// check is self
	if info.Uid == u.UID {
		resp.Status = "error"
		resp.Message = "cannot make firend with self"
		c.JSON(http.StatusUnprocessableEntity, resp)
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
	dbRes = contacts.Db.Model(contact).
		Find(contact, "user_uid = ? AND firend_uid = ?", u.UID, info.Uid)
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

	// save to request
	firendReq := &models.FirendRequest{}
	dbRes = contacts.Db.Model(firendReq).
		Find(firendReq, "from_user_uid = ? AND to_user_uid = ?", u.UID, info.Uid)
	if dbRes.Error != nil {
		resp.Status = "error"
		resp.Message = dbRes.Error.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	if dbRes.RowsAffected != 0 {
		resp.Status = "error"
		resp.Message = fmt.Sprintf("already sent a request to user %d", firend.UID)
		c.JSON(http.StatusUnprocessableEntity, resp)
		return
	}
	firendReq.FromUserUID = u.UID
	firendReq.ToUserUID = firend.UID
	contacts.Db.Save(firendReq)

	// save
	/* contact.UserUID = u.UID */
	/* contact.Firend = firend.UID */
	/* dbRes = contacts.Db.Save(contact) */
	/* if dbRes.Error != nil { */
	/* 	resp.Status = "error" */
	/* 	resp.Message = dbRes.Error.Error() */
	/* 	c.JSON(http.StatusInternalServerError, resp) */
	/* 	return */
	/* } */

	resp.Status = "ok"
	resp.Message = ""
	c.JSON(http.StatusOK, resp)
}

func (contacts Contacts) CheckRequest(c *gin.Context) {
	// response
	resp := &utils.BasicRes{}
	u := utils.GetContextUser(c, resp)

	log.Println(u)
}
