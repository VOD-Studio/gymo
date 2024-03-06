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

// 从 gin 上下文中获取当前登录的用户
func GetContextUser(c *gin.Context, resp *BasicRes) *models.User {
	var u *models.User
	current, ok := c.Get("user")
	if !ok {
		FailedAndReturn(c, resp, http.StatusInternalServerError, "parse token failed")
		return nil
	}
	u = current.(*models.User)
	return u
}

// 格式化失败响应，并通过 `c.JSON` 返回
func FailedAndReturn(c *gin.Context, resp *BasicRes, code int, message string) {
	resp.Status = "error"
	resp.Message = message
	c.JSON(code, resp)
}
