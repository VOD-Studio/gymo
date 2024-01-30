package controllers

import (
	"io"
	"log"
	"net/http"
)

type Root struct{}

func (root Root) GetRoot(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v / request", r.Method)
	io.WriteString(w, "Hello Gymo!")
}
