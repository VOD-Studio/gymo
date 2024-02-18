package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"rua.plus/gymo/models"
)

type User struct {
	Db *gorm.DB
}

type UserResponse struct {
	Status   string `json:"status"`
	Username string `json:"username"`
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
		Status:   "ok",
		Username: "",
	}

	var u models.User
	if err := u.GetSingle(userInfo.Username, user.Db); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			res.Status = "user not found"
			c.JSON(http.StatusOK, &res)
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, &res)
}
