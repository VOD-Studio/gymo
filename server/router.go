package server

import (
	"github.com/gin-gonic/gin"
	"rua.plus/gymo/controllers"
	"rua.plus/gymo/db"
)

func InitRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1 := router.Group("/v1")

	root := controllers.RootController{}
	v1.GET("/", root.Root)

	user := controllers.User{
		Db: db.Db,
	}
	v1.GET("/user/", user.GetUser)

	return router
}
