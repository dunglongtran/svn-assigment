package process

import (
	"SVN-interview/internal/common"
	"SVN-interview/internal/entities"
	"SVN-interview/pkg/repository"
	"gorm.io/gorm"
	"math"
	"sort"
)

func GenerateOHLCPriceData(db *gorm.DB, idCoin string, startTs, endTs int64, period common.PeriodEnum) error {
	prices, err := repository.LoadCoinPriceData(db, idCoin, startTs, endTs)
	if err != nil {
		return err
	}

	ohlcData := CalculateOHLC(prices, period)
	err = repository.SaveOHLCFromPriceData(db, ohlcData)
	if err != nil {
		return err
	}

	return nil
}
func CalculateOHLC(prices []entities.CoinPrice, period common.PeriodEnum) []entities.CoinOHLC {
	var ohlcData []entities.CoinOHLC

	// Giả sử mỗi khoảng thời gian được biểu diễn bằng số giây
	var periodSeconds int64
	switch period {
	case common.Period30M:
		periodSeconds = 1800 // 30 phút = 1800 giây
	case common.Period1H:
		periodSeconds = 3600 // 1 giờ = 3600 giây
	case common.Period1D:
		periodSeconds = 86400 // 1 ngày = 86400 giây
	default:
		periodSeconds = 1800 // Mặc định là 30 phút nếu không xác định được
	}

	// Tạo một map để lưu trữ dữ liệu OHLC theo khoảng thời gian
	ohlcMap := make(map[int64]*entities.CoinOHLC)
	var prevIndex int64
	if len(prices) > 0 {
		prevIndex = prices[0].Time / (periodSeconds)
	}

	for i, price := range prices {
		// Xác định khoảng thời gian cho mỗi giá
		periodIndex := price.Time / (periodSeconds)

		ohlc, exists := ohlcMap[periodIndex]
		if !exists {
			ohlc = &entities.CoinOHLC{
				IDCoin: price.IDCoin,
				Time:   periodIndex * periodSeconds * 1000, // Chuyển lại sang ms
				Open:   price.Price,
				High:   price.Price,
				Low:    price.Price,
				Close:  price.Price,
			}
			if i > 0 {
				ohlc.Open = prices[i-1].Price
				ohlcMap[prevIndex].Close = price.Price
				ohlcMap[prevIndex].High = findMax(ohlcMap[prevIndex].High, ohlcMap[prevIndex].Open, ohlcMap[prevIndex].Close)
				ohlcMap[prevIndex].Low = findMin(ohlcMap[prevIndex].Low, ohlcMap[prevIndex].Open, ohlcMap[prevIndex].Close)
			}
			if i == len(prices)-1 {
				ohlc.High = findMax(ohlc.High, ohlc.Open, ohlc.Close)
				ohlc.Low = findMin(ohlc.Low, ohlc.Open, ohlc.Close)
			}
			ohlcMap[periodIndex] = ohlc
			prevIndex = periodIndex
		} else {
			ohlc.High = findMax(ohlc.High, price.Price, ohlc.Open, ohlc.Close)
			ohlc.Low = findMin(ohlc.Low, price.Price, ohlc.Open, ohlc.Close)
		}
	}

	// Chuyển dữ liệu từ map vào slice để trả về
	for _, ohlc := range ohlcMap {
		ohlcData = append(ohlcData, *ohlc)
	}

	// Sắp xếp ohlcData theo thời gian
	sort.Slice(ohlcData, func(i, j int) bool {
		return ohlcData[i].Time < ohlcData[j].Time
	})

	return ohlcData
}
func findMax(numbers ...float64) float64 {
	if len(numbers) == 0 {
		return math.Inf(-1) // Trả về âm vô cùng nếu không có số nào
	}

	max := numbers[0] // Giả sử số đầu tiên là số lớn nhất

	for _, num := range numbers {
		if num > max {
			max = num
		}
	}

	return max
}
func findMin(numbers ...float64) float64 {
	if len(numbers) == 0 {
		return math.Inf(1) // Trả về dương vô cùng nếu không có số nào
	}

	min := numbers[0] // Giả sử số đầu tiên là số nhỏ nhất

	for _, num := range numbers {
		if num < min {
			min = num
		}
	}

	return min
}
