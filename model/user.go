package model

import (
	"context"
	"log"
	"rua.plus/gymo/db"
)

type User struct{}

func (user User) GetSingel() {
	rows, err := db.Pool.Query(context.Background(), "SELECT * FROM users LIMIT 1;")
	if err != nil {
		return
	}
	log.Println(rows)
}
