package main

import (
	"log"

	"github.com/joho/godotenv"

	"rua.plus/gymo/db"
	"rua.plus/gymo/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("env file not found")
	}

	db.InitPostgres()
	err = db.InitRedis()
	if err != nil {
		log.Fatalf("Cannot connect to redis %s\n", err)
	}
	server.InitServer()
}
