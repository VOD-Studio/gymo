package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

var pool *pgxpool.Pool

func InitPostgres() {
	var err error
	pool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}
}
