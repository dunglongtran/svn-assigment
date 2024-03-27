package coingecko

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func fetchURL(url, apiKey string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if apiKey != "" {
		// Chỉ thêm header này nếu apiKey được cung cấp
		req.Header.Add("x-cg-demo-api-key", apiKey)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if !strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		return "", fmt.Errorf("unexpected content type, want application/json")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func fetchGecko(endPoint string) (string, error) {
	url := os.Getenv("COINGK_BASE_URL") + endPoint
	apiKey := os.Getenv("COINGK_API_KEY")

	body, err := fetchURL(url, apiKey)
	if err != nil {
		return "", fmt.Errorf("error fetching URL: %v", err) // Trả về lỗi thay vì in lỗi
	}

	return body, nil
}
