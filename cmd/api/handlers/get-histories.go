package handlers

import (
	"SVN-interview/internal/common"
	"SVN-interview/pkg/cache"
	"SVN-interview/pkg/coingecko"
	"SVN-interview/pkg/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// GetHistoriesParams
type GetHistoriesParams struct {
	StartDate time.Time         `form:"start_date" binding:"required" time_format:"2006-01-02"`
	EndDate   time.Time         `form:"end_date" binding:"required" time_format:"2006-01-02"`
	Symbol    string            `form:"symbol" binding:"required"`
	Period    common.PeriodEnum `form:"period" binding:"required,oneof=30M 1H 1D"`
}

// CoinOHLCResponse
type CoinOHLCResponse struct {
	//ID        int       `json:"id"`
	IDCoin string  `json:"idCoin"`
	Time   int64   `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	//CreatedAt time.Time `json:"createdAt"`
	//UpdatedAt time.Time `json:"updatedAt"`
}

// @BasePath /

// GetHistoriesHandler godoc
// @Summary Get historical OHLC data for a coin
// @Description get historical open, high, low, and close data for a coin
// @Tags coin
// @Accept  json
// @Produce  json
// @Param   symbol     query    string     true  "Coin Symbol"
// @Param   start_date query    string     true  "Start Date" format(date)
// @Param end_date query string true "End Date" format(date)
// @Param period query string true "Data period (30M, 1H, 1D)"
// @Success 200 {object} []CoinOHLCResponse
// @Router /get-histories [get]
func GetHistoriesHandler(c *gin.Context, appCtx *common.AppContext) ([]CoinOHLCResponse, error) {
	var params GetHistoriesParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, nil
	}

	processData(params, appCtx)
	idCoin, _ := repository.GetCoinIDBySymbol(appCtx.DB, params.Symbol)
	data, err := repository.LoadOHLCData(appCtx.DB, idCoin, params.StartDate, params.EndDate, params.Period)
	if err != nil {
		fmt.Println("Error get histories of ", idCoin, "from ", params.StartDate, "to ", params.EndDate)
		return nil, err
	}

	//Update Redis
	_ = cache.IncreaseCoinAccess(appCtx.Cache, idCoin)
	farDate := findFarthestDate(params.StartDate, params.EndDate)
	_ = cache.UpdateFurthestDateInRedis(appCtx.Cache, idCoin, farDate.Unix())
	newFields := map[string]interface{}{
		"updatedAt": time.Now().Format(time.RFC3339),
	}
	_ = cache.UpdateCoinDatesInRedis(appCtx.Cache, idCoin, newFields)

	var responseData []CoinOHLCResponse
	for _, ohlc := range data {
		responseData = append(responseData, CoinOHLCResponse{
			//ID:        ohlc.ID,
			IDCoin: ohlc.IDCoin,
			Time:   ohlc.Time,
			Open:   ohlc.Open,
			High:   ohlc.High,
			Low:    ohlc.Low,
			Close:  ohlc.Close,
			//CreatedAt: ohlc.CreatedAt,
			//UpdatedAt: ohlc.UpdatedAt,
		})
	}

	// return JSON response
	c.JSON(http.StatusOK, responseData)
	return responseData, nil
}
func CalculateDaysAgo(days int) int {
	//days := int(endDate.Sub(startDate).Hours() / 24)

	switch {
	case days <= 1:
		return 1
	case days <= 7:
		return 7
	case days <= 14:
		return 14
	case days <= 30:
		return 30
	case days <= 90:
		return 90
	case days <= 180:
		return 180
	case days <= 365:
		return 365
	default:
		return -1 // "max"
	}
}

func processData(params GetHistoriesParams, appCtx *common.AppContext) {
	farDate := findFarthestDate(params.StartDate, params.EndDate)
	now := time.Now()
	// calc distance of start_date and end_date
	//days := params.EndDate.Sub(params.StartDate).Hours() / 24
	days := now.Sub(farDate).Hours() / 24
	daysAgo := CalculateDaysAgo(int(days))
	idCoin, _ := repository.GetCoinIDBySymbol(appCtx.DB, params.Symbol)
	if daysAgo == -1 {
		coingecko.FetchAndSaveOHLC(idCoin, "max", appCtx)
		//coingecko.FetchAndSaveOHLC(idCoin, strconv.Itoa(365), appCtx)

	} else if daysAgo > 1 {
		coingecko.FetchAndSaveOHLC(idCoin, strconv.Itoa(daysAgo), appCtx)
	}
	coingecko.FetchAndSaveOHLC(idCoin, strconv.Itoa(1), appCtx)
}
func findOldest() {

}
func findFarthestDate(date1, date2 time.Time) time.Time {
	now := time.Now()
	// Calc date range to now
	diff1 := now.Sub(date1)
	diff2 := now.Sub(date2)

	// Compare and then get the greater
	if diff1 > diff2 {
		return date1
	}
	return date2
}
