package repository

import (
	"SVN-interview/internal/entities"
	"fmt"
	"gorm.io/gorm"
)

func GetCoinIDBySymbol(db *gorm.DB, symbol string) (string, error) {
	var coin entities.CoinInfo
	result := db.Where("symbol = ?", symbol).First(&coin)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("coin with symbol %s not found", symbol)
		}
		return "", result.Error
	}
	return coin.ID, nil
}
