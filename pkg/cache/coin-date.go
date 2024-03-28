package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
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

	// Save into Redis
	_, err = redisClient.HSet(ctx, "coin_dates:"+idCoin, "recentDate", recentDate, "furthestDate", furthestDate,
		"updatedAt", time.Now().Format(time.RFC3339), "latestCheck", time.Now().Format(time.RFC3339)).Result()

	return err
}
func UpdateFurthestDateInRedis(redisClient *redis.Client, idCoin string, startDate int64) error {
	ctx := context.Background()
	// startDate is UNIX timestamp
	_, err := redisClient.HSet(ctx, "coin_dates:"+idCoin, "furthestDate", startDate, "updatedAt", time.Now().Format(time.RFC3339)).Result()
	return err
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
		return "", err
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
