package jobs

import (
	"SVN-interview/cmd/process"
	"SVN-interview/internal/common"
	"SVN-interview/internal/entities"
	cache2 "SVN-interview/pkg/cache"
	"SVN-interview/pkg/coingecko"
	"context"
	"fmt"
	"strconv"
	"time"
)

func FillHistoriesPriceJobHandler(appCtx *common.AppContext) {
	sortedCoinDates, err := cache2.FetchSortedCoinDates(appCtx.Cache, "coin_dates_latestCheck", 0, -1)
	if err != nil {
		fmt.Println("Error sorting coin dates from Redis:", err)
		return
	}

	var validCoins []entities.CoinDate

	// verify capacity
	var enoughCoins bool

	for _, coinLatest := range sortedCoinDates {
		if enoughCoins {
			break
		}

		idCoin := coinLatest.Member
		key := fmt.Sprintf("%v", idCoin)
		var info, _ = cache2.ReadCoinDatesFromRedis(appCtx.Cache, key)
		recentDate, _ := strconv.ParseInt(info["recentDate"], 10, 64)
		furthestDate, _ := strconv.ParseInt(info["furthestDate"], 10, 64)
		diffDays := processDates(recentDate, furthestDate)

		if diffDays != -1 {
			// add to valid coin
			validCoins = append(validCoins, entities.CoinDate{
				IDCoin:       key,
				LatestCheck:  parseTime(info["latestCheck"]),
				UpdatedAt:    parseTime(info["updatedAt"]),
				RecentDate:   recentDate,
				FurthestDate: furthestDate,
			})

			if len(validCoins) > 2 {
				enoughCoins = true // got the expectation
			}
		}

	}
	for _, coin := range validCoins {
		nextDate := processDates(coin.RecentDate, coin.FurthestDate)

		// Fetch from Coin Gecko maximum 90 days for get hourly prices histories in date range, and then save in db
		coingecko.FetchAndSavePrice(coin.IDCoin, coin.RecentDate, nextDate, appCtx)
		err := process.GenerateOHLCPriceData(appCtx.DB, coin.IDCoin, coin.RecentDate, nextDate, common.Period1H)
		if err != nil {
			fmt.Println("Error generating OHLC Data of", coin.IDCoin, "from", coin.RecentDate, "to", coin.FurthestDate)
		}
		latestCheck := time.Now()
		// Update record in Redis
		err = cache2.AddCoinDataToRedis(appCtx.Cache, context.Background(), "coin_dates_latestCheck", coin.IDCoin, float64(latestCheck.Unix()))
		if err != nil {
			fmt.Println("Error update coin_dates_latestCheck in Redis of", coin.IDCoin, "from", coin.RecentDate, "to", coin.FurthestDate)
		}
		newFields := map[string]interface{}{
			"recentDate":  nextDate,
			"latestCheck": latestCheck.Format(time.RFC3339),
			"updatedAt":   time.Now().Format(time.RFC3339),
		}
		err = cache2.UpdateCoinDatesInRedis(appCtx.Cache, coin.IDCoin, newFields)
		if err != nil {
			if err != nil {
				fmt.Println("Error update coin data in Redis of", coin.IDCoin, "from", coin.RecentDate, "to", coin.FurthestDate)
			}
		}
	}

}

// convert string to  time.Time
func parseTime(timeStr string) time.Time {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		fmt.Println("Error parse time %s", timeStr, err)
	}
	return t
}

func dateDiffInDays(t1, t2 time.Time) int {
	const day = 24 * time.Hour
	diff := t2.Sub(t1)

	return int(diff.Hours() / 24)
}

func processDates(unixDate1, unixDate2 int64) int64 {
	date1 := time.Unix(unixDate1, 0)
	date2 := time.Unix(unixDate2, 0)

	// Verify both is same date
	if date1.Equal(date2) {
		return -1
	}

	// Calc days range
	diffDays := dateDiffInDays(date1, date2)

	// Verify if diff > 90
	if diffDays > 90 {
		return date1.Add(90 * 24 * time.Hour).Unix()
	}

	// return second date
	return unixDate2
}
