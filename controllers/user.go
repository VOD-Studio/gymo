package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"rua.plus/gymo/models"
)

type User struct {
	Db *gorm.DB
}

type UserResponse struct {
	Message  string `json:"message"`
	Status   string `json:"status"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserQuery struct {
	Username string `form:"username"`
}

func (user User) GetUser(c *gin.Context) {
	var userInfo UserQuery
	if err := c.ShouldBindQuery(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userInfo.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is empty"})
		return
	}

	res := &UserResponse{
		Status: "ok",
	}

	var u models.User
	if err := u.GetSingle(userInfo.Username, user.Db); err != nil {
		res.Status = "error"
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res.Message = "user not found"
			c.JSON(http.StatusOK, &res)
			return
		} else {
			res.Message = err.Error()
			c.JSON(http.StatusInternalServerError, &res)
			return
		}
	}

	res.Username = u.Username
	res.Email = u.Email
	c.JSON(http.StatusOK, &res)
}

type UserJson struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (user User) AddUser(c *gin.Context) {
	var userInfo UserJson
	if err := c.ShouldBindJSON(&userInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := &UserResponse{
		Status:   "ok",
		Username: "",
	}

	u := &models.User{
		Username: userInfo.Username,
		Password: userInfo.Password,
		Email:    userInfo.Email,
	}
	if err := u.Create(user.Db); err != nil {
		res.Message = err.Error()
		res.Status = "error"
		if errors.Is(err, models.UserAlreadyExisty) {
			c.JSON(http.StatusConflict, &res)
			return
		} else {
			c.JSON(http.StatusInternalServerError, &res)
			return
		}
	}

	res.Username = u.Username
	res.Email = u.Email

	c.JSON(http.StatusOK, &res)
}
