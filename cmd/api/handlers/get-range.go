package handlers

import (
	"SVN-interview/internal/common"
	"SVN-interview/pkg/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
func GetHistoriesHandler1(c *gin.Context, appCtx *common.AppContext) ([]CoinOHLCResponse, error) {
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

	// Trả về dữ liệu dưới dạng JSON
	c.JSON(http.StatusOK, responseData)
	return responseData, nil
}
