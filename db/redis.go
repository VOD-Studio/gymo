package db

import (
	"context"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var Rdb *redis.Client

func InitRedis() error {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return err
	}
	Rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWD"),
		DB:       db,
	})
	return nil
}
