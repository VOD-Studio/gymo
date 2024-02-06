package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct{}

type UserResponse struct {
	Username string `json:"username"`
}

func (user User) GetUser(c *gin.Context) {
	username := c.Param("username")

	if username == "" {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, UserResponse{Username: username})
}
