package handlers

import (
	"SVN-interview/internal/common"
	"SVN-interview/pkg/coingecko"
	"SVN-interview/pkg/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// GetHistoriesParams
type GetRangeParams struct {
	StartDate time.Time `form:"start_date" binding:"required" time_format:"2006-01-02"`
	EndDate   time.Time `form:"end_date" binding:"required" time_format:"2006-01-02"`
	Symbol    string    `form:"symbol" binding:"required"`
}

// CoinRangeResponse
type CoinRangeResponse struct {
	//ID        int       `json:"id"`
	IDCoin string  `json:"idCoin"`
	Time   int64   `json:"time"`
	Price  float64 `json:"price"`
	//CreatedAt time.Time `json:"createdAt"`
	//UpdatedAt time.Time `json:"updatedAt"`
}

// @BasePath /

// GetRangeHandler godoc
// @Summary Get historical price data for a coin
// @Description get price a coin on market
// @Tags coin
// @Accept  json
// @Produce  json
// @Param   symbol     query    string     true  "Coin Symbol"
// @Param   start_date query    string     true  "Start Date" format(date)
// @Param end_date query string true "End Date" format(date)
// @Success 200 {object} []CoinRangeResponse
// @Router /get-range [get]
func GetRangeHandler(c *gin.Context, appCtx *common.AppContext) ([]CoinOHLCResponse, error) {
	var params GetHistoriesParams
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, nil
	}

	processData(params, appCtx)
	idCoin, _ := repository.GetCoinIDBySymbol(appCtx.DB, params.Symbol)
	data, err := repository.LoadOHLCData(appCtx.DB, idCoin, params.StartDate, params.EndDate, params.Period)
	if err != nil {
		return nil, err
	}
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
	// calc distance of start_date and end_date
	days := params.EndDate.Sub(params.StartDate).Hours() / 24
	daysAgo := CalculateDaysAgo(int(days))
	idCoin, _ := repository.GetCoinIDBySymbol(appCtx.DB, params.Symbol)
	if daysAgo == -1 {
		coingecko.FetchAndSaveOHLC(idCoin, "max", appCtx)

	} else if daysAgo > 1 {
		coingecko.FetchAndSaveOHLC(idCoin, strconv.Itoa(daysAgo), appCtx)
	}
	coingecko.FetchAndSaveOHLC(idCoin, strconv.Itoa(1), appCtx)
}
