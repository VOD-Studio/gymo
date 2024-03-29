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
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{
		"Content-Type",
		"Content-Length",
		"Accept-Encoding",
		"X-CSRF-Token",
		"Authorization",
		"accept",
		"origin",
		"Cache-Control",
		"X-Requested-With",
	}
	router.Use(cors.New(config))

	v1 := router.Group("/v1")

	root := controllers.RootController{}
	v1.GET("/", root.Root)

	// user
	user := controllers.User{
		Db: db.Db,
	}
	v1.GET("/user", user.GetUser)      // query single user by query
	v1.POST("/register", user.AddUser) // register account
	v1.POST("/login", user.Login)      // login

	// authorization
	{
		v1.Use(middlewares.TokenAuth())
		v1.Use(middlewares.TokenTimeAuth(db.Db))
		v1.PATCH("/user", user.ModifyUser) // modify user infomation
		v1.POST("/user", user.UserSelf)    // get current logged in user infomation
		v1.DELETE("/user", user.Delete)    // cancel this account

		// contacts
		contacts := controllers.Contacts{
			Db: db.Db,
		}
		v1.POST("/follow", contacts.FollowUser) // follow a user
		v1.POST("/make_firend", contacts.MakeFirend)
		v1.GET("/firends", contacts.FirendList)
		v1.GET("/firend_requests", contacts.RequestList)
		v1.PATCH("/firend_request", contacts.AcceptRequest)

		// websocket
		ws := controllers.WS{
			Db: db.Db,
		}
		v1.GET("/ws", ws.Connect)
	}

	return router
}
