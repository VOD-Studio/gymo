package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type User struct{}

func (user User) GetUser(c *gin.Context) {
	username := c.Param("username")

	if username == "" {
		c.Status(http.StatusNotFound)
		return
	}
}
