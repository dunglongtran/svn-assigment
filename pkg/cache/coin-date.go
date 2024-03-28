package cache

import (
	"SVN-interview/internal/entities"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"time"
)

func InitializeCoinDatesInRedis(db *gorm.DB, redisClient *redis.Client, idCoin string) error {
	var recentDate, furthestDate int64
	ctx := context.Background()
	// Get recentDate from coin_price
	var recentPrice struct{ Time int64 }
	if err := db.Table("coin_price").Select("MIN(time) as time").Where("\"idCoin\" = ?", idCoin).Scan(&recentPrice).Error; err != nil || recentPrice.Time == 0 {
		recentDate = time.Now().Unix()
	} else {
		recentDate = recentPrice.Time
	}

	// Get furthestDate from coin_ohlc
	var furthestOHLC struct{ Time int64 }
	if err := db.Table("coin_ohlc").Select("MIN(time) as time").Where("\"idCoin\" = ?", idCoin).Scan(&furthestOHLC).Error; err != nil || furthestOHLC.Time == 0 {
		furthestDate = time.Now().Unix()
	} else {
		furthestDate = furthestOHLC.Time / 1000
	}

	// Kiểm tra xem Redis đã có dữ liệu chưa
	exists, err := redisClient.HExists(ctx, "coin_dates:"+idCoin, "recentDate").Result()
	if err != nil || exists {
		// not update if exist
		return nil
	}

	latestCheck := time.Now()
	// Save into Redis
	_, err = redisClient.HSet(ctx, "coin_dates:"+idCoin, "recentDate", recentDate, "furthestDate", furthestDate,
		"updatedAt", time.Now().Format(time.RFC3339), "latestCheck", latestCheck.Format(time.RFC3339)).Result()
	err = AddCoinDataToRedis(redisClient, ctx, "coin_dates_latestCheck", idCoin, float64(latestCheck.Unix()))
	err = AddCoinDataToRedis(redisClient, ctx, "coin_dates_access", idCoin, 0)
	if err != nil {
		return err
	}
	return err
}
func UpdateFurthestDateInRedis(redisClient *redis.Client, idCoin string, startDate int64) error {
	ctx := context.Background()
	furthestDateStr, _ := ReadFurthestDateFromRedis(redisClient, idCoin)
	furthestDate, _ := strconv.ParseInt(furthestDateStr, 10, 64)
	if startDate < furthestDate {
		// startDate is UNIX timestamp
		_, err := redisClient.HSet(ctx, "coin_dates:"+idCoin, "furthestDate", startDate, "updatedAt", time.Now().Format(time.RFC3339)).Result()
		return err
	}
	return nil
}
func UpdateRecentDateInRedis(redisClient *redis.Client, idCoin string, endDate int64) error {
	ctx := context.Background()
	// endDate is UNIX timestamp
	_, err := redisClient.HSet(ctx, "coin_dates:"+idCoin, "recentDate", endDate, "updatedAt", time.Now().Format(time.RFC3339)).Result()
	return err
}
func InitializeAllCoinsDates(db *gorm.DB, redisClient *redis.Client) error {
	var coins []struct{ ID string }
	if err := db.Table("coin_info").Select("id").Scan(&coins).Error; err != nil {
		return err
	}

	for _, coin := range coins {
		if err := InitializeCoinDatesInRedis(db, redisClient, coin.ID); err != nil {
			return err
		}
	}

	return nil
}
func ReadCoinDatesFromRedis(redisClient *redis.Client, idCoin string) (map[string]string, error) {
	ctx := context.Background()
	result, err := redisClient.HGetAll(ctx, "coin_dates:"+idCoin).Result()
	if err != nil {
		return nil, err
	}

	return result, nil
}
func ReadRecentDateFromRedis(redisClient *redis.Client, idCoin string) (string, error) {
	ctx := context.Background()
	result, err := redisClient.HGet(ctx, "coin_dates:"+idCoin, "recentDate").Result()
	if err != nil {
		return "", nil
	}

	return result, nil
}
func ReadFurthestDateFromRedis(redisClient *redis.Client, idCoin string) (string, error) {
	ctx := context.Background()
	result, err := redisClient.HGet(ctx, "coin_dates:"+idCoin, "furthestDate").Result()
	if err != nil {
		return "", err
	}

	return result, nil
}
func UpdateCoinDatesInRedis(redisClient *redis.Client, idCoin string, newFields map[string]interface{}) error {
	ctx := context.Background()

	err := redisClient.HMSet(ctx, "coin_dates:"+idCoin, newFields).Err()
	if err != nil {
		return err
	}

	return nil
}
func FetchCoinDateWithSmallestLatestCheck(rdb *redis.Client, key string) (*entities.CoinDate, error) {
	ctx := context.Background()
	// Fetch all hash values
	result, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	// Initialize minDate to a very high value
	minDate := time.Now().Add(100 * 365 * 24 * time.Hour) // 100 years into the future
	minKey := ""
	for idCoin, latestCheckStr := range result {
		// Parse the latestCheck string as time.Time
		latestCheck, err := time.Parse(time.RFC3339, latestCheckStr)
		if err != nil {
			return nil, err
		}

		// Check if this date is earlier than our current minDate
		if latestCheck.Before(minDate) {
			minDate = latestCheck
			minKey = idCoin
		}
	}

	if minKey == "" {
		return nil, errors.New("No data found")
	}

	// Return the CoinDate with the smallest latestCheck
	return &entities.CoinDate{IDCoin: minKey, LatestCheck: minDate}, nil
}
func FetchAndSortCoinDates(rdb *redis.Client, key string) ([]entities.CoinDate, error) {
	ctx := context.Background()
	// Fetch all hash values
	result, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// Convert to slice of CoinDate for sorting
	var coinDates []entities.CoinDate
	for idCoin, latestCheckStr := range result {
		// Parse the latestCheck as time.Time
		// Assuming latestCheck is stored as a timestamp string in Redis
		latestCheck, err := time.Parse(time.RFC3339, latestCheckStr)
		if err != nil {
			return nil, err
		}

		coinDates = append(coinDates, entities.CoinDate{IDCoin: idCoin, LatestCheck: latestCheck})
	}

	// Sort by LatestCheck
	sort.Slice(coinDates, func(i, j int) bool {
		return coinDates[i].LatestCheck.Before(coinDates[j].LatestCheck)
	})

	return coinDates, nil
}

