package repository

import (
	"SVN-interview/internal/entities"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func SaveCoinPriceData(db *gorm.DB, idCoin string, timestamp float64, price float64) error {
	timeInt := int64(timestamp / 1000) // Chuyển timestamp sang giây
	record := entities.CoinPrice{
		IDCoin: idCoin,
		Time:   timeInt,
		Price:  price,
	}

	// Upsert logic
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "idCoin"}, {Name: "time"}},
		DoUpdates: clause.AssignmentColumns([]string{"price"}),
	}).Create(&record).Error

	if err != nil {
		fmt.Printf("Error saving coin price data: %v\n", err)
		return err
	}

	return nil
}

func FetchCoinPriceData(db *gorm.DB, idCoin string, startTs, endTs int64) ([]entities.CoinPrice, error) {
	var prices []entities.CoinPrice
	err := db.Where("\"idCoin\" = ? AND time >= ? AND time <= ?", idCoin, startTs, endTs).Order("time asc").Find(&prices).Error
	if err != nil {
		return nil, err
	}
	return prices, nil
}
