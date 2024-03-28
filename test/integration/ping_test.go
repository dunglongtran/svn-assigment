// ping_test.go
package integration

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestPingEndpoint(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		t.Fatal("PORT environment variable is not set")
	}

	// Tạo URL với port từ biến môi trường
	url := fmt.Sprintf("http://localhost:%s/ping", port)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		t.Fatalf("Failed to send request to %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

}
