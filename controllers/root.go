package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"rua.plus/gymo/utils"
)

type RootController struct{}

func (root RootController) Root(c *gin.Context) {
	resp := &utils.BasicRes{
		Status:  "ok",
		Message: "hello gymo",
	}

	c.JSON(http.StatusOK, resp)
}
