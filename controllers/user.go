package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct{}

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

	c.JSON(http.StatusOK, UserResponse{Status: "Ok", Username: userInfo.Username})
}
