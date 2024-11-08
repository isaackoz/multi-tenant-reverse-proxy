package db

import (
	"context"
	"database/sql"
	"log"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

var RedisClient *redis.Client
var PostgresDB *sql.DB
var ctx = context.Background()

func InitRedis(redisAddress string) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	log.Println("Connected to Redis")
	return nil
}

func InitPostgres(postgresURI string) error {
	var err error
	PostgresDB, err = sql.Open("postgres", postgresURI)
	if err != nil {
		return err
	}

	if err = PostgresDB.Ping(); err != nil {
		return err
	}

	log.Println("Connected to PostgreSQL")
	return nil
}
