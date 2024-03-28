package repository

import (
	"SVN-interview/internal/common"
	"SVN-interview/internal/entities"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"sort"
	"time"
)

//func SaveCoinOHLCData(db *gorm.DB, idCoin string, ohlcData [][]float64) {
//	for i, ohlc := range ohlcData {
//
//		output := fmt.Sprintf("%s %i %f %f %f %f", idCoin, ohlc[0], ohlc[1], ohlc[2], ohlc[3], ohlc[4])
//		println(output)
//
//		record := entities.CoinOHLC{
//			IDCoin: idCoin,
//			Time:   int64(ohlc[0]),
//			Open:   ohlc[1],
//			High:   ohlc[2],
//			Low:    ohlc[3],
//			Close:  ohlc[4],
//		}
//
//		//if err := db.Clauses(clause.OnConflict{DoNothing: true}).(&record).Error; err != nil {
//		//	// Log lỗi và tiếp tục vòng lặp thay vì trả về lỗi
//		//	fmt.Printf("Error inserting record %d: %v\n", i, err)
//		//	continue // Tiếp tục với bản ghi tiếp theo trong mảng
//		//}
//
//		err := db.Create(&record)
//		if err != nil {
//			fmt.Printf("Error inserting record %d: %v\n", i, err)
//		}
//
//	}
//
//}

func SaveCoinOHLCData(db *gorm.DB, idCoin string, ohlcData [][]float64) {

	var records []entities.CoinOHLC

	for _, ohlc := range ohlcData {
		record := entities.CoinOHLC{
			IDCoin: idCoin,
			Time:   int64(ohlc[0]),
			Open:   ohlc[1],
			High:   ohlc[2],
			Low:    ohlc[3],
			Close:  ohlc[4],
		}
		records = append(records, record)
	}

	// Sử dụng OnConflict để update nếu gặp conflict trên cặp khóa unique (IDCoin và Time), hoặc insert nếu không có conflict
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "idCoin"}, {Name: "time"}},                  // Định nghĩa cột bị conflict
		DoUpdates: clause.AssignmentColumns([]string{"open", "high", "low", "close"}), // Cập nhật giá trị nếu conflict
	}).Create(&records).Error // Chèn records

	if err != nil {
		log.Printf("Error upserting records: %v", err)
	}

}

func LoadOHLCData(db *gorm.DB, idCoin string, startDate, endDate time.Time, period common.PeriodEnum) ([]entities.CoinOHLC, error) {
	var results []entities.CoinOHLC
	var err error

	// Chuyển đổi startDate và endDate sang Unix timestamp để so sánh
	startTs := startDate.Unix() * 1000
	endTs := endDate.Unix() * 1000
	intervalSeconds := getIntervalSeconds(period) * 1000
	//err = db.Raw(`
	//        SELECT
	//            FLOOR(EXTRACT(EPOCH FROM "time") / ?) * ? as time,
	//            FIRST_VALUE(open) OVER w as open,
	//            MAX(high) OVER w as high,
	//            MIN(low) OVER w as low,
	//            LAST_VALUE(close) OVER w as close
	//        FROM coin_ohlc
	//        WHERE idCoin = ? AND "time" BETWEEN ? AND ?
	//        WINDOW w AS (PARTITION BY FLOOR(EXTRACT(EPOCH FROM "time") / ?))
	//        ORDER BY time
	//    `, intervalSeconds, intervalSeconds, idCoin, startTs, endTs, intervalSeconds).Scan(&results).Error

	err = db.Raw(`
            SELECT 
                "idCoin" ,
                floor(extract(epoch from to_timestamp("time" / ?))) * ? as time,
            	FIRST_VALUE(open) OVER w as open, 
  				MAX(high) OVER w as high,
  				MIN(low) OVER w as low,
  				LAST_VALUE(close) OVER w as close
            FROM coin_ohlc
            WHERE "idCoin"= ? AND "time" BETWEEN ? AND ?
            WINDOW w AS (PARTITION BY floor(extract(epoch from to_timestamp("time" / ?))))
            ORDER BY time
        `, intervalSeconds, intervalSeconds, idCoin, startTs, endTs, intervalSeconds).Scan(&results).Error

	if err != nil {
		log.Printf("Error load records: %v", err)
		return nil, err
	}
	reduceData := reduceData(results, period)
	return reduceData, nil
}
func getIntervalSeconds(period common.PeriodEnum) int64 {
	var intervalSeconds int64
	switch period {
	case common.Period30M:
		intervalSeconds = 1800
	case common.Period1H:
		intervalSeconds = 3600
	case common.Period1D:
		intervalSeconds = 86400
	}
	return intervalSeconds
}
func reduceData(data []entities.CoinOHLC, interval common.PeriodEnum) []entities.CoinOHLC {
	var reducedData []entities.CoinOHLC
	tempMap := make(map[int64]entities.CoinOHLC)

	// Định nghĩa biến intervalSeconds dựa trên interval
	intervalSeconds := getIntervalSeconds(interval)

	for _, ohlc := range data {
		// Làm tròn thời gian xuống khoảng thời gian gần nhất dựa trên interval
		roundedTime := ohlc.Time - (ohlc.Time % intervalSeconds)

		// Kiểm tra xem khoảng thời gian này đã có trong tempMap chưa
		if existing, exists := tempMap[roundedTime]; exists {
			// Cập nhật High và Low nếu cần
			if ohlc.High > existing.High {
				existing.High = ohlc.High
			}
			if ohlc.Low < existing.Low {
				existing.Low = ohlc.Low
			}
			// Close luôn được cập nhật từ record cuối cùng trong khoảng thời gian
			existing.Close = ohlc.Close
			tempMap[roundedTime] = existing
		} else {
			// Nếu khoảng thời gian này chưa có, thêm mới vào tempMap
			tempMap[roundedTime] = ohlc
		}
	}

	// Chuyển dữ liệu từ tempMap vào reducedData
	for _, v := range tempMap {
		reducedData = append(reducedData, v)
	}
	// Sắp xếp reducedData trước khi trả về
	sort.Slice(reducedData, func(i, j int) bool {
		return reducedData[i].Time < reducedData[j].Time
	})

	return reducedData
}
func SaveOHLCFromPriceData(db *gorm.DB, ohlcData []entities.CoinOHLC) error {
	// Sử dụng OnConflict để update nếu gặp conflict trên cặp khóa unique (IDCoin và Time), hoặc insert nếu không có conflict
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "idCoin"}, {Name: "time"}},                  // Định nghĩa cột bị conflict
		DoUpdates: clause.AssignmentColumns([]string{"open", "high", "low", "close"}), // Cập nhật giá trị nếu conflict
	}).Create(&ohlcData).Error // Chèn dữ liệu

	if err != nil {
		log.Printf("Error upserting OHLC data: %v", err)
		return err
	}

	return nil
}
