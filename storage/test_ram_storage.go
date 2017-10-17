package storage

import (
	"github.com/KyberNetwork/reserve-data/market"
	"sync"
)

func NewTestRamStorage() *RamStorage {
	return &RamStorage{
		mu:        sync.RWMutex{},
		version:   8891,
		timestamp: map[int64]string{8891: "timestamp1"},
		data: map[int64]market.AllPriceData{
			8891: market.AllPriceData{
				Version:   1,
				Timestamp: "allpricetimestamp",
				AllPairData: map[string]market.OnePairData{
					"eth-omg": market.OnePairData{
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
					},
				},
			},
		},
	}
}
