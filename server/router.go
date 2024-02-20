package server

import (
	"github.com/gin-gonic/gin"

	"rua.plus/gymo/controllers"
	"rua.plus/gymo/db"
	"rua.plus/gymo/middlewares"
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
	v1.POST("/user/", user.AddUser)

	v1.POST("/login/", user.Login)
	v1.Use(middlewares.TokenAuth())
	v1.PATCH("/user/", user.ModifyUser)

	return router
}
