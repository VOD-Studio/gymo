package controllers

import (
	"encoding/json"
	"log"
	"net/http"
)

type User struct {
	Path string
}

func (u User) GetUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v request", r.Method, u.Path)

	user := map[string]any{
		"name": "xfy",
		"age":  14,
		"test": nil,
	}
	res, err := json.Marshal(user)
	if err != nil {
		log.Printf("parse json failed %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
