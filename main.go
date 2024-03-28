package main

import (
	"SVN-interview/cmd/api"
	"SVN-interview/infra/cache"
	"SVN-interview/infra/db"
	"SVN-interview/internal/common"
	"SVN-interview/internal/entities"
	cache2 "SVN-interview/pkg/cache"
	"github.com/joho/godotenv"
	"os"
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
	err := cache2.InitializeAllCoinsDates(appCtx.DB, appCtx.Cache)
	if err != nil {
		panic("Failed to init to cache")
	}
	router := api.SetupRouter(appCtx)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if PORT is not set
	}

	// Start the server on the specified port
	router.Run(":" + port)

}
