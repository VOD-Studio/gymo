package main

import (
	"log"
	"rua.plus/gymo/db"

	"github.com/joho/godotenv"

	"rua.plus/gymo/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("env file not found")
	}

	db.InitPostgres()
	server.InitServer()
}
