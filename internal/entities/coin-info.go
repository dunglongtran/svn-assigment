package entities

import (
	"time"
)

type CoinInfo struct {
	ID        string    `gorm:"primaryKey"`
	Symbol    string    `gorm:"unique;not null"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	LatestAt  time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	Order     int       `gorm:"default:1"`
}

// TableName sets the insert table name for this struct type
func (CoinInfo) TableName() string {
	return "coin_info"
}
