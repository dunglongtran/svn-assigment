package entities

import (
	"time"
)

type CoinOHLC struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	IDCoin    string    `gorm:"column:idCoin;type:varchar(50);not null;index:idx_idcoin_time,unique"`
	Time      int64     `gorm:"not null;index:idx_idcoin_time,unique"`
	Open      float64   `gorm:"type:numeric(20,10);not null"`
	High      float64   `gorm:"type:numeric(20,10);not null"`
	Low       float64   `gorm:"type:numeric(20,10);not null"`
	Close     float64   `gorm:"type:numeric(20,10);not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName sets the insert table name for this struct type
func (CoinOHLC) TableName() string {
	return "coin_ohlc"
}
