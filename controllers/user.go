package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"rua.plus/gymo/models"
	"rua.plus/gymo/utils"
)

type User struct {
	Db *gorm.DB
}

// 查询用户
type UserQuery struct {
	Email string `form:"email" binding:"required"`
}

// 通过 email 查询用户
// 仅支持 query
func (user User) GetUser(c *gin.Context) {
	// response
	resp := &utils.BasicRes{}

	var userInfo UserQuery
	if err := c.ShouldBindQuery(&userInfo); err != nil {
		resp.Status = "error"
		resp.Message = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if userInfo.Email == "" {
		resp.Status = "error"
		resp.Message = "email is empty"
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	u := &models.User{}
	res := user.Db.Model(u).Find(u, "email = ?", userInfo.Email)
	if res.Error != nil {
		resp.Status = "error"
		resp.Message = res.Error.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	if res.RowsAffected == 0 {
		resp.Status = "error"
		resp.Message = "user not exist"
		c.JSON(http.StatusNoContent, resp)
		return
	}

	resp.Status = "ok"
	resp.Data = u
	c.JSON(http.StatusOK, resp)
}

// 用户注册
type UserJson struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"    binding:"required"`
}

// 添加用户
// 仅支持 json body
func (user User) AddUser(c *gin.Context) {
	// response
	resp := &utils.BasicRes{}

	var userInfo UserJson
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		resp.Status = "error"
		resp.Message = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	u := &models.User{
		Username: userInfo.Username,
		Password: userInfo.Password,
		Email:    userInfo.Email,
	}

	res := user.Db.Model(u).Where("email = ?", u.Email).FirstOrCreate(&u)
	if res.Error != nil {
		resp.Status = "error"
		resp.Message = res.Error.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	if res.RowsAffected == 0 {
		resp.Status = "error"
		resp.Message = "user already exist"
		c.JSON(http.StatusConflict, resp)
		return
	}

	resp.Status = "ok"
	resp.Data = u
	c.JSON(http.StatusOK, resp)
}

type UserModify struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (user User) ModifyUser(c *gin.Context) {
	// response
	resp := &utils.BasicRes{}

	var u *models.User
	current, ok := c.Get("user")
	if !ok {
		resp.Status = "error"
		resp.Message = "parse token failed"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	u = current.(*models.User)

	userInfo := &UserModify{}
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		resp.Status = "error"
		resp.Message = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if userInfo.Username != "" {
		u.Username = userInfo.Username
	}
	if userInfo.Email != "" {
		u.Email = userInfo.Email
	}
	if userInfo.Password != "" {
		u.Password = userInfo.Password
		u.HashPassword()
	}

	res := user.Db.Save(u)
	if res.Error != nil {
		resp.Status = "error"
		resp.Message = res.Error.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Status = "ok"
	resp.Data = u
	c.JSON(http.StatusOK, resp)

}

// 用户登录 json
type UserLogin struct {
	Email    string `json:"email"    binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginResponse struct {
	*models.User
	Token string `json:"token"`
}

// 用户登录
// 仅支持 json
func (user User) Login(c *gin.Context) {
	// response
	resp := &utils.BasicRes{}

	var userInfo UserLogin
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		resp.Status = "error"
		resp.Message = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// query the user
	u := &models.User{}
	dbResult := user.Db.Model(&models.User{}).Find(&u, "email = ?", userInfo.Email)
	if dbResult.Error != nil {
		resp.Status = "error"
		resp.Message = dbResult.Error.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	if dbResult.RowsAffected == 0 {
		resp.Status = "error"
		resp.Message = "user not exist"
		c.JSON(http.StatusNoContent, resp)
		return
	}

	// check the password
	if err := models.CheckPasswordHash(userInfo.Password, u.Password); err != nil {
		resp.Status = "error"
		resp.Message = "password not correct"
		c.JSON(http.StatusUnauthorized, resp)
		return
	}

	// generate token
	lastLogin := time.Now().Unix()
	token, err := utils.GenerateToken(int(u.ID), lastLogin)
	if err != nil {
		resp.Status = "error"
		resp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// update last login
	u.LastLogin = lastLogin
	user.Db.Save(u)

	resp.Status = "ok"
	resp.Data = &LoginResponse{
		User:  u,
		Token: token,
	}
	c.JSON(http.StatusOK, resp)
}

// 当前登录的用户信息
// 通过 Token 获取
func (user User) UserSelf(c *gin.Context) {
	// response
	resp := &utils.BasicRes{}

	var u *models.User
	current, ok := c.Get("user")
	if !ok {
		resp.Status = "error"
		resp.Message = "parse token failed"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	u = current.(*models.User)

	resp.Status = "ok"
	resp.Data = u
	c.JSON(http.StatusOK, resp)
	return
}

// 删除当前用户
func (user User) Delete(c *gin.Context) {
	// response
	resp := &utils.BasicRes{}

	var u *models.User
	current, ok := c.Get("user")
	if !ok {
		resp.Status = "error"
		resp.Message = "parse token failed"
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	u = current.(*models.User)

	res := user.Db.Model(&models.User{}).Delete(u, "email = ?", u.Email)
	if res.Error != nil {
		resp.Status = "error"
		resp.Message = res.Error.Error()
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	msg := fmt.Sprintf("account %s has been deleted", u.Email)

	resp.Status = "ok"
	resp.Message = msg
	c.JSON(http.StatusOK, resp)
}