func AddCoinDataToRedis(rdb *redis.Client, ctx context.Context, key string, idCoin string, score float64) error {
	if err := rdb.ZAdd(ctx, key, &redis.Z{Score: score, Member: idCoin}).Err(); err != nil {
		return err
	}
	return nil
}
func ReadCoinDataToRedis(rdb *redis.Client, ctx context.Context, key string, idCoin string) float64 {
	result, err := rdb.ZScore(ctx, key, idCoin).Result()
	if err != nil {
		fmt.Println("Error read data of", idCoin, "in ", key)
	}
	return result
}
func FetchSortedCoinDates(redisClient *redis.Client, sortedSetKey string, start, stop int64) ([]redis.Z, error) {
	ctx := context.Background()
	results, err := redisClient.ZRangeWithScores(ctx, sortedSetKey, start, stop).Result()
	if err != nil {
		return nil, fmt.Errorf("error retrieving sorted values from %s: %v", sortedSetKey, err)
	}
	return results, nil
}

func IncreaseCoinAccess(rdb *redis.Client, idCoin string) error {
	ctx := context.Background()
	count := ReadCoinDataToRedis(rdb, ctx, "coin_dates_access", idCoin)
	err := AddCoinDataToRedis(rdb, ctx, "coin_dates_access", idCoin, count+1)
	if err != nil {
		fmt.Errorf("Error update access count of ", idCoin)
	}
	return nil
}
