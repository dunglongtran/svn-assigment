package main

import (
	"SVN-interview/cmd/process"
	"SVN-interview/infra/cache"
	"SVN-interview/infra/db"
	"SVN-interview/internal/common"
	"SVN-interview/internal/entities"
	"github.com/joho/godotenv"
	"time"

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
	idCoin := "bitcoin"
	layout := "2006-01-02"
	startTime, _ := time.Parse(layout, "2024-02-20")
	startTs := startTime.Unix()
	endTime, _ := time.Parse(layout, "2024-03-20")
	endTs := endTime.Unix()

	process.GenerateOHLCPriceData(appCtx.DB, idCoin, startTs, endTs, common.Period1H)
}
