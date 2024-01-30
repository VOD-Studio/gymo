package main

import (
	"log"

	"github.com/joho/godotenv"

	"rua.plus/gymo/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("env file not found")
	}

	server.InitServer()
}
