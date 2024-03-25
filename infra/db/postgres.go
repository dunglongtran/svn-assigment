package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Initialize connects to the database with environment variables
func Initialize() (*gorm.DB, error) {
	// Load environment variables
	err := godotenv.Load() // This will load variables from a .env file located in the same directory as the binary or where you run the command. You can specify a path to .env file if it's located elsewhere.
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	// Construct DSN from environment variables
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_PASSWORD"))

	// Open the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Optionally, auto-migrate your schema here
	// db.AutoMigrate(&User{}) // You'd define a User model for this.

	return db, nil
}
