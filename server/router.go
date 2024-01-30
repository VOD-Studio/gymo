package server

import (
	"net/http"

	"rua.plus/gymo/controllers"
)

func NewRouter() {
	root := new(controllers.Root)
	http.HandleFunc("/", root.GetRoot)
}
