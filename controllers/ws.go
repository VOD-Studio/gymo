package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"

	"rua.plus/gymo/utils"
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

	// response
	resp := &utils.BasicRes{}
	u := utils.GetContextUser(c, resp)
	log.Println(u.Username)

	defer ctx.Close()
	for {
		// message type:
		// 1: text
		// 2: binary
		mt, message, err := ctx.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s message type %d \n", message, mt)
		err = ctx.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}
