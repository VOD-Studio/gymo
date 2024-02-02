package controllers

import (
	"io"
	"log"
	"net/http"
)

type Root struct{}

func (root Root) GetRoot(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v / request", r.Method)
	_, err := io.WriteString(w, "Hello Gymo!")
	if err != nil {
		log.Printf("write to client failed %v", err)
		return
	}
}
