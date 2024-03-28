package unit

import (
	"SVN-interview/cmd/process"
	"SVN-interview/internal/common"
	"SVN-interview/internal/entities"
	"reflect"
	"testing"
)

func TestCalculateOHLC(t *testing.T) {
	// Mock dữ liệu
	prices := []entities.CoinPrice{
		{IDCoin: "bitcoin", Time: 1609459200, Price: 29000}, // Giả sử 2021-01-01 00:00:00 UTC
		{IDCoin: "bitcoin", Time: 1609462800, Price: 29500}, // 2021-01-01 01:00:00 UTC, 1 giờ sau
		{IDCoin: "bitcoin", Time: 1609466400, Price: 28500}, // 2021-01-01 02:00:00 UTC, 2 giờ sau
	}
	period := common.Period1H

	got := process.CalculateOHLC(prices, period)

	// Define expect
	want := []entities.CoinOHLC{
		{IDCoin: "bitcoin", Time: 1609459200000, Open: 29000, High: 29500, Low: 29000, Close: 29500},
		{IDCoin: "bitcoin", Time: 1609462800000, Open: 29000, High: 29500, Low: 28500, Close: 28500},
		{IDCoin: "bitcoin", Time: 1609466400000, Open: 29500, High: 29500, Low: 28500, Close: 28500},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("CalculateOHLC() got = %v, want %v", got, want)
	}
}
