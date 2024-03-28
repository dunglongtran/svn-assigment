package entities

import "time"

type CoinDate struct {
	IDCoin       string
	LatestCheck  time.Time
	UpdatedAt    time.Time
	RecentDate   int64
	FurthestDate int64
}
