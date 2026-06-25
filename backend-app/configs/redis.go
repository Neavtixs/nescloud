package configs

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

func NewAccess() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	if strings.TrimSpace(host) == "" {
		host = "localhost"
	}

	port := os.Getenv("REDIS_PORT")
	if strings.TrimSpace(port) == "" {
		port = "6379"
	}

	password := os.Getenv("REDIS_PASS")

	db := 0
	if rawDB := os.Getenv("REDIS_DB"); strings.TrimSpace(rawDB) != "" {
		parsedDB, err := strconv.Atoi(rawDB)
		if err != nil {
			panic(err)
		}
		db = parsedDB
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return redisClient
}
