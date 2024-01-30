package server

import (
	"log"
	"os"
)

func InitServer() {
	log.Println("init server")
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}
}
