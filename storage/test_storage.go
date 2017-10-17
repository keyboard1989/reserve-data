package storage

import (
	"github.com/KyberNetwork/reserve-data/market"
)

type TestStorage struct{}

func NewTestStorage() *TestStorage {
	return &TestStorage{}
}

func (self *TestStorage) GetVersion() int64 {
	return 1
}

func (self *TestStorage) GetTimestamp(version int64) string {
	return "timestamp"
}

func (self *TestStorage) GetOnePairData(base, quote string, version int64) (market.OnePairData, error) {
	return market.OnePairData{
		"1": market.ExchangeData{
			Valid:     true,
			Timestamp: "1",
			BuyPrices: []market.Price{
				market.Price{10.1, 0.2}, market.Price{1.2, 1},
			},
			SellPrices: []market.Price{
				market.Price{3.4, 0.21}, market.Price{2.0, 1.1},
			},
		},
		"2": market.ExchangeData{
			Valid:     false,
			Timestamp: "1",
			BuyPrices: []market.Price{
				market.Price{10.1, 0.3}, market.Price{1.2, 2},
			},
			SellPrices: []market.Price{
				market.Price{3.4, 0.31}, market.Price{2.0, 2.1},
			},
		},
	}, nil
}
