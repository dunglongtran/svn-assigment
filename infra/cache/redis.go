package cache

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var ctx = context.Background()

func NewRedisClient() *redis.Client {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve Redis connection details from environment variables
	redisAddr := os.Getenv("REDIS_ADDR")     // Example: "localhost:6379"
	redisPassword := os.Getenv("REDIS_PASS") // Example: "", if no password is set
	redisDB := os.Getenv("REDIS_DB")         // Example: "0", use default DB if not specified

	db, err := strconv.Atoi(redisDB)
	if err != nil {
		log.Fatalf("Invalid Redis DB: %v", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,     // Redis address
		Password: redisPassword, // Redis password
		DB:       db,            // Redis DB
	})

	return client
}
