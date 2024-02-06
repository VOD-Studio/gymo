package server

import (
	"log"
	"os"
)

func InitServer() {
	r := InitRouter()
	err := r.Run("0.0.0.0:" + os.Getenv("PORT"))
	if err != nil {
		log.Fatal("start server failed", err)
	}
}
