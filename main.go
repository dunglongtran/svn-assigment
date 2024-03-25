package main

import (
	"SVN-interview/cmd/api"
	"SVN-interview/infra/cache"
	"SVN-interview/infra/db"
	"SVN-interview/internal/common"
	"github.com/joho/godotenv"
	"os"
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
	// Setup Redis
	redisClient := cache.NewRedisClient()

	// Define Context
	appCtx := &common.AppContext{
		DB:    dbInstance,
		Cache: redisClient,
	}

	router := api.SetupRouter(appCtx)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to port 8080 if PORT is not set
	}

	// Start the server on the specified port
	router.Run(":" + port)
	//router := gin.Default()
	////router.GET("/get_histories", handlers.GetHistoriesHandler)
	//router.GET("/ping", func(context *gin.Context) {
	//	context.JSON(http.StatusOK, gin.H{
	//		"message": "OK",
	//	})
	//})
	//router.Run(":8080")
}
