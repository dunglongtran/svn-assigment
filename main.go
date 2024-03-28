package main

import (
	"SVN-interview/cmd/api"
	"SVN-interview/cmd/cron"
	"SVN-interview/infra/cache"
	"SVN-interview/infra/db"
	"SVN-interview/internal/common"
	"SVN-interview/internal/entities"
	cache2 "SVN-interview/pkg/cache"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"

	// Import your generated docs package
	_ "SVN-interview/docs"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
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
	err = cache2.InitializeAllCoinsDates(appCtx.DB, appCtx.Cache)

	if err != nil {
		panic("Failed to init to cache")
	}
	router := api.SetupRouter(appCtx)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if PORT is not set
	}

	go cron.StartCronJob(ctx, appCtx)

	//// Start the server on the specified port
	//router.Run(":" + port)

	// Start the server and wait interrupt signal
	go func() {
		// Start the server on the specified port
		if err := router.Run(":" + port); err != nil {
			fmt.Println("Error starting server:", err)
			stop()
		}
	}()

	<-ctx.Done()
	fmt.Println("Server stopping.")

}
