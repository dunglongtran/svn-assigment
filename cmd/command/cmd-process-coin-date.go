package main

import (
	"SVN-interview/infra/cache"
	"SVN-interview/infra/db"
	"SVN-interview/internal/common"
	"SVN-interview/internal/entities"
	cache2 "SVN-interview/pkg/cache"
	"fmt"
	"github.com/joho/godotenv"
	// Import your generated docs package
	_ "SVN-interview/docs"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load env file")
	}
	// Initialize database
	dbInstance, err := db.Initialize()
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto migrate
	dbInstance.AutoMigrate(&entities.CoinOHLC{})
	dbInstance.AutoMigrate(&entities.CoinPrice{})

	// Setup Redis
	redisClient := cache.NewRedisClient()

	// Define Context
	appCtx := &common.AppContext{
		DB:    dbInstance,
		Cache: redisClient,
	}
	//cache2.InitializeAllCoinsDates(appCtx.DB, appCtx.Cache)
	//cache2.InitializeCoinDatesInRedis(appCtx.DB, appCtx.Cache, "bitcoin")
	coinDates, err := cache2.ReadCoinDatesFromRedis(appCtx.Cache, "bitcoin")
	if err != nil {
		fmt.Println("Error reading coin dates from Redis:", err)
		return
	}
	fmt.Println("Coin Dates:", coinDates["latestCheck"])
}
