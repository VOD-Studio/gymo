package main

import (
	"github.com/joho/godotenv"
	"log"
	"rua.plus/gymo/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("env file not found")
	}

	server.InitServer()
}
