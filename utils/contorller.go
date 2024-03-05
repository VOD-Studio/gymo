package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"rua.plus/gymo/models"
)

type BasicRes struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func GetContextUser(c *gin.Context, resp *BasicRes) *models.User {
	var u *models.User
	current, ok := c.Get("user")
	if !ok {
		resp.Status = "error"
		resp.Message = "parse token failed"
		c.JSON(http.StatusInternalServerError, resp)
		return nil
	}
	u = current.(*models.User)
	return u
}
