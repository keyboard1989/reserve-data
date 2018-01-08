package exchange

type OkexInfo struct {
	Data []struct {
		BaseCurrency   int     `json:"baseCurrency"`
		PricePrecision int     `json:"maxPriceDigit"`
		AmounPrecision int     `json:"maxSizeDigit"`
		MinAmount      float32 `json:"minTradeSize"`
		Symbol         string  `json:"symbol"`
	}
}

type OkexDepth struct {
	Asks [][]float32 `json:"asks"`
	Bids [][]float32 `json:"bids"`
}
