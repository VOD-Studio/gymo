package server

import (
	"net/http"

	"rua.plus/gymo/controllers"
)

func NewRouter() {
	root := new(controllers.Root)
	http.HandleFunc("/", root.GetRoot)

	user := new(controllers.User)
	user.Path = "/xfy"
	http.HandleFunc(user.Path, user.GetUser)
}