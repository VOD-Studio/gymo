package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"rua.plus/gymo/controllers"
	"rua.plus/gymo/db"
	"rua.plus/gymo/middlewares"
)

func InitRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	v1 := router.Group("/v1")

	root := controllers.RootController{}
	v1.GET("/", root.Root)

	user := controllers.User{
		Db: db.Db,
	}
	v1.GET("/user", user.GetUser)      // query single user by query
	v1.POST("/register", user.AddUser) // register account
	v1.POST("/login", user.Login)      // login
	v1.Use(middlewares.TokenAuth())
	v1.Use(middlewares.TokenTimeAuth(db.Db))
	v1.PATCH("/user", user.ModifyUser) // modify user infomation
	v1.POST("/user", user.UserSelf)    // get current logged in user infomation
	v1.DELETE("/user", user.Delete)    // cancel this account

	return router
}
