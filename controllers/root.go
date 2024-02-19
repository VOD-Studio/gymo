package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type RootController struct{}

type RootStatus struct {
	Status string `json:"status"`
}

func (root RootController) Root(c *gin.Context) {
	status := &RootStatus{
		Status: "ok",
	}
	c.JSON(http.StatusOK, status)
}
