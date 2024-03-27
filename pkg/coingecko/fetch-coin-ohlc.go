package coingecko

import (
	"SVN-interview/internal/common"
	"SVN-interview/pkg/repository"
	"encoding/json"
	"fmt"
)

func FetchCoinOHLC(idCoin string, days string) (string, error) {
	endpoint := "/coins/" + idCoin + "/ohlc?vs_currency=usd&days=" + days
	return fetchGecko(endpoint)
}
func FetchAndSaveOHLC(idCoin string, days string, appCtx *common.AppContext) string {
	body, _ := FetchCoinOHLC(idCoin, days)
	var ohlcData [][]float64

	if err := json.Unmarshal([]byte(body), &ohlcData); err != nil {
		fmt.Printf("Error parsing JSON body: %v\n", err)
	}
	repository.SaveCoinOHLCData(appCtx.DB, idCoin, ohlcData)
	//if err != nil {
	//	return ""
	//}
	return body
}
