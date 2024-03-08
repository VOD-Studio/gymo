package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{} // use default option

type WS struct {
	Db *gorm.DB
}

func (ws WS) Connect(c *gin.Context) {
	w, r := c.Writer, c.Request
	ctx, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer ctx.Close()
	for {
		mt, message, err := ctx.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s \n", message)
		err = ctx.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}
