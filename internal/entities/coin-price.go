package entities

import "time"

// CoinPrice represents the price data for a coin
type CoinPrice struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	IDCoin    string    `gorm:"column:idCoin;type:varchar(50);not null;index:coin_price_idcoin_time_key,unique"`
	Time      int64     `gorm:"not null;index:coin_price_idcoin_time_key,unique"`
	Price     float64   `gorm:"type:numeric(20,10);not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName sets the insert table name for this struct type
func (CoinPrice) TableName() string {
	return "coin_price"
}
