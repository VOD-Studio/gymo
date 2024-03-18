package controllers

import (
	"errors"
	"fmt"
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
		utils.FailedAndReturn(c, resp, http.StatusBadRequest, err.Error())
		return
	}

	// check is self
	if info.Uid == u.UID {
		utils.FailedAndReturn(
			c,
			resp,
			http.StatusUnprocessableEntity,
			"cannot make firend with self",
		)
		return
	}

	// find target user
	firend := &models.User{}
	dbRes := contacts.Db.Model(firend).First(firend, "uid = ?", info.Uid)
	if dbRes.Error != nil {
		if errors.Is(dbRes.Error, gorm.ErrRecordNotFound) {
			utils.FailedAndReturn(
				c,
				resp,
				http.StatusUnprocessableEntity,
				"target user not exist",
			)
			return
		}
		utils.FailedAndReturn(
			c,
			resp,
			http.StatusInternalServerError,
			dbRes.Error.Error(),
		)
		return
	}

	// check is already in contect
	contact := &models.Contact{}
	dbRes = contacts.Db.Model(contact).
		First(contact, "user_id = ? AND firend_id = ?", u.UID, info.Uid)
	if dbRes.Error != nil && !errors.Is(dbRes.Error, gorm.ErrRecordNotFound) {
		utils.FailedAndReturn(
			c,
			resp,
			http.StatusInternalServerError,
			dbRes.Error.Error(),
		)
		return
	}
	if dbRes.RowsAffected != 0 {
		utils.FailedAndReturn(
			c,
			resp,
			http.StatusConflict,
			"target user is already firend",
		)
		return
	}

	// save to request
	firendReq := &models.FirendRequest{}
	dbRes = contacts.Db.Model(firendReq).
		First(firendReq, "from_user_id = ? AND to_user_id = ?", u.ID, firend.ID)
	if dbRes.Error != nil && !errors.Is(dbRes.Error, gorm.ErrRecordNotFound) {
		utils.FailedAndReturn(
			c,
			resp,
			http.StatusInternalServerError,
			dbRes.Error.Error(),
		)
		return
	}
	if dbRes.RowsAffected != 0 {
		utils.FailedAndReturn(
			c,
			resp,
			http.StatusConflict,
			fmt.Sprintf("already sent a request to user %d", firend.UID),
		)
		return
	}
	/* firendReq.FromUserID = u.ID */
	/* firendReq.ToUserID = firend.ID */
	firendReq.FromUser = *u
	firendReq.ToUser = *firend
	dbRes = contacts.Db.Save(firendReq)
	if dbRes.Error != nil {
		utils.FailedAndReturn(
			c,
			resp,
			http.StatusInternalServerError,
			dbRes.Error.Error(),
		)
		return
	}

	resp.Status = "ok"
	resp.Message = ""
	c.JSON(http.StatusOK, resp)
}

// 获取当前账号的好友列表
func (contacts Contacts) FirendList(c *gin.Context) {
	// response
	resp := &utils.BasicRes{}
	u := utils.GetContextUser(c, resp)

	var list = []models.Contact{}
	dbRes := contacts.Db.Model(&models.Contact{}).
		Find(&list, "user_id = ?", u.ID)
	if dbRes.Error != nil {
		utils.FailedAndReturn(
			c,
			resp,
			http.StatusInternalServerError,
			dbRes.Error.Error(),
		)
		return
	}
	if dbRes.RowsAffected == 0 {
		utils.FailedAndReturn(
			c,
			resp,
			http.StatusUnprocessableEntity,
			"firend list is empty",
		)
		return
	}

	resp.Status = "ok"
	resp.Data = list
	c.JSON(http.StatusOK, resp)
}

type RequestListInfo struct {
	// 是否是自己收到的好友请求
	Received bool `form:"received"`
}

// 当前账户的好友请求列表
func (contacts Contacts) RequestList(c *gin.Context) {
	// response
	resp := &utils.BasicRes{}
	u := utils.GetContextUser(c, resp)

	var info RequestListInfo
	if err := c.ShouldBindWith(&info, binding.Query); err != nil {
		utils.FailedAndReturn(
			c,
			resp,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	var list = []models.FirendRequest{}
	var condition string
	if info.Received {
		// 是自身收到的好友请求，则 to_user_id 为自身 ID
		condition = "to_user_id = ?"
	} else {
		// 否则 from_user_id 为自身 ID
		condition = "from_user_id = ?"
	}
	dbRes := contacts.Db.Preload("FromUser").Preload("ToUser").Find(&list, condition, u.ID)
	if dbRes.Error != nil {
		utils.FailedAndReturn(
			c,
			resp,
			http.StatusInternalServerError,
			dbRes.Error.Error(),
		)
		return
	}
	if dbRes.RowsAffected == 0 {
		utils.FailedAndReturn(
			c,
			resp,
			http.StatusUnprocessableEntity,
			"firend request is empty",
		)
		return
	}

	resp.Status = "ok"
	resp.Data = list
	c.JSON(http.StatusOK, resp)
}

type AcceptRequestInfo struct {
	ID     uint `json:"id"     binding:"required"`
	Accept bool `json:"accpet" binding:"required"`
}

// 接受好友请求的相关信息
func (contacts Contacts) AcceptRequest(c *gin.Context) {
	// response
	resp := &utils.BasicRes{}
	u := utils.GetContextUser(c, resp)

	info := &AcceptRequestInfo{}
	if err := c.ShouldBindWith(&info, binding.JSON); err != nil {
		utils.FailedAndReturn(c, resp, http.StatusBadRequest, err.Error())
		return
	}

	// find target reqeust
	firendReq := &models.FirendRequest{}
	dbRes := contacts.Db.First(firendReq, "id = ? AND to_user_id = ?", info.ID, u.ID)
	if dbRes.Error != nil {
		if errors.Is(dbRes.Error, gorm.ErrRecordNotFound) {
			utils.FailedAndReturn(
				c,
				resp,
				http.StatusUnprocessableEntity,
				fmt.Sprintf("firend request %d not exist", info.ID),
			)
			return
		} else {
			utils.FailedAndReturn(
				c,
				resp,
				http.StatusInternalServerError,
				dbRes.Error.Error(),
			)
			return
		}
	}

	// check it's already responsed
	if firendReq.Accepted != 0 {
		utils.FailedAndReturn(
			c,
			resp,
			http.StatusBadRequest,
			fmt.Sprintf("request %d already responsed", info.ID),
		)
		return
	}

	if info.Accept {
		firendReq.Accepted = 1
	} else {
		firendReq.Accepted = 2
	}
	dbRes = contacts.Db.Save(firendReq)
	if dbRes.Error != nil {
		utils.FailedAndReturn(
			c,
			resp,
			http.StatusInternalServerError,
			dbRes.Error.Error(),
		)
		return
	}

	resp.Status = "ok"
	resp.Message = "accepted"
	c.JSON(http.StatusOK, resp)
}
