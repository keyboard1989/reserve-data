package market

import (
	"errors"
)

func TestGetPriceWithError(base, quote string) (PriceData, error) {
	return PriceData{}, errors.New("test")
}

func TestGetPrice(base, quote string) (PriceData, error) {
	return PriceData{
		Version:   1,
		Timestamp: "123456789",
		ExchangeData: OnePairData{
			"1": ExchangeData{
				Valid:     true,
				Timestamp: "1",
				BuyPrices: []Price{
					Price{10.1, 0.2}, Price{1.2, 1},
				},
				SellPrices: []Price{
					Price{3.4, 0.21}, Price{2.0, 1.1},
				},
			},
			"2": ExchangeData{
				Valid:     false,
				Timestamp: "1",
				BuyPrices: []Price{
					Price{10.1, 0.3}, Price{1.2, 2},
				},
				SellPrices: []Price{
					Price{3.4, 0.31}, Price{2.0, 2.1},
				},
			},
		},
	}, nil
}
