package coingecko

import (
	"SVN-interview/internal/common"
	"SVN-interview/pkg/repository"
	"encoding/json"
	"fmt"
	"strconv"
)

func FetchCoinPrice(idCoin string, startTs, endTs int64) (string, error) {
	startStr := strconv.FormatInt(startTs, 10)
	endStr := strconv.FormatInt(endTs, 10)
	endpoint := "/coins/" + idCoin + "/market_chart/range?vs_currency=usd&from=" + startStr + "&to=" + endStr + "&precision=3"
	return fetchGecko(endpoint)
}
func FetchAndSavePrice(idCoin string, startTs int64, endTs int64, appCtx *common.AppContext) string {
	body, err := FetchCoinPrice(idCoin, startTs, endTs)
	if err != nil {
		fmt.Printf("Error fetching coin price: %v\n", err)
		return ""
	}

	var result map[string][][]float64
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		fmt.Printf("Error parsing JSON body: %v\n", err)
		return ""
	}

	prices := result["prices"]
	for _, price := range prices {
		repository.SaveCoinPriceData(appCtx.DB, idCoin, price[0], price[1])
	}
	return body
}
