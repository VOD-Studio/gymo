package server

import (
	"log"
	"net/http"
	"os"
)

func InitServer() {
	log.Println("init server")
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	NewRouter()
	log.Println("server listening on", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
